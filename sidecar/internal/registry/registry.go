// Package registry is the adapter to zot's OCI distribution API. It exposes the
// single operation the reaper needs: delete a manifest by tag.
package registry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"
)

// ErrManifestReferenced is zot's 405/DENIED: the reference names a manifest an
// index still points at. Despite the "access denied" wording it is not an
// authorization failure, and it will not clear on retry — only deleting the
// index releases the child.
var ErrManifestReferenced = errors.New("manifest is referenced by an index")

type Client struct {
	baseURL string
	http    *http.Client
}

// New returns a Client targeting the given zot base URL (e.g.
// "http://zot:5000"). The base URL should not have a trailing slash.
func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// gone lists the distribution error codes that mean the delete has nothing left
// to do: the repository or the manifest is already absent. zot returns
// NAME_UNKNOWN with a 400, not a 404, once a repo's last tag is reaped and the
// repo itself disappears — a cosign .sig tag outliving its image hits this.
var gone = []string{"NAME_UNKNOWN", "MANIFEST_UNKNOWN"}

// isGone reports whether an error body carries one of the `gone` codes. Matching
// on the body rather than the status keeps a genuinely malformed request, which
// is also a 400, retryable and visible.
func isGone(body []byte) bool {
	var payload struct {
		Errors []struct {
			Code string `json:"code"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return false
	}
	for _, e := range payload.Errors {
		if slices.Contains(gone, e.Code) {
			return true
		}
	}
	return false
}

// DeleteManifest issues DELETE /v2/<repo>/manifests/<tag>; zot resolves the tag
// to its digest. 200/202/204/404, and any error naming a missing repo or
// manifest, all mean the tag is gone; any other status is a transient error to
// retry on the next tick. ctx bounds the request, so a wedged registry cannot
// outlive a shutdown.
func (c *Client) DeleteManifest(ctx context.Context, repo, tag string) error {
	endpoint := fmt.Sprintf("%s/v2/%s/manifests/%s", c.baseURL, repo, tag)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	body, _ := io.ReadAll(resp.Body)
	switch {
	case resp.StatusCode == http.StatusOK,
		resp.StatusCode == http.StatusAccepted,
		resp.StatusCode == http.StatusNoContent,
		resp.StatusCode == http.StatusNotFound,
		isGone(body):
		return nil
	case resp.StatusCode == http.StatusMethodNotAllowed:
		return fmt.Errorf("DELETE %s: %w", endpoint, ErrManifestReferenced)
	default:
		// zot describes the failure in the body; %q keeps it on one log line.
		return fmt.Errorf("DELETE %s -> %d: %q", endpoint, resp.StatusCode, body)
	}
}
