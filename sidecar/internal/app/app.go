// Package app wires the concrete components together and runs them until the
// context is cancelled. It is separate from package main so the full
// startup/shutdown path can be tested; main only loads config and handles
// signals.
package app

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/nullbytelabs/zot-ephemeral-ttl/internal/config"
	"github.com/nullbytelabs/zot-ephemeral-ttl/internal/reaper"
	"github.com/nullbytelabs/zot-ephemeral-ttl/internal/registry"
	"github.com/nullbytelabs/zot-ephemeral-ttl/internal/server"
	"github.com/nullbytelabs/zot-ephemeral-ttl/internal/store"
)

// storeBackend is the union of persistence capabilities app injects into the
// server and reaper, so neither consumer names the concrete adapter.
type storeBackend interface {
	server.Store // Upsert
	reaper.Store // Expired, Delete
	io.Closer
}

// Run connects to Redis, starts the sweep loop and HTTP server, and blocks
// until ctx is cancelled (or the HTTP server fails). It returns the first
// error encountered during startup or shutdown.
func Run(ctx context.Context, cfg config.Config) error {
	var st storeBackend
	st, err := store.Open(cfg.RedisURL)
	if err != nil {
		return err
	}
	defer func() {
		if err := st.Close(); err != nil {
			log.Printf("closing store: %v", err)
		}
	}()

	// Bind up front so a bad address fails synchronously and a ":0" caller can
	// read back the chosen port.
	ln, err := net.Listen("tcp", cfg.Listen)
	if err != nil {
		return err
	}

	srv := server.New(st, cfg.DefaultTTL, cfg.MaxTTL, time.Now)
	httpSrv := &http.Server{
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	rp := reaper.New(cfg.SweepInterval, st, registry.New(cfg.ZotURL))
	go rp.Run(ctx)

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", ln.Addr())
		serveErr <- httpSrv.Serve(ln)
	}()

	select {
	case <-ctx.Done():
	case err := <-serveErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	log.Printf("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return httpSrv.Shutdown(shutdownCtx)
}
