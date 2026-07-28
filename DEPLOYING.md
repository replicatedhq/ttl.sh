# Deploying ttl.sh manually

Infrastructure provisioning is no longer part of this repo. The playbook expects an
already-running Ubuntu host that you can reach as `root` over SSH, listed in
`ansible/inventory/hosts`, with its public IP set as `target_ip` in
`ansible/inventory/group_vars/all.yml` (that value is what the Cloudflare A record is
pointed at).

## Build the ttl.sh images

The `web` and `zot-ephemeral-ttl` services are built here — zot and redis are used
off the shelf.

1. Docker and Docker Compose are installed
2. Authenticated to `ghcr.io` (`docker login ghcr.io`)
3. `./build-and-push.sh`

## Server setup and ttl.sh workloads

1. Doppler is installed and authenticated for use with the `ttl-sh` project
    a. Have a valid `DOPPLER_TOKEN` from `ttl-sh` set
    b. `echo $DOPPLER_TOKEN | doppler configure set token --scope /`
2. These secrets exist in Doppler and describe the S3-compatible bucket backing the
   registry: `S3_BUCKET`, `S3_REGION`, `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`
3. `cd ansible`
4. `./ansible.sh`

The playbook installs Nginx and Certbot, obtains a certificate via Cloudflare DNS-01,
installs Docker, then copies `docker-compose.yaml` and the rendered `zot-config.json` to
`/opt/ttlsh` and brings the stack up.
