// Package store is the Redis persistence layer: it records
// (repository, tag, manifest_digest, expires_at) rows and answers which
// of them have expired.
//
// State lives in exactly two keys:
//
//	zot-ephemeral-ttl:rows   HASH  field = <repository>\x00<tag>, value = Row as JSON
//	zot-ephemeral-ttl:index  ZSET  member = <repository>\x00<tag>, score = expires_at
//
// The sorted set is the expiry index: "which tags are due" is a range query by
// score against it. The NUL separator keeps (repository, tag) pairs unambiguous
// — no repository or tag can contain one — and makes the sorted set's
// lexicographic tie-break for equal scores fall out as (repository, tag) order.
//
// Rows deliberately do NOT carry a Redis key TTL. The reaper has to *see* an
// expired row in order to issue the manifest delete to zot, and it removes the
// row only once zot confirms; letting Redis expire the data itself would drop
// tags on the floor without ever deleting their manifests.
package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// The two keys this service owns. The Redis instance is dedicated to the
// sidecar, so they are fixed rather than configurable.
const (
	rowsKey  = "zot-ephemeral-ttl:rows"
	indexKey = "zot-ephemeral-ttl:index"
)

// opTimeout bounds every individual Redis round trip. The consumer interfaces
// take no context, so each operation carries its own deadline; a wedged Redis
// therefore surfaces as an error rather than a stuck sweep or event handler.
const opTimeout = 5 * time.Second

// memberSep joins repository and tag into a single sorted-set member / hash
// field. NUL cannot appear in either, so the encoding is injective.
const memberSep = "\x00"

// Row is a single tracked tag and its expiry.
type Row struct {
	Repository     string `json:"repository"`
	Tag            string `json:"tag"`
	ManifestDigest string `json:"manifest_digest"`
	ExpiresAt      int64  `json:"expires_at"`
	CreatedAt      int64  `json:"created_at"`
}

type Store struct {
	rdb *redis.Client
}

// Open dials the Redis server named by rawURL (a redis:// or rediss:// URL as
// understood by redis.ParseURL) and verifies the connection with a PING, so a
// bad address or credential fails at startup rather than on the first push.
func Open(rawURL string) (*Store, error) {
	opt, err := redis.ParseURL(rawURL)
	if err != nil {
		// Deliberately not echoing rawURL: it may carry a password.
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	rdb := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &Store{rdb: rdb}, nil
}

func (s *Store) Close() error { return s.rdb.Close() }

func (s *Store) opCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), opTimeout)
}

func member(repo, tag string) string { return repo + memberSep + tag }

// Upsert inserts or refreshes the row for (repo, tag). Timestamps are truncated
// to whole seconds. The row payload and its index entry are written in one
// transaction so a reader never sees an index entry without its row.
func (s *Store) Upsert(repo, tag, digest string, expiresAt, createdAt time.Time) error {
	row := Row{
		Repository:     repo,
		Tag:            tag,
		ManifestDigest: digest,
		ExpiresAt:      expiresAt.Unix(),
		CreatedAt:      createdAt.Unix(),
	}
	payload, err := json.Marshal(row)
	if err != nil {
		return err
	}

	ctx, cancel := s.opCtx()
	defer cancel()
	m := member(repo, tag)
	_, err = s.rdb.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.HSet(ctx, rowsKey, m, payload)
		p.ZAdd(ctx, indexKey, redis.Z{Score: float64(row.ExpiresAt), Member: m})
		return nil
	})
	return err
}

// Expired returns every row whose expiry is at or before now, ordered by expiry
// ascending with ties broken by (repository, tag) — the sorted set orders equal
// scores lexicographically by member, which is exactly that pair.
func (s *Store) Expired(now time.Time) ([]Row, error) {
	ctx, cancel := s.opCtx()
	defer cancel()

	members, err := s.rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:     indexKey,
		Start:   "-inf",
		Stop:    strconv.FormatInt(now.Unix(), 10), // inclusive: expires_at <= now
		ByScore: true,
	}).Result()
	if err != nil {
		return nil, err
	}
	return s.rows(ctx, members)
}

// Delete removes the row for (repo, tag); a missing key is a no-op.
func (s *Store) Delete(repo, tag string) error {
	ctx, cancel := s.opCtx()
	defer cancel()

	m := member(repo, tag)
	_, err := s.rdb.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.HDel(ctx, rowsKey, m)
		p.ZRem(ctx, indexKey, m)
		return nil
	})
	return err
}

// rows fetches the row payloads for members, preserving their order. An index
// entry whose payload is missing is dropped from the index (best effort) and
// skipped: the pair is only ever written in one transaction, so a gap means
// something outside this service removed the hash field.
func (s *Store) rows(ctx context.Context, members []string) ([]Row, error) {
	if len(members) == 0 {
		return nil, nil
	}
	vals, err := s.rdb.HMGet(ctx, rowsKey, members...).Result()
	if err != nil {
		return nil, err
	}

	var out []Row
	var orphans []string
	for i, v := range vals {
		payload, ok := v.(string)
		if !ok { // nil: no such hash field
			orphans = append(orphans, members[i])
			continue
		}
		var r Row
		if err := json.Unmarshal([]byte(payload), &r); err != nil {
			return nil, fmt.Errorf("decode row %q: %w", members[i], err)
		}
		out = append(out, r)
	}
	if len(orphans) > 0 {
		args := make([]any, len(orphans))
		for i, m := range orphans {
			args[i] = m
		}
		_ = s.rdb.ZRem(ctx, indexKey, args...).Err()
	}
	return out, nil
}
