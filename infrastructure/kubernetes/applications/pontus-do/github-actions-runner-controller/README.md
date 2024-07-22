# GitHub Actions Runner Controller

Actions Runner Controller (https://github.com/summerwind/actions-runner-controller v0.10.0) is a kubernetes operator that handles the creation of self-hosted github action runners.


# Prerequisites
1. kubectl v1.18.6 - https://kubernetes.io/docs/tasks/tools/install-kubectl/
2. vault v1.4.3 - https://learn.hashicorp.com/vault/getting-started/install
3. envconsul v0.9.3 - https://github.com/hashicorp/envconsul#installation
4. docker v19.03.13 - https://docs.docker.com/get-docker/

run `make init` to get access to the vault secrets that are necessary for the deployment/upgrade process

# Deploy/Upgrade

Apply via `make deploy`

# Dry run

To check what changes will be made by the upgrade run:
`make dry-run`

# Docker image

To build and push the docker image that the runners use - 
    - make sure to set the DOCKER_IMAGE variable in the Makefile to update the image version
    - run `make image`
