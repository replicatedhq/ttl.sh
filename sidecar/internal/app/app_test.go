package app

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/replicatedhq/ttl.sh/sidecar/internal/config"
	"github.com/replicatedhq/ttl.sh/sidecar/internal/redistest"
)

// The startup-path test needs a real Redis; redistest starts one per package.
func TestMain(m *testing.M) { os.Exit(redistest.Run(m)) }

// baseConfig returns a config pointing at that server.
func baseConfig(t *testing.T) config.Config {
	t.Helper()
	return config.Config{
		Listen:        "127.0.0.1:0",
		RedisURL:      redistest.URL(t),
		SweepInterval: time.Hour,
		ZotURL:        "http://127.0.0.1:1",
		DefaultTTL:    24 * time.Hour,
		MaxTTL:        24 * time.Hour,
	}
}

// TestRunCleanShutdown exercises the real startup/shutdown path. SweepInterval
// is an hour so the reaper does its one immediate sweep against an empty store
// and then idles, never reaching the (unreachable) registry.
func TestRunCleanShutdown(t *testing.T) {
	cfg := baseConfig(t)

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- Run(ctx, cfg)
	}()

	// Give Run a moment to bind the listener and start serving.
	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run returned error on clean shutdown: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return within 2s after context cancellation")
	}
}

// TestRunStoreOpenFailure points Run at an address nothing is listening on, so
// the store's startup PING fails and Run returns before binding anything. It
// needs no server of its own.
func TestRunStoreOpenFailure(t *testing.T) {
	cfg := config.Config{
		Listen:        "127.0.0.1:0",
		RedisURL:      "redis://127.0.0.1:1", // connection refused
		SweepInterval: time.Hour,
		ZotURL:        "http://127.0.0.1:1",
		DefaultTTL:    24 * time.Hour,
		MaxTTL:        24 * time.Hour,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- Run(context.Background(), cfg)
	}()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("Run returned nil; expected store-open error")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Run did not return promptly on store-open failure")
	}
}
