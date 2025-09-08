#!/usr/bin/env bash

# builds, tags, and pushes the images referenced in docker-compose.yaml

docker compose build
docker compose push