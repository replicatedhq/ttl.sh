// Package reaper implements the background sweep loop that periodically queries
// the store for expired tags and asks the registry to remove their manifests. A
// row is removed from the store only after the registry confirms the manifest
// is gone, so a failed delete is retried on the next tick.
package reaper

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/replicatedhq/ttl.sh/sidecar/internal/registry"
	"github.com/replicatedhq/ttl.sh/sidecar/internal/store"
)

// Store is the subset of the persistence layer the reaper needs.
type Store interface {
	Expired(now time.Time) ([]store.Row, error)
	Delete(repo, tag string) error
}

// ManifestDeleter removes a manifest by tag from the registry.
type ManifestDeleter interface {
	DeleteManifest(ctx context.Context, repo, tag string) error
}

type Reaper struct {
	interval time.Duration
	store    Store
	registry ManifestDeleter
	clock    func() time.Time
}

// New constructs a Reaper with clock set to time.Now.
func New(interval time.Duration, st Store, reg ManifestDeleter) *Reaper {
	return &Reaper{
		interval: interval,
		store:    st,
		registry: reg,
		clock:    time.Now,
	}
}

// Run sweeps every interval until ctx is cancelled.
func (r *Reaper) Run(ctx context.Context) {
	t := time.NewTicker(r.interval)
	defer t.Stop()

	// Sweep once on start so a restart doesn't wait a full interval to reap
	// already-expired rows.
	r.sweepOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			r.sweepOnce(ctx)
		}
	}
}

func (r *Reaper) sweepOnce(ctx context.Context) {
	now := r.clock()
	rows, err := r.store.Expired(now)
	if err != nil {
		log.Printf("sweep: query expired: %v", err)
		return
	}
	if len(rows) == 0 {
		return
	}
	log.Printf("sweep: %d expired tag(s)", len(rows))
	for _, row := range rows {
		switch err := r.registry.DeleteManifest(ctx, row.Repository, row.Tag); {
		case err == nil:
			log.Printf("sweep: deleted %s:%s", row.Repository, row.Tag)
		case errors.Is(err, registry.ErrManifestReferenced):
			// Retrying can never succeed, so drop the row instead of failing
			// this pair on every tick. The index that holds this manifest has a
			// row of its own; once that expires, zot's untagged retention
			// collects the child.
			log.Printf("sweep: %s:%s held by an index, untracking: %v", row.Repository, row.Tag, err)
		default:
			log.Printf("sweep: delete %s:%s: %v", row.Repository, row.Tag, err)
			continue
		}
		if err := r.store.Delete(row.Repository, row.Tag); err != nil {
			log.Printf("sweep: row delete %s:%s: %v", row.Repository, row.Tag, err)
		}
	}
}
