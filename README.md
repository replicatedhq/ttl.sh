[![Develop on Okteto](https://okteto.com/develop-okteto.svg)](https://replicated.okteto.dev/deploy?repository=https://github.com/replicatedhq/ttl.sh&branch=main)

# ttl.sh

## An ephemeral container registry for CI workflows.

## What is ttl.sh?

ttl.sh is an anonymous, expiring Docker container registry using the official Docker Registry image. This is a set of tools and configurations that can be used to deploy the registry without authentication, but with self-expiring images.

# Development

Development for the services in this project is done through [Okteto](https://replicated.okteto.dev).

## Setup

1. Install the Okteto CLI (`brew install okteto`)
2. Setup Okteto CLI (`okteto context use https://replicated.okteto.dev`)
3. Setup Okteto context in kubectl (`okteto context update-kubeconfig`)
4. Deploy your current branch. (from the ttl.sh root directory: `okteto pipeline deploy`)

## Debugging

Okteto is utilized for debugging. New build targets have been added to allow building and running each service in debug mode.

1. Replace the default container in your Okteto environment with a development container.
   1. From the root directory: `okteto up` or `okteto up <service name>`
2. Run the build targets for the desired service:
   1. ttl-hooks: `make deps build hooks`
   2. ttl-reaper: `make deps build reap`
3. Stop development and go back to the default container.
   1. From the root directory: `okteto down` or `okteto down <service name>`

## Example workflows

### Switching branches or rebasing

1. `git checkout my-new-branch` 
2. `okteto pipeline deploy` 
3. (make code changes)
4. `okteto up`
5. (test changes, find they don't work, make more changes)...
6. `okteto down`
7. (commit code, and be happy)


## Running the Registry Locally (MacBook Example)

1. **Prepare the data directory:**

   ```sh
   cd registry/
   mkdir -p data
   chmod 777 data
   ```

2. **Build the Docker image:**

   ```sh
   docker build -t my-ttlsh:latest .
   ```

3. **Run the registry container:**

   ```sh
   docker run -d \
     -p 5000:5000 \
     --name ttlsh \
     -e PORT=5000 \
     -e HOOK_TOKEN=localtoken \
     -e HOOK_URI=http://localhost:9999/hook \
     -e REPLREG_HOST=localhost:5000 \
     -e REPLREG_SECRET=localsecret \
     -v ttlsh-data:/var/lib/registry \
     my-ttlsh:latest
   ```

4. **Verify the registry is running:**

   - Check logs:
     ```sh
     docker logs ttlsh
     ```
   - Test the registry API:
     ```sh
     curl http://localhost:5000/v2/
     ```

This will start the Docker registry with your local configuration and environment variables.