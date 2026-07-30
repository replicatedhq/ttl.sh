// Package ttl parses TTL-encoding tag references (e.g. "1h", "30m", "7d") and
// applies the expiry policy: the parsed TTL, or a default when the tag doesn't
// match, clamped to a maximum.
package ttl

import (
	"math"
	"regexp"
	"strconv"
	"time"
)

var tagRe = regexp.MustCompile(`^([0-9]+)(s|m|h|d|w)$`)

// ParseTag returns the TTL encoded in reference and whether it matched the
// `^\d+(s|m|h|d|w)$` form.
func ParseTag(reference string) (time.Duration, bool) {
	m := tagRe.FindStringSubmatch(reference)
	if m == nil {
		return 0, false
	}
	n, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil || n <= 0 {
		return 0, false
	}
	var unit time.Duration
	switch m[2] {
	case "s":
		unit = time.Second
	case "m":
		unit = time.Minute
	case "h":
		unit = time.Hour
	case "d":
		unit = 24 * time.Hour
	case "w":
		unit = 7 * 24 * time.Hour
	default:
		return 0, false
	}
	// A tag can name a period far longer than time.Duration can hold. Saturate
	// instead of letting the multiply wrap: a negative duration would slip past
	// the clamp in ComputeExpiry and expire the tag immediately, which is the
	// opposite of what a huge TTL asks for.
	if n > int64(math.MaxInt64)/int64(unit) {
		return time.Duration(math.MaxInt64), true
	}
	return time.Duration(n) * unit, true
}

// ComputeExpiry applies the policy: the parsed TTL (or defaultTTL if reference
// doesn't encode one), clamped to maxTTL.
func ComputeExpiry(reference string, now time.Time, defaultTTL, maxTTL time.Duration) time.Time {
	ttl := defaultTTL
	if parsed, ok := ParseTag(reference); ok {
		ttl = parsed
	}
	if ttl > maxTTL {
		ttl = maxTTL
	}
	return now.Add(ttl)
}
