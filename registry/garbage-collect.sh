#!/bin/sh

echo "Script started"

while true; do
  echo "Starting garbage collection..."
  registry garbage-collect /etc/docker/registry/config.yml || true
  echo "Garbage collection finished"
  sleep 72h # every 3 days
done
