package registry

import (
	"context"
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
