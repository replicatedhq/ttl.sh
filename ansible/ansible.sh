#!/usr/bin/env bash

set -euo pipefail

# Minimal wrapper to run Ansible with Cloudflare vars resolved via Doppler.
# Usage:
#   bash ansible.sh

HERE="$(cd "$(dirname "$0")" && pwd)"

require_bin() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "error: required command not found: $1" >&2
    exit 1
  }
}

require_env() {
  local name=$1
  if [[ -z "${!name:-}" ]]; then
    echo "error: required env var not set: $name" >&2
    exit 1
  fi
}

get_secret() {
  local key=$1
  # Uses DOPPLER_TOKEN from environment; Doppler CLI will read it.
  # --plain prints only the value with a trailing newline, strip it.
  local val
  if ! val=$(doppler secrets get "$key" --plain 2>/dev/null); then
    echo "error: failed to fetch secret '$key' from Doppler" >&2
    exit 1
  fi
  printf %s "$val"
}

main() {
  require_bin ansible-playbook
  require_bin doppler
  require_env DOPPLER_TOKEN

  # Always run from the ansible dir so relative paths match
  cd "$HERE"

  ansible-playbook \
    -e "cloudflare_api_token=$(get_secret CF_API_TOKEN)" \
    -e "cloudflare_zone_id=$(get_secret CF_ZONE_ID)" \
    -e "cloudflare_email=$(get_secret LE_EMAIL)" \
    -e "le_email=$(get_secret LE_EMAIL)" \
    -e "gcloud_sa_email=$(get_secret GOOGLE_APPLICATION_CREDENTIALS_EMAIL)" \
    -e "gcloud_sa_key_json=$(get_secret GCS_KEY_ENCODED)" \
    -e "gcloud_registry_host=$(get_secret GOOGLE_CLOUD_ARTIFACT_REGISTRY_URI)" \
    -e "gcs_key_encoded=$(get_secret GCS_KEY_ENCODED)" \
    -e "hook_token=$(get_secret HOOK_TOKEN)" \
    -e "hook_uri=$(get_secret HOOK_URI)" \
    -e "redis_password=$(get_secret REDIS_PASSWORD)" \
    -e "rediscloud_url=$(get_secret REDISCLOUD_URL)" \
    -e "replreg_host=$(get_secret REPLREG_HOST)" \
    -e "replreg_secret=$(get_secret REPLREG_SECRET)" \
    -e "registry_url=$(get_secret REGISTRY_URL)" \
    playbooks/site.yml
}

main "$@"

