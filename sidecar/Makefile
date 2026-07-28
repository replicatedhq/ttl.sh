# The Go toolchain is the primary interface to this repo: `go test ./...` works
# on a clean checkout. This file covers the two things that need more than a
# one-word go command — the checks CI gates on, and a version-stamped image.

IMAGE ?= zot-ephemeral-ttl
TAG   ?= dev

.DEFAULT_GOAL := test

.PHONY: test
test:
	go test -race ./...

# Everything CI gates on, in the order CI runs it.
.PHONY: check
check:
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "not gofmt-formatted:"; echo "$$unformatted"; \
		echo "run 'gofmt -w .' to fix."; exit 1; \
	fi
	go vet ./...
	@cp go.mod .go.mod.tidycheck && cp go.sum .go.sum.tidycheck; \
	go mod tidy; \
	rc=0; \
	if ! cmp -s go.mod .go.mod.tidycheck || ! cmp -s go.sum .go.sum.tidycheck; then \
		echo "go.mod/go.sum were not tidy; 'go mod tidy' has fixed them — commit the result."; \
		rc=1; \
	fi; \
	rm -f .go.mod.tidycheck .go.sum.tidycheck; \
	exit $$rc
	go test -race ./...

.PHONY: image
image:
	docker build \
		--build-arg VERSION=$$(git describe --tags --always 2>/dev/null || echo dev) \
		--build-arg COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo unknown) \
		--build-arg BUILD_DATE=$$(date -u +%Y-%m-%dT%H:%M:%SZ) \
		--tag $(IMAGE):$(TAG) \
		.
