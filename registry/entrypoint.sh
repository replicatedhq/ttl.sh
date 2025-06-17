#!/bin/sh
set -e

sed -i "s|__PORT__|$PORT|g" /etc/docker/registry/config.yml
sed -i "s|__HOOK_TOKEN__|$HOOK_TOKEN|g" /etc/docker/registry/config.yml
sed -i "s|__HOOK_URI__|$HOOK_URI|g" /etc/docker/registry/config.yml
sed -i "s|__REPLREG_HOST__|$REPLREG_HOST|g" /etc/docker/registry/config.yml
sed -i "s|__REPLREG_SECRET__|$REPLREG_SECRET|g" /etc/docker/registry/config.yml

case "$1" in
    *.yaml|*.yml) set -- registry serve "$@" ;;
    serve|garbage-collect|help|-*) set -- registry "$@" ;;
esac

exec "$@"
