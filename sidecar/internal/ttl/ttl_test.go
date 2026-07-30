package ttl

import (
	"math"
	"testing"
	"time"
)

func TestParseTag(t *testing.T) {
	cases := []struct {
		name      string
		reference string
		want      time.Duration
		wantOK    bool
	}{
		{"empty", "", 0, false},
		{"latest", "latest", 0, false},
		{"semver", "v1.2.3", 0, false},
		{"unknown unit", "10y", 0, false},
		{"missing unit", "10", 0, false},
		{"missing number", "h", 0, false},
		{"zero seconds", "0s", 0, false},
		{"negative not allowed by regex", "-5m", 0, false},
		{"trailing junk", "10ss", 0, false},
		{"leading junk", "x10s", 0, false},
		{"upper case unit", "10H", 0, false},
		{"seconds", "30s", 30 * time.Second, true},
		{"minutes", "30m", 30 * time.Minute, true},
		{"hours", "2h", 2 * time.Hour, true},
		{"days", "7d", 7 * 24 * time.Hour, true},
		{"weeks", "2w", 2 * 7 * 24 * time.Hour, true},
		{"single second", "1s", time.Second, true},
		{"leading zeros", "007h", 7 * time.Hour, true},
		{"large", "9999h", 9999 * time.Hour, true},
		{"int64 overflow", "99999999999999999999s", 0, false},
		// 20000w exceeds what time.Duration can represent. Saturating keeps the
		// value positive so ComputeExpiry clamps it; wrapping would make it
		// negative and expire the tag on the spot.
		{"duration overflow saturates", "20000w", time.Duration(math.MaxInt64), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := ParseTag(tc.reference)
			if ok != tc.wantOK {
				t.Fatalf("ok=%v want=%v", ok, tc.wantOK)
			}
			if got != tc.want {
				t.Fatalf("got=%v want=%v", got, tc.want)
			}
		})
	}
}

func TestComputeExpiry(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	defaultTTL := 24 * time.Hour
	maxTTL := 24 * time.Hour

	cases := []struct {
		name      string
		reference string
		want      time.Duration // expected delta from now
	}{
		{"unknown reference uses default", "latest", defaultTTL},
		{"empty reference uses default", "", defaultTTL},
		{"parsed under max", "1h", time.Hour},
		{"parsed equal to max", "24h", 24 * time.Hour},
		{"parsed over max gets clamped", "7d", maxTTL},
		{"parsed weeks clamped", "2w", maxTTL},
		{"overflowing tag clamps to max", "20000w", maxTTL},
		{"parsed seconds untouched", "30s", 30 * time.Second},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ComputeExpiry(tc.reference, now, defaultTTL, maxTTL)
			want := now.Add(tc.want)
			if !got.Equal(want) {
				t.Fatalf("got=%v want=%v", got, want)
			}
		})
	}
}

func TestComputeExpiryNonDefaultMax(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	// Default larger than max: default itself should be clamped.
	got := ComputeExpiry("latest", now, 48*time.Hour, 12*time.Hour)
	want := now.Add(12 * time.Hour)
	if !got.Equal(want) {
		t.Fatalf("default-clamping: got=%v want=%v", got, want)
	}
}
