# syntax=docker/dockerfile:1.7
#
# Multi-stage build for the zot-ephemeral-ttl sidecar. This file builds; it
# does not gate — formatting, vet, lint, and tests run in CI.

# ---- build arguments -------------------------------------------------------
ARG GO_VERSION=1.26
ARG ALPINE_VERSION=3.22

# Cross-compilation: BuildKit injects TARGETOS/TARGETARCH for the
# requested --platform; default to the build host's values when invoked
# without --platform (e.g. plain `docker compose build`).
ARG TARGETOS=linux
ARG TARGETARCH

# Version metadata wired into the binary via -ldflags.
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown


# ---- deps: download modules once and cache them ---------------------------
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS deps
WORKDIR /src

# git for VCS info, ca-certs in case `go mod download` talks to a private proxy.
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download -x


# ---- build: produce the static binary -------------------------------------
FROM deps AS build
COPY . .
ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG COMMIT
ARG BUILD_DATE

# CGO_ENABLED=0 + pure-Go deps => the binary runs from `FROM scratch`.
# -buildvcs=false keeps host git state out; the metadata is set via -X below.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build \
        -trimpath \
        -buildvcs=false \
        -ldflags="-s -w \
            -X main.version=${VERSION} \
            -X main.commit=${COMMIT} \
            -X main.buildDate=${BUILD_DATE}" \
        -o /out/zot-ephemeral-ttl \
        ./cmd/zot-ephemeral-ttl


# ---- runtime: scratch image with just the binary --------------------------
FROM scratch AS runtime

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

LABEL org.opencontainers.image.title="zot-ephemeral-ttl" \
      org.opencontainers.image.description="Tag-driven TTL sidecar for zot. Subscribes to image.updated CloudEvents and DELETEs expired manifests." \
      org.opencontainers.image.source="https://github.com/nullbytelabs/zot-ephemeral-ttl" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.revision="${COMMIT}" \
      org.opencontainers.image.created="${BUILD_DATE}"

COPY --from=build /out/zot-ephemeral-ttl /zot-ephemeral-ttl

# scratch has no shell, so HEALTHCHECK is intentionally omitted; probe
# /healthz from outside.

EXPOSE 8080
ENTRYPOINT ["/zot-ephemeral-ttl"]
