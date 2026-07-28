package config

import (
	"strings"
	"testing"
	"time"
)

func TestEnvStr(t *testing.T) {
	t.Setenv("X_ZOT_EPHEMERAL_TTL_TEST", "")
	if got := envStr("X_ZOT_EPHEMERAL_TTL_TEST", "fallback"); got != "fallback" {
		t.Errorf("empty env should return default, got %q", got)
	}
	t.Setenv("X_ZOT_EPHEMERAL_TTL_TEST", "value")
	if got := envStr("X_ZOT_EPHEMERAL_TTL_TEST", "fallback"); got != "value" {
		t.Errorf("set env should win, got %q", got)
	}
}

func TestLoad(t *testing.T) {
	keys := []string{
		"ZOT_EPHEMERAL_TTL_LISTEN",
		"ZOT_EPHEMERAL_TTL_REDIS_URL",
		"ZOT_EPHEMERAL_TTL_SWEEP_INTERVAL",
		"ZOT_EPHEMERAL_TTL_ZOT_URL",
		"ZOT_EPHEMERAL_TTL_DEFAULT_TTL",
		"ZOT_EPHEMERAL_TTL_MAX_TTL",
	}

	t.Run("defaults", func(t *testing.T) {
		for _, k := range keys {
			t.Setenv(k, "")
		}
		got, err := Load()
		if err != nil {
			t.Fatalf("defaults: unexpected error: %v", err)
		}
		want := Config{
			Listen:        ":8080",
			RedisURL:      "redis://redis:6379/0",
			SweepInterval: 30 * time.Second,
			ZotURL:        "http://zot:5000",
			DefaultTTL:    24 * time.Hour,
			MaxTTL:        24 * time.Hour,
		}
		if got != want {
			t.Fatalf("defaults: got %+v want %+v", got, want)
		}
	})

	t.Run("overrides", func(t *testing.T) {
		t.Setenv("ZOT_EPHEMERAL_TTL_LISTEN", "127.0.0.1:9090")
		t.Setenv("ZOT_EPHEMERAL_TTL_REDIS_URL", "rediss://cache.internal:6380/3")
		t.Setenv("ZOT_EPHEMERAL_TTL_SWEEP_INTERVAL", "5s")
		t.Setenv("ZOT_EPHEMERAL_TTL_ZOT_URL", "http://zot:5000/")
		t.Setenv("ZOT_EPHEMERAL_TTL_DEFAULT_TTL", "1h")
		t.Setenv("ZOT_EPHEMERAL_TTL_MAX_TTL", "2h")
		got, err := Load()
		if err != nil {
			t.Fatalf("overrides: unexpected error: %v", err)
		}
		want := Config{
			Listen:        "127.0.0.1:9090",
			RedisURL:      "rediss://cache.internal:6380/3",
			SweepInterval: 5 * time.Second,
			ZotURL:        "http://zot:5000",
			DefaultTTL:    time.Hour,
			MaxTTL:        2 * time.Hour,
		}
		if got != want {
			t.Fatalf("overrides: got %+v want %+v", got, want)
		}
	})
}

// TestRedactedRedisURL: the startup banner logs the Redis URL, which commonly
// carries a password, so it must go through redaction first.
func TestRedactedRedisURL(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"no credentials", "redis://redis:6379/0", "redis://redis:6379/0"},
		{"password only", "redis://:hunter2@redis:6379/0", "redis://:xxxxx@redis:6379/0"},
		{"user and password", "rediss://app:hunter2@cache:6380/1", "rediss://app:xxxxx@cache:6380/1"},
		{"unparseable", "redis://%zz@host", "(unparseable redis url)"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := (Config{RedisURL: tc.in}).RedactedRedisURL(); got != tc.want {
				t.Errorf("got %q want %q", got, tc.want)
			}
		})
	}
}

// TestLoadReportsEveryBadDuration checks that a malformed duration surfaces as
// an error naming the offending variable, and that all of them are reported
// rather than just the first.
func TestLoadReportsEveryBadDuration(t *testing.T) {
	t.Setenv("ZOT_EPHEMERAL_TTL_SWEEP_INTERVAL", "not-a-duration")
	t.Setenv("ZOT_EPHEMERAL_TTL_MAX_TTL", "also-bad")

	cfg, err := Load()
	if err == nil {
		t.Fatal("expected an error for malformed durations, got nil")
	}
	for _, want := range []string{"ZOT_EPHEMERAL_TTL_SWEEP_INTERVAL", "ZOT_EPHEMERAL_TTL_MAX_TTL"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error should name %s; got %v", want, err)
		}
	}
	// The bad values fall back to their defaults so the caller sees a usable
	// Config alongside the error.
	if cfg.SweepInterval != 30*time.Second || cfg.MaxTTL != 24*time.Hour {
		t.Errorf("bad durations should fall back to defaults, got %+v", cfg)
	}
}

func TestEnvDur(t *testing.T) {
	t.Setenv("X_ZOT_EPHEMERAL_TTL_DUR", "")
	got, err := envDur("X_ZOT_EPHEMERAL_TTL_DUR", 5*time.Second)
	if err != nil || got != 5*time.Second {
		t.Errorf("empty env: got %v, %v", got, err)
	}
	t.Setenv("X_ZOT_EPHEMERAL_TTL_DUR", "2m")
	got, err = envDur("X_ZOT_EPHEMERAL_TTL_DUR", time.Second)
	if err != nil || got != 2*time.Minute {
		t.Errorf("set env: got %v, %v", got, err)
	}
	t.Setenv("X_ZOT_EPHEMERAL_TTL_DUR", "nope")
	if _, err := envDur("X_ZOT_EPHEMERAL_TTL_DUR", time.Second); err == nil {
		t.Error("malformed duration: expected error, got nil")
	}
}
