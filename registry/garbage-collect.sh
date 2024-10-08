#!/bin/sh

while true; do
  sleep 1m
  echo "Starting garbage collection..."
  registry garbage-collect /etc/docker/registry/config.yml || true
done
