# Deploying ttl.sh manually

## Hetzner Resources via Terraform

1. Doppler is installed and authenticated for use with the `ttl-sh` project
    a. Have a valid `DOPPLER_TOKEN` from `ttl-sh` set
    b. `echo $DOPPLER_TOKEN | doppler configure set token --scope /`
2. Export necessary doppler values to ENVs
    a. `export HCLOUD_TOKEN=$(doppler secrets get HCLOUD_TOKEN --plain)`
    b. `export AWS_ACCESS_KEY_ID=$(doppler secrets get HCLOUD_S3_ACCESS_KEY --plain)`
    c. `export AWS_SECRET_ACCESS_KEY=$(doppler secrets get HCLOUD_S3_SECRET_KEY --plain)`
3. `cd terraform`
4. `terraform plan`
    a. Make sure plan looks as expected
5. `terraform apply`
    a. approve with `yes`
6. Successful run states `Apply complete!`

**NOTE**: For extra keys to be associated they must be added in the ttl-sh project. This is crucial to enable proper execution of installers.

## Build ttl.sh Docker images
1. Docker and Docker Compose are installed
2. `gcloud` CLI is installed and authenticated
3. `gcloud auth configure-docker us-east4-docker.pkg.dev`
4. `./build-and-push.sh`

## Server Setup and ttl.sh workloads

1. Doppler is installed and authenticated for use with the `ttl-sh` project
    a. Have a valid `DOPPLER_TOKEN` from `ttl-sh` set
    b. `echo $DOPPLER_TOKEN | doppler configure set token --scope /`
2. `cd ansible`
3. `./ansible.sh`
4. Successful run indicates approximately the following:
    a. `ok=42   changed=6    unreachable=0    failed=0    skipped=3`
