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

	"github.com/replicatedhq/ttl.sh/sidecar/internal/config"
	"github.com/replicatedhq/ttl.sh/sidecar/internal/reaper"
	"github.com/replicatedhq/ttl.sh/sidecar/internal/registry"
	"github.com/replicatedhq/ttl.sh/sidecar/internal/server"
	"github.com/replicatedhq/ttl.sh/sidecar/internal/store"
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

	// The sweep loop gets its own cancellable context so it can be stopped on
	// the way out even when shutdown was triggered by a server error rather
	// than by ctx.
	reaperCtx, stopReaper := context.WithCancel(ctx)
	defer stopReaper()

	rp := reaper.New(cfg.SweepInterval, st, registry.New(cfg.ZotURL))
	reaperDone := make(chan struct{})
	go func() {
		defer close(reaperDone)
		rp.Run(reaperCtx)
	}()

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
	shutdownErr := httpSrv.Shutdown(shutdownCtx)

	// Wait for an in-flight sweep to unwind before returning, so it cannot
	// still be talking to Redis when the deferred Close runs. Cancelling
	// reaperCtx aborts the delete request in flight, so this does not block for
	// long.
	stopReaper()
	<-reaperDone
	return shutdownErr
}
