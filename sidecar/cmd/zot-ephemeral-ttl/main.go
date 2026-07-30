// Command zot-ephemeral-ttl is a sidecar for zot that deletes tags whose names
// encode a TTL (e.g. ":1h", ":30m", ":7d"). It loads config, handles signals,
// and hands off to package app; package ttl documents the expiry policy.
package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/nullbytelabs/zot-ephemeral-ttl/internal/app"
	"github.com/nullbytelabs/zot-ephemeral-ttl/internal/config"
)

// Build metadata, set via -ldflags "-X main.version=..." in the Dockerfile.
var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	log.Printf("zot-ephemeral-ttl %s (commit=%s built=%s) starting: listen=%s redis=%s sweep=%s zot=%s default_ttl=%s max_ttl=%s",
		version, commit, buildDate,
		cfg.Listen, cfg.RedactedRedisURL(), cfg.SweepInterval, cfg.ZotURL, cfg.DefaultTTL, cfg.MaxTTL)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := app.Run(ctx, cfg); err != nil {
		log.Fatalf("run: %v", err)
	}
}
