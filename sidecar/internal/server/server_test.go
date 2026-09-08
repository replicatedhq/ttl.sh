package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/replicatedhq/ttl.sh/sidecar/internal/events"
	"github.com/replicatedhq/ttl.sh/sidecar/internal/store"
)

// Server.Store is documented as safe for concurrent use, so the fake locks its
// own state rather than leaning on the caller to serialize.
type fakeStore struct {
	mu        sync.Mutex
	rows      []store.Row
	upsertErr error // if set, Upsert returns it
}

func (f *fakeStore) Upsert(repo, tag, digest string, expiresAt, createdAt time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.upsertErr != nil {
		return f.upsertErr
	}
	f.rows = append(f.rows, store.Row{
		Repository:     repo,
		Tag:            tag,
		ManifestDigest: digest,
		ExpiresAt:      expiresAt.Unix(),
		CreatedAt:      createdAt.Unix(),
	})
	return nil
}

const testEpoch = 1_700_000_000

func fixedClock() time.Time { return time.Unix(testEpoch, 0) }

// newTestServer returns a Server backed by an in-memory fake plus the fake
// itself so tests can assert on what was recorded.
func newTestServer() (*Server, *fakeStore) {
	fs := &fakeStore{}
	return New(fs, 24*time.Hour, 24*time.Hour, fixedClock), fs
}

func TestHandleEventsRejectsNonPost(t *testing.T) {
	srv, _ := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	rec := httptest.NewRecorder()
	srv.handleEvents(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got %d want 405", rec.Code)
	}
	if rec.Header().Get("Allow") != http.MethodPost {
		t.Errorf("Allow = %q", rec.Header().Get("Allow"))
	}
}

func TestHandleEventsImageUpdatedUpserts(t *testing.T) {
	srv, fs := newTestServer()
	body, _ := json.Marshal(events.ImageUpdatedData{
		Name: "foo/bar", Reference: "1h", Digest: "sha256:deadbeef",
	})
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Ce-Type", events.ImageUpdatedType)
	rec := httptest.NewRecorder()

	srv.handleEvents(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("got %d want 204", rec.Code)
	}
	if len(fs.rows) != 1 {
		t.Fatalf("got %d rows want 1", len(fs.rows))
	}
	r := fs.rows[0]
	if r.Repository != "foo/bar" || r.Tag != "1h" || r.ManifestDigest != "sha256:deadbeef" {
		t.Errorf("row = %+v", r)
	}
	wantExpires := fixedClock().Add(time.Hour).Unix()
	if r.ExpiresAt != wantExpires {
		t.Errorf("expires_at = %d want %d", r.ExpiresAt, wantExpires)
	}
}

// A multi-arch push emits an image.updated per child manifest, referenced by
// digest, before the one carrying the index's tag. Only the tag is trackable:
// zot answers 405/DENIED to a delete of an index's child, so recording the
// children would wedge the reaper on rows it can never clear.
func TestHandleEventsIgnoresDigestReference(t *testing.T) {
	srv, fs := newTestServer()
	digest := "sha256:d3d669c9a5ef6483b05164101265237d0fff3a6495f659242262a1d8d68e2dda"
	incoming := []events.ImageUpdatedData{
		{Name: "foo/bar", Reference: digest, Digest: digest},             // child manifest
		{Name: "foo/bar", Reference: "sha256:abc", Digest: "sha256:abc"}, // child manifest
		{Name: "foo/bar", Reference: "1h", Digest: "sha256:index"},       // the index's tag
	}
	for _, data := range incoming {
		body, _ := json.Marshal(data)
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
		req.Header.Set("Ce-Type", events.ImageUpdatedType)
		rec := httptest.NewRecorder()
		srv.handleEvents(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("reference %q: got %d want 204", data.Reference, rec.Code)
		}
	}

	if len(fs.rows) != 1 {
		t.Fatalf("rows = %+v, want only the tagged index", fs.rows)
	}
	if fs.rows[0].Tag != "1h" {
		t.Errorf("tracked tag = %q, want %q", fs.rows[0].Tag, "1h")
	}
}

func TestHandleEventsOtherTypeAcked(t *testing.T) {
	srv, fs := newTestServer()
	body := []byte(`{"name":"foo","reference":"1h","digest":"sha256:x"}`)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Ce-Type", "zotregistry.image.deleted")
	rec := httptest.NewRecorder()

	srv.handleEvents(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("got %d want 204", rec.Code)
	}
	if len(fs.rows) != 0 {
		t.Fatalf("non-update event should not store; got %d rows", len(fs.rows))
	}
}

func TestHandleEventsMissingFieldsAcked(t *testing.T) {
	srv, fs := newTestServer()
	body := []byte(`{"digest":"sha256:x"}`) // no name, no reference
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Ce-Type", events.ImageUpdatedType)
	rec := httptest.NewRecorder()

	srv.handleEvents(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("got %d want 204", rec.Code)
	}
	if len(fs.rows) != 0 {
		t.Fatalf("missing fields should not store; got %d rows", len(fs.rows))
	}
}

func TestHandleEventsMalformedReturns400(t *testing.T) {
	srv, _ := newTestServer()
	body := []byte(`{not json`)
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Ce-Type", events.ImageUpdatedType)
	rec := httptest.NewRecorder()

	srv.handleEvents(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d want 400", rec.Code)
	}
}

func TestHandleEventsStoreErrorReturns500(t *testing.T) {
	fs := &fakeStore{upsertErr: errors.New("boom")}
	srv := New(fs, 24*time.Hour, 24*time.Hour, fixedClock)
	body, _ := json.Marshal(events.ImageUpdatedData{
		Name: "foo", Reference: "1h", Digest: "sha256:x",
	})
	req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
	req.Header.Set("Ce-Type", events.ImageUpdatedType)
	rec := httptest.NewRecorder()

	srv.handleEvents(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("got %d want 500", rec.Code)
	}
}

func TestHandleHealthz(t *testing.T) {
	srv, _ := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	srv.handleHealthz(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("got %d want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "ok") {
		t.Errorf("body = %q", rec.Body.String())
	}
}

func TestShortDigest(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", ""},
		{"sha256:abc", "sha256:abc"}, // shorter than threshold
		{"sha256:0123456789ab", "sha256:0123456789ab"},     // exactly at threshold (len 19)
		{"sha256:0123456789abcdef", "sha256:0123456789ab"}, // longer; truncated
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got := shortDigest(tc.in)
			if got != tc.want {
				t.Errorf("got %q want %q", got, tc.want)
			}
		})
	}
}

// TestRoutesWiring exercises New + Routes over real HTTP round-trips.
func TestRoutesWiring(t *testing.T) {
	srv := New(&fakeStore{}, 24*time.Hour, 24*time.Hour, fixedClock)
	ts := httptest.NewServer(srv.Routes())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET /healthz = %d want 200", resp.StatusCode)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	body, _ := json.Marshal(events.ImageUpdatedData{
		Name: "foo/bar", Reference: "1h", Digest: "sha256:deadbeef",
	})
	req, err := http.NewRequest(http.MethodPost, ts.URL+"/events", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Ce-Type", events.ImageUpdatedType)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST /events: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("POST /events = %d want 204", resp.StatusCode)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()
}
