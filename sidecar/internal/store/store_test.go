package store

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/replicatedhq/ttl.sh/sidecar/internal/redistest"
)

// These tests run against a real Redis, started for the package by redistest.
func TestMain(m *testing.M) { os.Exit(redistest.Run(m)) }

// farFuture is later than any expiry these tests set, so Expired(farFuture)
// returns the whole keyspace — the read-back used to inspect stored state.
var farFuture = time.Unix(1<<40, 0)

// newTestStore opens a store against that server. The keys are fixed and the
// whole package shares one container, so each test starts from a clean slate
// and drops them again on the way out.
func newTestStore(t *testing.T) *Store {
	t.Helper()
	st, err := Open(redistest.URL(t))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	reset := func() {
		ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
		defer cancel()
		if err := st.rdb.Del(ctx, rowsKey, indexKey).Err(); err != nil {
			t.Errorf("clearing test keys: %v", err)
		}
	}
	reset()
	t.Cleanup(func() {
		reset()
		if err := st.Close(); err != nil {
			t.Errorf("cleanup close: %v", err)
		}
	})
	return st
}

func TestOpenRejectsBadURL(t *testing.T) {
	// Not a redis:// URL: ParseURL fails before any dial, so this needs no
	// server.
	if _, err := Open("http://127.0.0.1:6379"); err == nil {
		t.Fatal("Open with a non-redis scheme: expected error, got nil")
	}
}

func TestOpenPingFails(t *testing.T) {
	// Port 1 refuses connections, so the startup PING fails and Open reports
	// it rather than returning a store that only breaks on first use.
	if _, err := Open("redis://127.0.0.1:1"); err == nil {
		t.Fatal("Open against a dead address: expected error, got nil")
	}
}

func TestOpenEmptyStore(t *testing.T) {
	st := newTestStore(t)
	rows, err := st.Expired(farFuture)
	if err != nil {
		t.Fatalf("Expired on empty store: %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("expected empty store, got %d rows", len(rows))
	}
}

// TestKeyLayout pins the documented storage layout: state lives in exactly two
// keys, with (repository, tag) encoded as the hash field.
func TestKeyLayout(t *testing.T) {
	st := newTestStore(t)
	if rowsKey != "zot-ephemeral-ttl:rows" || indexKey != "zot-ephemeral-ttl:index" {
		t.Errorf("keys = %q/%q, want zot-ephemeral-ttl:rows / :index", rowsKey, indexKey)
	}

	now := time.Now()
	if err := st.Upsert("repo/a", "v1", "sha256:aaa", now.Add(time.Hour), now); err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	keys, err := st.rdb.Keys(ctx, "zot-ephemeral-ttl*").Result()
	if err != nil {
		t.Fatalf("Keys: %v", err)
	}
	if len(keys) != 2 {
		t.Errorf("got keys %v, want exactly the rows hash and the index zset", keys)
	}

	// Field/member encoding is <repository>NUL<tag>.
	fields, err := st.rdb.HKeys(ctx, rowsKey).Result()
	if err != nil {
		t.Fatalf("HKeys: %v", err)
	}
	if len(fields) != 1 || fields[0] != "repo/a\x00v1" {
		t.Errorf("hash fields = %q, want [\"repo/a\\x00v1\"]", fields)
	}
}

func TestUpsertInsertThenUpdate(t *testing.T) {
	st := newTestStore(t)
	now := time.Now()

	if err := st.Upsert("repo/a", "v1", "sha256:aaa", now.Add(time.Hour), now); err != nil {
		t.Fatalf("first upsert: %v", err)
	}
	rows, err := st.Expired(farFuture)
	if err != nil {
		t.Fatalf("Expired: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("got %d rows want 1", len(rows))
	}
	if rows[0].ManifestDigest != "sha256:aaa" {
		t.Fatalf("digest = %q", rows[0].ManifestDigest)
	}

	// Same (repo, tag): should update digest, expires_at, created_at in place.
	later := now.Add(2 * time.Hour)
	if err := st.Upsert("repo/a", "v1", "sha256:bbb", later.Add(time.Hour), later); err != nil {
		t.Fatalf("second upsert: %v", err)
	}
	rows, err = st.Expired(farFuture)
	if err != nil {
		t.Fatalf("Expired: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("upsert duplicated row, got %d", len(rows))
	}
	if rows[0].ManifestDigest != "sha256:bbb" {
		t.Fatalf("digest not updated: %q", rows[0].ManifestDigest)
	}
	if rows[0].CreatedAt != later.Unix() {
		t.Fatalf("created_at not updated: got %d want %d", rows[0].CreatedAt, later.Unix())
	}
	if rows[0].ExpiresAt != later.Add(time.Hour).Unix() {
		t.Fatalf("expires_at not updated: got %d want %d", rows[0].ExpiresAt, later.Add(time.Hour).Unix())
	}
}

// TestUpsertRefreshesIndexScore proves the sorted-set score moves with the
// row: a re-push must actually postpone reaping, not just rewrite the payload.
func TestUpsertRefreshesIndexScore(t *testing.T) {
	st := newTestStore(t)
	now := time.Unix(1_700_000_000, 0)

	if err := st.Upsert("r", "t", "sha256:x", now.Add(-time.Hour), now.Add(-2*time.Hour)); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	expired, err := st.Expired(now)
	if err != nil {
		t.Fatalf("Expired: %v", err)
	}
	if len(expired) != 1 {
		t.Fatalf("row should be expired before refresh, got %d rows", len(expired))
	}

	if err := st.Upsert("r", "t", "sha256:x", now.Add(time.Hour), now); err != nil {
		t.Fatalf("refresh upsert: %v", err)
	}
	expired, err = st.Expired(now)
	if err != nil {
		t.Fatalf("Expired after refresh: %v", err)
	}
	if len(expired) != 0 {
		t.Fatalf("row should no longer be expired after refresh, got %d rows", len(expired))
	}
}

func TestExpiredFiltersByTime(t *testing.T) {
	st := newTestStore(t)
	now := time.Unix(1_700_000_000, 0)

	// Past, exactly-now, and future rows.
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(st.Upsert("r", "past", "sha256:1", now.Add(-time.Hour), now.Add(-2*time.Hour)))
	must(st.Upsert("r", "now", "sha256:2", now, now.Add(-time.Minute)))
	must(st.Upsert("r", "future", "sha256:3", now.Add(time.Hour), now))

	expired, err := st.Expired(now)
	if err != nil {
		t.Fatalf("Expired: %v", err)
	}
	tags := map[string]bool{}
	for _, r := range expired {
		tags[r.Tag] = true
	}
	if !tags["past"] {
		t.Errorf("past row should be expired")
	}
	if !tags["now"] {
		t.Errorf("now row should be expired (expires_at <= now)")
	}
	if tags["future"] {
		t.Errorf("future row should not be expired")
	}
}

// TestExpiredReturnsFullRow checks the payload survives the round trip: the
// reaper needs repository and tag to build the manifest DELETE.
func TestExpiredReturnsFullRow(t *testing.T) {
	st := newTestStore(t)
	now := time.Unix(1_700_000_000, 0)
	if err := st.Upsert("team/app", "30m", "sha256:abc", now.Add(-time.Second), now.Add(-time.Hour)); err != nil {
		t.Fatal(err)
	}
	rows, err := st.Expired(now)
	if err != nil {
		t.Fatalf("Expired: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("got %d rows want 1", len(rows))
	}
	want := Row{
		Repository:     "team/app",
		Tag:            "30m",
		ManifestDigest: "sha256:abc",
		ExpiresAt:      now.Add(-time.Second).Unix(),
		CreatedAt:      now.Add(-time.Hour).Unix(),
	}
	if rows[0] != want {
		t.Errorf("got %+v want %+v", rows[0], want)
	}
}

func TestDelete(t *testing.T) {
	st := newTestStore(t)
	now := time.Now()
	if err := st.Upsert("r", "t", "sha256:x", now.Add(time.Hour), now); err != nil {
		t.Fatal(err)
	}
	if err := st.Delete("r", "t"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	rows, err := st.Expired(farFuture)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("expected empty after delete, got %d", len(rows))
	}

	// Both keys must be cleaned up, not just the payload.
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	n, err := st.rdb.ZCard(ctx, indexKey).Result()
	if err != nil {
		t.Fatalf("ZCard: %v", err)
	}
	if n != 0 {
		t.Errorf("index still holds %d member(s) after delete", n)
	}

	// Deleting a non-existent row is a no-op (no error).
	if err := st.Delete("r", "missing"); err != nil {
		t.Fatalf("Delete missing: %v", err)
	}
}

// TestDeleteIsScopedToOneTag guards the (repository, tag) encoding: deleting
// one tag must not disturb its neighbours.
func TestDeleteIsScopedToOneTag(t *testing.T) {
	st := newTestStore(t)
	now := time.Now()
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(st.Upsert("r", "keep", "sha256:1", now.Add(time.Hour), now))
	must(st.Upsert("r", "drop", "sha256:2", now.Add(time.Hour), now))
	must(st.Upsert("r/nested", "keep", "sha256:3", now.Add(time.Hour), now))

	must(st.Delete("r", "drop"))

	rows, err := st.Expired(farFuture)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("got %d rows after delete, want 2", len(rows))
	}
	for _, r := range rows {
		if r.Tag == "drop" {
			t.Errorf("deleted row still present: %+v", r)
		}
	}
}

func TestOrdering(t *testing.T) {
	st := newTestStore(t)
	now := time.Now()
	// Insert in non-sorted order.
	rows := []struct {
		repo, tag string
		expires   time.Time
	}{
		{"r", "c", now.Add(3 * time.Hour)},
		{"r", "a", now.Add(1 * time.Hour)},
		{"r", "b", now.Add(2 * time.Hour)},
	}
	for _, r := range rows {
		if err := st.Upsert(r.repo, r.tag, "sha256:x", r.expires, now); err != nil {
			t.Fatal(err)
		}
	}
	got, err := st.Expired(farFuture)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Fatalf("got %d rows, want 3", len(got))
	}
	wantOrder := []string{"a", "b", "c"}
	for i, w := range wantOrder {
		if got[i].Tag != w {
			t.Errorf("position %d: got %q want %q", i, got[i].Tag, w)
		}
	}
}

// TestOrderingTieBreak pins the (repository, tag) tie-break for rows that
// share an expiry — it comes from the sorted set's lexicographic ordering of
// equal-score members, which the NUL separator makes match pair ordering.
func TestOrderingTieBreak(t *testing.T) {
	st := newTestStore(t)
	exp := time.Unix(1_700_000_000, 0)
	now := time.Now()
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	must(st.Upsert("repo/b", "x", "sha256:1", exp, now))
	must(st.Upsert("repo/a", "y", "sha256:2", exp, now))
	must(st.Upsert("repo/a", "x", "sha256:3", exp, now))

	got, err := st.Expired(farFuture)
	if err != nil {
		t.Fatal(err)
	}
	want := []struct{ repo, tag string }{
		{"repo/a", "x"},
		{"repo/a", "y"},
		{"repo/b", "x"},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d rows, want %d", len(got), len(want))
	}
	for i, w := range want {
		if got[i].Repository != w.repo || got[i].Tag != w.tag {
			t.Errorf("position %d: got %s:%s want %s:%s", i, got[i].Repository, got[i].Tag, w.repo, w.tag)
		}
	}
}

// TestOrphanedIndexEntryIsHealed covers an index member whose payload has gone
// missing (something outside this service touched the hash): it is skipped and
// dropped from the index rather than surfacing as a zero-valued row.
func TestOrphanedIndexEntryIsHealed(t *testing.T) {
	st := newTestStore(t)
	now := time.Now()
	if err := st.Upsert("r", "t", "sha256:x", now.Add(time.Hour), now); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	if err := st.rdb.HDel(ctx, rowsKey, member("r", "t")).Err(); err != nil {
		t.Fatalf("HDel: %v", err)
	}

	rows, err := st.Expired(farFuture)
	if err != nil {
		t.Fatalf("Expired with an orphaned index entry: %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("got %d rows, want the orphan skipped", len(rows))
	}
	n, err := st.rdb.ZCard(ctx, indexKey).Result()
	if err != nil {
		t.Fatalf("ZCard: %v", err)
	}
	if n != 0 {
		t.Errorf("orphaned index entry not cleaned up: %d member(s) remain", n)
	}
}

// TestCorruptPayloadErrors: an undecodable payload is real corruption, so it
// surfaces as an error instead of being silently dropped.
func TestCorruptPayloadErrors(t *testing.T) {
	st := newTestStore(t)
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	m := member("r", "t")
	if err := st.rdb.HSet(ctx, rowsKey, m, "not json").Err(); err != nil {
		t.Fatalf("HSet: %v", err)
	}
	if err := st.rdb.ZAdd(ctx, indexKey, redis.Z{Score: 1, Member: m}).Err(); err != nil {
		t.Fatalf("ZAdd: %v", err)
	}

	if _, err := st.Expired(farFuture); err == nil {
		t.Error("Expired with a corrupt payload: expected error, got nil")
	}
}

func TestOperationsOnClosedStore(t *testing.T) {
	st, err := Open(redistest.URL(t))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	now := time.Now()
	if err := st.Upsert("r", "t", "sha256:x", now.Add(time.Hour), now); err == nil {
		t.Error("Upsert on closed store: expected error, got nil")
	}
	if err := st.Delete("r", "t"); err == nil {
		t.Error("Delete on closed store: expected error, got nil")
	}
	if _, err := st.Expired(now); err == nil {
		t.Error("Expired on closed store: expected error, got nil")
	}
}

// TestConcurrentAccess runs concurrent writers and readers under the race
// detector: the server's event handler and the reaper's sweep share one store.
func TestConcurrentAccess(t *testing.T) {
	st := newTestStore(t)
	now := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 25; i++ {
		wg.Add(2)
		go func(n int) {
			defer wg.Done()
			tag := string(rune('a' + n%26))
			if err := st.Upsert("r", tag, "sha256:x", now.Add(time.Hour), now); err != nil {
				t.Errorf("concurrent upsert: %v", err)
			}
		}(i)
		go func() {
			defer wg.Done()
			if _, err := st.Expired(farFuture); err != nil {
				t.Errorf("concurrent Expired: %v", err)
			}
		}()
	}
	wg.Wait()
}
