#!/usr/bin/env bash

# builds, tags, and pushes the images referenced in docker-compose.yaml
# (zot and redis are used off the shelf and are never built here)

docker compose build web zot-ephemeral-ttl
docker compose push web zot-ephemeral-ttl