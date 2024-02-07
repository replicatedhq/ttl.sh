#!/bin/sh

set -e

sed -i "s/__PORT__/$PORT/g" /etc/docker/registry/config.yml
sed -i "s/__HOOK_TOKEN__/$HOOK_TOKEN/g" /etc/docker/registry/config.yml
sed -i "s/__HOOK_URI__/$HOOK_URI/g" /etc/docker/registry/config.yml
sed -i "s/__REPLREG_HOST__/$REPLREG_HOST/g" /etc/docker/registry/config.yml
sed -i "s/__REPLREG_SECRET__/$REPLREG_SECRET/g" /etc/docker/registry/config.yml

if [[ -z "${GCS_KEY_ENCODED}" ]]; then
  echo "Set GCS_KEY_ENCODED variable"
else
  echo ${GCS_KEY_ENCODED} | base64 -d > /etc/gcs.json
  chmod 0400 /etc/gcs.json
fi

# Run garbage collection job in background
# /garbage-collect.sh &

case "$1" in
    *.yaml|*.yml) set -- registry serve "$@" ;;
    serve|garbage-collect|help|-*) set -- registry "$@" ;;
esac

exec "$@"
