# Kubed

Kubed is a kubernetes tool for copying a secret across multiple namespaces. It is useful for cert-manager when several namespaces need the same tls certificate.

# Prerequisites
1. kubectl v1.18.6 - https://kubernetes.io/docs/tasks/tools/install-kubectl/
2. helm v3.3.0 - https://helm.sh/docs/intro/install/

run `make init` add the repo necessary for helm

# Deploy/Upgrade

Apply via `make deploy`

# Check

To check what changes will be made by the upgrade run:
`make check`

