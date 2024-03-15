#!/bin/sh

set -e

sed -i "s/__PORT__/$PORT/g" /etc/nginx/nginx.conf

exec "$@"
