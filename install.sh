#!/usr/bin/env bash

REMOTE_HOST=178.156.198.215
REMOTE_USER=root
INSTALL_DIR=/opt/ttlsh

function setup_docker_compose() {
    ssh $REMOTE_USER@$REMOTE_HOST "mkdir -p $INSTALL_DIR"
    scp docker-compose.yaml $REMOTE_USER@$REMOTE_HOST:$INSTALL_DIR/docker-compose.yaml
}

function setup_provisioner_script() {
    scp provision.sh $REMOTE_USER@$REMOTE_HOST:$INSTALL_DIR/provision.sh
    ssh $REMOTE_USER@$REMOTE_HOST "chmod a+x $INSTALL_DIR/provision.sh"
}

function run_provisioner_script() {
    ssh $REMOTE_USER@$REMOTE_HOST "bash $INSTALL_DIR/provision.sh $1"
}

function install() {
    setup_docker_compose
    setup_provisioner_script
    run_provisioner_script $DOPPLER_TOKEN
}

install