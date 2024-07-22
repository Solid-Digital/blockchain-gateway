# Harbor

# Prerequisites
1. kubectl v1.18.6 - https://kubernetes.io/docs/tasks/tools/install-kubectl/
2. vault v1.4.3 - https://learn.hashicorp.com/vault/getting-started/install
3. envconsul v0.9.3 - https://github.com/hashicorp/envconsul#installation
4. helm v3.3.0 - https://helm.sh/docs/intro/install/

run `make init` to get access to the vault secrets that are necessary for the deployment/upgrade process and add the repo necessary for helm

# Deploy/Upgrade

Apply via `make deploy`

# Test

Test if the deployment/upgrade was successful via: 
`make test`

# Check

To check what changes will be made by the upgrade run:
`make check`

