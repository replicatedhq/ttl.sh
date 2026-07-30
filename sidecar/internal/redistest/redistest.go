// Package redistest supplies a real Redis to the tests that need one.
//
// The store exists to get Redis semantics right — transactional writes across
// two keys, score-range queries, lexicographic tie-breaks — so faking it would
// only assert what a reimplementation believes those semantics to be. A package
// opts in with a TestMain and gets a throwaway container for the run:
//
//	func TestMain(m *testing.M) { os.Exit(redistest.Run(m)) }
//
// One container is shared by all tests in the package, which are responsible
// for leaving the keyspace as they found it.
package redistest

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

// image pins the server under test. Keep it in step with the Redis the example
// stack runs.
const image = "redis:8-alpine"

// Exactly one of these is set before tests run, and neither changes after.
var (
	url         string
	unavailable string
)

// Run starts a Redis container, runs the package's tests against it, and
// terminates it. Pass the result to os.Exit from TestMain. A container that
// cannot be started is not fatal here; URL accounts for it per test.
func Run(m *testing.M) int {
	ctx := context.Background()

	container, err := tcredis.Run(ctx, image)
	if err != nil {
		unavailable = fmt.Sprintf("could not start %s: %v", image, err)
		return m.Run()
	}
	defer func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			fmt.Fprintf(os.Stderr, "redistest: terminating container: %v\n", err)
		}
	}()

	url, err = container.ConnectionString(ctx)
	if err != nil {
		unavailable = fmt.Sprintf("could not read connection string: %v", err)
		return m.Run()
	}

	return m.Run()
}

// URL returns the address of the Redis started by Run, in a form store.Open
// accepts. It skips the test when there is no server, or fails it under CI
// where a skip would quietly drop the coverage this suite exists to provide.
func URL(t *testing.T) string {
	t.Helper()
	if url != "" {
		return url
	}
	if os.Getenv("CI") != "" {
		t.Fatalf("no Redis available: %s", unavailable)
	}
	t.Skipf("no Redis available: %s\nThese tests start their own server — check that Docker is running.", unavailable)
	return ""
}
