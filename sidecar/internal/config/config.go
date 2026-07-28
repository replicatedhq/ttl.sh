// Package config loads zot-ephemeral-ttl's runtime configuration from the
// environment. All knobs use the ZOT_EPHEMERAL_TTL_ prefix and fall back
// to values suitable for the docker-compose deployment.
package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

type Config struct {
	Listen        string
	RedisURL      string
	SweepInterval time.Duration
	ZotURL        string
	DefaultTTL    time.Duration
	MaxTTL        time.Duration
}

// Load reads configuration from the environment, applying defaults for any
// unset values. Every malformed duration is reported, not just the first.
func Load() (Config, error) {
	var errs []error
	dur := func(key string, def time.Duration) time.Duration {
		d, err := envDur(key, def)
		if err != nil {
			errs = append(errs, err)
			return def
		}
		return d
	}

	cfg := Config{
		Listen:        envStr("ZOT_EPHEMERAL_TTL_LISTEN", ":8080"),
		RedisURL:      envStr("ZOT_EPHEMERAL_TTL_REDIS_URL", "redis://redis:6379/0"),
		SweepInterval: dur("ZOT_EPHEMERAL_TTL_SWEEP_INTERVAL", 30*time.Second),
		ZotURL:        strings.TrimRight(envStr("ZOT_EPHEMERAL_TTL_ZOT_URL", "http://zot:5000"), "/"),
		DefaultTTL:    dur("ZOT_EPHEMERAL_TTL_DEFAULT_TTL", 24*time.Hour),
		MaxTTL:        dur("ZOT_EPHEMERAL_TTL_MAX_TTL", 24*time.Hour),
	}
	return cfg, errors.Join(errs...)
}

// RedactedRedisURL returns RedisURL with any password replaced by "xxxxx",
// safe to write to logs. An unparseable URL is reported as such rather than
// echoed, since it may still contain a credential.
func (c Config) RedactedRedisURL() string {
	u, err := url.Parse(c.RedisURL)
	if err != nil {
		return "(unparseable redis url)"
	}
	return u.Redacted()
}

func envStr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envDur(key string, def time.Duration) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		return def, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("invalid duration for %s=%q: %w", key, v, err)
	}
	return d, nil
}
