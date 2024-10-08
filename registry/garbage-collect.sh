#!/bin/sh

while true; do
  sleep 72h # every 3 days
  echo "Starting garbage collection..."
  registry garbage-collect /etc/docker/registry/config.yml || true
done
