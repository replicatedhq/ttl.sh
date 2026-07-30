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

  # The inventory has no address in it; both the SSH target and the Cloudflare
  # A record value come from this one secret.
  local vm_ip
  vm_ip=$(get_secret LATITUDESH_VM_IP)

  ansible-playbook \
    -e "ansible_host=${vm_ip}" \
    -e "target_ip=${vm_ip}" \
    -e "cloudflare_api_token=$(get_secret CF_API_TOKEN)" \
    -e "cloudflare_zone_id=$(get_secret CF_ZONE_ID)" \
    -e "cloudflare_email=$(get_secret LE_EMAIL)" \
    -e "le_email=$(get_secret LE_EMAIL)" \
    -e "s3_bucket=$(get_secret S3_BUCKET)" \
    -e "s3_region=$(get_secret S3_REGION)" \
    -e "s3_endpoint=$(get_secret S3_ENDPOINT)" \
    -e "s3_access_key=$(get_secret S3_ACCESS_KEY)" \
    -e "s3_secret_key=$(get_secret S3_SECRET_KEY)" \
    playbooks/site.yml
}

main "$@"

