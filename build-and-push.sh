#!/usr/bin/env bash

# builds, tags, and pushes the web image referenced in docker-compose.yaml
# (zot is used off the shelf and is never built here)

docker compose build web
docker compose push web