#!/bin/sh

while true; do
  echo "Starting garbage collection..."
  registry garbage-collect /etc/docker/registry/config.yml || true
  sleep 24h
done
