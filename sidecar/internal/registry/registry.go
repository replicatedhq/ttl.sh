// Package registry is the adapter to zot's OCI distribution API. It exposes the
// single operation the reaper needs: delete a manifest by tag.
package registry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

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

// DeleteManifest issues DELETE /v2/<repo>/manifests/<tag>; zot resolves the tag
// to its digest. 200/202/204/404 all mean the tag is gone; any other status is
// a transient error to retry on the next tick. ctx bounds the request, so a
// wedged registry cannot outlive a shutdown.
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
	switch resp.StatusCode {
	case http.StatusOK, http.StatusAccepted, http.StatusNoContent, http.StatusNotFound:
		return nil
	default:
		// zot describes the failure in the body; %q keeps it on one log line.
		return fmt.Errorf("DELETE %s -> %d: %q", endpoint, resp.StatusCode, body)
	}
}
