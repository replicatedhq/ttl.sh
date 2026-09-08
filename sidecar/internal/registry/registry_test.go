package registry

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestDeleteManifestStatus checks the status-to-error mapping against a fake
// zot, and that the request method and path are what the API expects.
func TestDeleteManifestStatus(t *testing.T) {
	cases := []struct {
		status  int
		wantErr bool
	}{
		{http.StatusOK, false},
		{http.StatusAccepted, false},
		{http.StatusNoContent, false},
		{http.StatusNotFound, false},
		{http.StatusBadRequest, true},
		{http.StatusInternalServerError, true},
		{http.StatusServiceUnavailable, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(http.StatusText(tc.status), func(t *testing.T) {
			var gotMethod, gotPath string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotMethod = r.Method
				gotPath = r.URL.Path
				w.WriteHeader(tc.status)
			}))
			defer srv.Close()

			c := New(srv.URL)
			err := c.DeleteManifest(context.Background(), "foo/bar", "v1")

			if tc.wantErr && err == nil {
				t.Fatalf("status %d: expected error, got nil", tc.status)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("status %d: expected nil, got error: %v", tc.status, err)
			}

			if gotMethod != http.MethodDelete {
				t.Errorf("method = %q, want %q", gotMethod, http.MethodDelete)
			}
			if want := "/v2/foo/bar/manifests/v1"; gotPath != want {
				t.Errorf("path = %q, want %q", gotPath, want)
			}
		})
	}
}

// zot answers 405/DENIED when the reference is a manifest an index still
// points at. That has to arrive as ErrManifestReferenced, not as a generic
// error, or the reaper would retry it every tick forever.
func TestDeleteManifestReferencedByIndex(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"errors":[{"code":"DENIED","message":"requested access to the resource is denied"}]}`))
	}))
	defer srv.Close()

	err := New(srv.URL).DeleteManifest(context.Background(), "foo/bar", "sha256:abc")
	if !errors.Is(err, ErrManifestReferenced) {
		t.Fatalf("err = %v, want ErrManifestReferenced", err)
	}
}

// Once a repo's last tag is reaped the repo itself disappears, and zot answers a
// delete against it with 400/NAME_UNKNOWN rather than 404. A cosign .sig tag
// outliving its image lands here, so it has to read as success — otherwise the
// row is never dropped and the reaper retries it every tick forever.
func TestDeleteManifestRepoAlreadyGone(t *testing.T) {
	cases := []struct {
		name    string
		status  int
		body    string
		wantErr bool
	}{
		{
			name:   "NAME_UNKNOWN",
			status: http.StatusBadRequest,
			body:   `{"errors":[{"code":"NAME_UNKNOWN","message":"repository name not known to registry","detail":{"name":"c1e2d663-74a5-47d1-a395-d1b8cdacb137"}}]}`,
		},
		{
			name:   "MANIFEST_UNKNOWN",
			status: http.StatusBadRequest,
			body:   `{"errors":[{"code":"MANIFEST_UNKNOWN","message":"manifest unknown"}]}`,
		},
		{
			// A malformed request is also a 400, and that one is worth retrying
			// and surfacing rather than silently untracking.
			name:    "other 400 stays an error",
			status:  http.StatusBadRequest,
			body:    `{"errors":[{"code":"UNSUPPORTED","message":"the operation is unsupported"}]}`,
			wantErr: true,
		},
		{
			name:    "unparseable body stays an error",
			status:  http.StatusBadRequest,
			body:    `not json`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer srv.Close()

			err := New(srv.URL).DeleteManifest(context.Background(), "foo/bar", "sha256-abc.sig")
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected nil, got error: %v", err)
			}
		})
	}
}

// Port 1 refuses connections, so http.Client.Do fails and the error surfaces.
func TestDeleteManifestTransportError(t *testing.T) {
	c := New("http://127.0.0.1:1")
	if err := c.DeleteManifest(context.Background(), "foo/bar", "v1"); err == nil {
		t.Fatal("expected transport error, got nil")
	}
}

// TestDeleteManifestRequestBuildError uses a base URL containing a control
// character so that http.NewRequest fails before any network I/O.
func TestDeleteManifestRequestBuildError(t *testing.T) {
	// Sanity-check that this base URL really makes http.NewRequest fail, so
	// the test exercises the request-build error path and not something else.
	badBase := "http://\x7f-bad-host"
	if _, err := http.NewRequest(http.MethodDelete, badBase+"/v2/foo/bar/manifests/v1", nil); err == nil {
		t.Fatalf("precondition failed: expected http.NewRequest to fail for %q", badBase)
	}

	c := New(badBase)
	if err := c.DeleteManifest(context.Background(), "foo/bar", "v1"); err == nil {
		t.Fatal("expected request-build error, got nil")
	}
}
