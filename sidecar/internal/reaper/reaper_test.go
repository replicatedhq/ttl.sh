package reaper

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/replicatedhq/ttl.sh/sidecar/internal/store"
)

// In-memory fakes for the reaper's two dependencies. No real DB, no real HTTP.
type fakeStore struct {
	mu         sync.Mutex
	expired    []store.Row
	expiredErr error
	deleted    [][2]string // recorded (repo,tag)
	deleteErr  error
}

func (f *fakeStore) Expired(now time.Time) ([]store.Row, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.expired, f.expiredErr
}

func (f *fakeStore) Delete(repo, tag string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.deleted = append(f.deleted, [2]string{repo, tag})
	return f.deleteErr
}

func (f *fakeStore) deletedCalls() [][2]string {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([][2]string, len(f.deleted))
	copy(out, f.deleted)
	return out
}

type fakeDeleter struct {
	mu     sync.Mutex
	calls  [][2]string
	failOn map[string]bool // key "repo:tag" -> return error
}

func (d *fakeDeleter) DeleteManifest(_ context.Context, repo, tag string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.calls = append(d.calls, [2]string{repo, tag})
	if d.failOn[repo+":"+tag] {
		return errors.New("boom")
	}
	return nil
}

func (d *fakeDeleter) callList() [][2]string {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([][2]string, len(d.calls))
	copy(out, d.calls)
	return out
}

func row(repo, tag string) store.Row {
	return store.Row{
		Repository:     repo,
		Tag:            tag,
		ManifestDigest: "sha256:" + repo + tag,
		ExpiresAt:      1,
		CreatedAt:      0,
	}
}

func TestSweepOnceDeletesExpiredRow(t *testing.T) {
	fs := &fakeStore{expired: []store.Row{row("r", "old")}}
	fd := &fakeDeleter{}
	rp := New(time.Hour, fs, fd)

	rp.sweepOnce(context.Background())

	if got := fd.callList(); len(got) != 1 || got[0] != [2]string{"r", "old"} {
		t.Fatalf("DeleteManifest calls = %v, want one [r old]", got)
	}
	if got := fs.deletedCalls(); len(got) != 1 || got[0] != [2]string{"r", "old"} {
		t.Fatalf("store.Delete calls = %v, want one [r old]", got)
	}
}

func TestSweepOnceExpiredError(t *testing.T) {
	fs := &fakeStore{expiredErr: errors.New("db down")}
	fd := &fakeDeleter{}
	rp := New(time.Hour, fs, fd)

	rp.sweepOnce(context.Background())

	if got := fd.callList(); len(got) != 0 {
		t.Fatalf("DeleteManifest calls = %v, want none", got)
	}
	if got := fs.deletedCalls(); len(got) != 0 {
		t.Fatalf("store.Delete calls = %v, want none", got)
	}
}

func TestSweepOnceNoExpiredRows(t *testing.T) {
	fs := &fakeStore{expired: nil}
	fd := &fakeDeleter{}
	rp := New(time.Hour, fs, fd)

	rp.sweepOnce(context.Background())

	if got := fd.callList(); len(got) != 0 {
		t.Fatalf("DeleteManifest calls = %v, want none", got)
	}
	if got := fs.deletedCalls(); len(got) != 0 {
		t.Fatalf("store.Delete calls = %v, want none", got)
	}
}

func TestSweepOnceDeleterErrorPreservesRow(t *testing.T) {
	fs := &fakeStore{expired: []store.Row{row("r", "stuck")}}
	fd := &fakeDeleter{failOn: map[string]bool{"r:stuck": true}}
	rp := New(time.Hour, fs, fd)

	rp.sweepOnce(context.Background())

	if got := fd.callList(); len(got) != 1 || got[0] != [2]string{"r", "stuck"} {
		t.Fatalf("DeleteManifest calls = %v, want one [r stuck]", got)
	}
	if got := fs.deletedCalls(); len(got) != 0 {
		t.Fatalf("store.Delete calls = %v, want none (row preserved)", got)
	}
}

func TestSweepOnceStoreDeleteError(t *testing.T) {
	fs := &fakeStore{
		expired:   []store.Row{row("r", "old")},
		deleteErr: errors.New("delete failed"),
	}
	fd := &fakeDeleter{}
	rp := New(time.Hour, fs, fd)

	rp.sweepOnce(context.Background())

	if got := fd.callList(); len(got) != 1 || got[0] != [2]string{"r", "old"} {
		t.Fatalf("DeleteManifest calls = %v, want one [r old]", got)
	}
	if got := fs.deletedCalls(); len(got) != 1 || got[0] != [2]string{"r", "old"} {
		t.Fatalf("store.Delete calls = %v, want one [r old]", got)
	}
}

// A failure on one row must not stop the sweep: both manifest deletes are
// attempted, and only the row that succeeded is removed from the store.
func TestSweepOnceMixedResults(t *testing.T) {
	fs := &fakeStore{expired: []store.Row{row("r", "bad"), row("r", "good")}}
	fd := &fakeDeleter{failOn: map[string]bool{"r:bad": true}}
	rp := New(time.Hour, fs, fd)

	rp.sweepOnce(context.Background())

	gotCalls := fd.callList()
	if len(gotCalls) != 2 || gotCalls[0] != [2]string{"r", "bad"} || gotCalls[1] != [2]string{"r", "good"} {
		t.Fatalf("DeleteManifest calls = %v, want [r bad] then [r good]", gotCalls)
	}
	gotDel := fs.deletedCalls()
	if len(gotDel) != 1 || gotDel[0] != [2]string{"r", "good"} {
		t.Fatalf("store.Delete calls = %v, want only [r good]", gotDel)
	}
}

func TestRunImmediateSweepAndStopsOnContext(t *testing.T) {
	fs := &fakeStore{expired: []store.Row{row("r", "old")}}
	fd := &fakeDeleter{}
	rp := New(5*time.Millisecond, fs, fd)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		rp.Run(ctx)
		close(done)
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not return after cancel")
	}

	if got := fd.callList(); len(got) < 1 {
		t.Fatalf("expected at least one DeleteManifest call, got %d", len(got))
	}
	if got := fs.deletedCalls(); len(got) < 1 {
		t.Fatalf("expected at least one store.Delete call, got %d", len(got))
	}
}
