#!/bin/sh

echo "Starting garbage collection..."

registry garbage-collect /etc/docker/registry/config.yml
