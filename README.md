# ttl.sh

## An ephemeral container registry for CI workflows.

## What is ttl.sh?

ttl.sh is an anonymous, expiring container registry built on [zot](https://zotregistry.dev).
This repo holds the host configuration (Ansible), the Next.js marketing site (`web/`), and
the Docker Compose stack that runs on the server. It does not provision infrastructure —
it configures a host that already exists.

## Layout

| Path                  | What it is                                             |
| --------------------- | ------------------------------------------------------ |
| `ansible/`            | Nginx + TLS, Docker, and the Compose deployment        |
| `web/`                | Next.js site served at https://ttl.sh                  |
| `static/`             | Legacy static site assets                              |
| `docker-compose.yaml` | The two services that run on the host: `web` and `zot` |

## How it runs

Nginx terminates TLS and proxies:

- `/` to the `web` container on port 3000
- `/v2` to the `zot` container on port 5000 (`/v2/_catalog` is blocked)

Zot runs an off-the-shelf image (`ghcr.io/project-zot/zot-linux-amd64`) with no local
customization. Its only configuration is `zot-config.json`, rendered onto the host by
Ansible from `ansible/templates/zot-config.json.j2`, which points zot at an S3-compatible
bucket for blob storage.

Tag expiry is handled separately and is not yet wired up in this repo.

## Deploying

See [DEPLOYING.md](DEPLOYING.md). Pushes to the deploy branch also run
`.github/workflows/deploy.yml`, which builds and pushes the `web` image and then runs the
Ansible playbook.
