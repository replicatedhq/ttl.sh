#!/usr/bin/env bash

DOPPLER_TOKEN=$1

function setup_doppler() {
    # install for ubuntu: https://docs.doppler.com/docs/cli
    sudo apt-get update && sudo apt-get install -y apt-transport-https ca-certificates curl gnupg
    curl -sLf --retry 3 --tlsv1.2 --proto "=https" 'https://packages.doppler.com/public/cli/gpg.DE2A7741A397C129.key' | sudo gpg --dearmor -o /usr/share/keyrings/doppler-archive-keyring.gpg
    echo "deb [signed-by=/usr/share/keyrings/doppler-archive-keyring.gpg] https://packages.doppler.com/public/cli/deb/debian any-version main" | sudo tee /etc/apt/sources.list.d/doppler-cli.list
    sudo apt-get update && sudo apt-get install doppler

    # authenticate with service token
    echo $DOPPLER_TOKEN | doppler configure set token --scope /
}

setup_docker() {
    # Adapted from:
    # * https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository

    # Add Docker's official GPG key:
    sudo apt-get update
    sudo apt-get install ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc

    # Add the repository to Apt sources:
    echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
    $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
    sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt-get update

    # install necessary packages
    sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y
}

setup_gcloud() {
    # Adapted from:
    # * https://cloud.google.com/sdk/docs/downloads-interactive
    # * Artifact Registry setup commands

    INSTALL_PATH="$(doppler secrets get GOOGLE_APPLICATION_CREDENTIALS --plain)"
    INSTALL_EMAIL="$(doppler secrets get GOOGLE_APPLICATION_CREDENTIALS_EMAIL --plain)"
    REGISTRY_URI="$(doppler secrets get GOOGLE_CLOUD_ARTIFACT_REGISTRY_URI --plain)"

    mkdir -p "$(dirname "$INSTALL_PATH")"
    doppler secrets get GOOGLE_APPLICATION_CREDENTIALS_JSON --plain >"$INSTALL_PATH"
    chmod 600 "$INSTALL_PATH"

    echo "export GOOGLE_APPLICATION_CREDENTIALS=\"$INSTALL_PATH\"" >> "$HOME/.bashrc"

    # download the installer
    curl https://sdk.cloud.google.com | bash

    # Explicitly add gcloud to PATH for this session
    export PATH="$HOME/google-cloud-sdk/bin:$PATH"

    gcloud auth activate-service-account "$INSTALL_EMAIL" --key-file="$INSTALL_PATH" --quiet
    gcloud auth configure-docker "$REGISTRY_URI" --quiet
}

setup_ttlsh_env() {
  for VAR in \
    GCS_KEY_ENCODED \
    HOOK_TOKEN \
    HOOK_URI \
    REDIS_PASSWORD \
    REDISCLOUD_URL \
    REPLREG_HOST \
    REPLREG_SECRET \
    REGISTRY_URL
  do
    echo "$VAR=$(doppler secrets get "$VAR" --plain)" >> /opt/ttlsh/.env
  done

  chmod 600 "$ENV_FILE"
}

function run_ttlsh()
{
    docker compose -f /opt/ttlsh/docker-compose.yaml pull
    docker compose -f /opt/ttlsh/docker-compose.yaml up -d
}

function provision() {
    setup_doppler
    setup_docker
    setup_gcloud
    setup_ttlsh_env
    run_ttlsh
}

provision