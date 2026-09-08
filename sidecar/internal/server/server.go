// Package server implements the HTTP surface for the TTL service: a CloudEvents
// sink for zot image events and a health check.
package server

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/replicatedhq/ttl.sh/sidecar/internal/events"
	"github.com/replicatedhq/ttl.sh/sidecar/internal/ttl"
)

// Store is the subset of the persistence layer the server needs.
// Implementations must be safe for concurrent use: the HTTP handlers call it
// from whatever goroutine net/http hands them, without serializing first.
type Store interface {
	Upsert(repo, tag, digest string, expiresAt, createdAt time.Time) error
}

// Server is the HTTP server: events sink + healthz. It holds no lock of its
// own; concurrent event POSTs are safe by way of the Store contract above.
type Server struct {
	store      Store
	defaultTTL time.Duration
	maxTTL     time.Duration
	clock      func() time.Time
}

// New constructs a Server. defaultTTL and maxTTL drive the expiry policy applied
// to incoming events.
func New(st Store, defaultTTL, maxTTL time.Duration, clock func() time.Time) *Server {
	return &Server{
		store:      st,
		defaultTTL: defaultTTL,
		maxTTL:     maxTTL,
		clock:      clock,
	}
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		http.Error(w, "read body: "+err.Error(), http.StatusBadRequest)
		return
	}
	evt, err := events.Parse(r, body)
	if err != nil {
		log.Printf("events: parse error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if evt.Type != events.ImageUpdatedType {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if evt.Data.Name == "" || evt.Data.Reference == "" {
		log.Printf("events: image.updated missing name/reference; ignoring (digest=%q)", evt.Data.Digest)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// A multi-arch push emits one event per child manifest, referenced by
	// digest, before the event for the index's tag. zot refuses to delete a
	// manifest an index still references, so tracking those children would only
	// produce rows the reaper can never clear; they go away with their index.
	if isDigest(evt.Data.Reference) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	now := s.clock()
	expires := ttl.ComputeExpiry(evt.Data.Reference, now, s.defaultTTL, s.maxTTL)

	err = s.store.Upsert(evt.Data.Name, evt.Data.Reference, evt.Data.Digest, expires, now)
	if err != nil {
		log.Printf("events: upsert %s:%s: %v", evt.Data.Name, evt.Data.Reference, err)
		http.Error(w, "store error", http.StatusInternalServerError)
		return
	}
	log.Printf("events: recorded %s:%s digest=%s expires_at=%s",
		evt.Data.Name, evt.Data.Reference, shortDigest(evt.Data.Digest), expires.UTC().Format(time.RFC3339))
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ok\n")
}

// Routes returns the HTTP handler wiring the server's endpoints.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/events", s.handleEvents)
	mux.HandleFunc("/healthz", s.handleHealthz)
	return mux
}

// isDigest reports whether reference is a digest rather than a tag. An OCI tag
// cannot contain a colon, so the separator alone is decisive.
func isDigest(reference string) bool {
	return strings.Contains(reference, ":")
}

func shortDigest(d string) string {
	if len(d) > 19 { // "sha256:" + 12 hex chars
		return d[:19]
	}
	return d
}
