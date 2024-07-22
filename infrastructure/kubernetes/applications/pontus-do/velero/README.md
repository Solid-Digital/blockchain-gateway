# Velero

Velero is a backup tool for kubernetes.

# Prerequisites
1. kubectl v1.18.6 - https://kubernetes.io/docs/tasks/tools/install-kubectl/
2. vault v1.4.3 - https://learn.hashicorp.com/vault/getting-started/install
3. envconsul v0.9.3 - https://github.com/hashicorp/envconsul#installation
4. velero CLI v1.4.2 - https://velero.io/docs/v1.4/basic-install/

run `make init` to get access to the vault secrets that are necessary for the deployment/upgrade process and add the repo necessary for helm

# Deploy/Upgrade

Apply via `make deploy`

# Test

Test if the deployment/upgrade was successful via: 
`make test`

# Check

To check what changes will be made by the upgrade run:
`make check`

# Usage

## Single backup

`velero backup create <backup-name> --include-namespaces <namespace>`

## Restore from backup

`velero restore create --from-backup <backup-name>`

## Scheduled backup

`velero schedule create <schedule-name> --schedule "0 7 * * *" --include-namespaces <namespace>`

## Restore from scheduled backup

`velero restore create --from-schedule <schedule-name>`

## Further reading

1. https://velero.io/docs/v1.4/examples/
2. https://velero.io/docs/v1.4/disaster-case/
3. https://velero.io/docs/v1.4/migration-case/
