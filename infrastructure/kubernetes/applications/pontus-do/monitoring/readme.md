# do-pontus Monitoring
-  https://monitoring.unchain.io/grafana
-  https://monitoring.unchain.io/alertmanager
-  https://monitoring.unchain.io/prometheus

# Deploying Updates
## Required Packages.

- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/ )
- [vault](https://www.vaultproject.io/docs/install)
- [doctl](https://www.digitalocean.com/docs/apis-clis/doctl/how-to/install/)

## Required Environment Variables

- `KUBECONFIG`: needed by `kubectl`
- `DIGITALOCEAN_ACCESS_TOKEN`: needed by `doctl`
  - alternatively, if you authenticate via `doctl auth init`  (more [info](https://github.com/digitalocean/doctl#authenticating-with-digitalocean))
  you don't need to set the digital ocean env var
- `VAULT_TOKEN`: needed by `vault`
  - this can be copied from the vault ui (top right dropdown menu -> copy token)
  -  alternatively you can authenticate via `vault login -method=ldap -address=https://vault.tools.unchain.io username={your user name}`
  
  
## Deploy
- `make deploy`
- `make deploy` uses the `kubectl apply` with the `--prune` option, this means when resources are removed
and `make deploy` is run, they will be automatically removed
- if you want to see what would deleted run `make check` and look for `pruned` resources
- updating configsmaps or secrets will trigger a redeploy of all the affected statefulsets/deployments/etc since
kustomize appends a hash to the configmap/secret name based on the resources' contents

## Check/Dry run
- `make check`

## View Kubernetes yaml
- `make build`

## Delete local secrets
- when running the above commands, secrets are pulled from vault & digital ocean
- running `make clean` deletes these secrets

# Additional Information

## Oauth2

- the oauth2 deployment authenticates via a [github application](https://github.com/organizations/unchain/settings/applications/1444776)
- alertmanager, grafana and prometheus all use this same github applcation
- [documentation](https://oauth2-proxy.github.io/oauth2-proxy/docs/)

## Alertmanager

- provides alerts in slack via [webhook](https://projectvloed.slack.com/services/B01HR73E7LJ)
- can also be configured to alert via email (or other [integrations](https://prometheus.io/docs/operating/integrations/))

## Prometheus

- copied alerting files from [helm chart](https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack/templates/prometheus/rules)
- only copied what I though were the [essential/useful ones](./configmaps/prometheus/prometheus_rules)


## Grafana

### High availability
- alertmanager & prometheus are running with more than one replica each
  - more [information](https://prometheus.io/docs/introduction/faq/#can-prometheus-be-made-highly-available)


### Updating graphs

- In order to commit changes to graphs back to version control, once a graph is updated and saved
you need to  `Dashboard settings> JSON Model` and copy it. Replace the version in [this repo](./configmaps/grafana/dashboards/) with the new version.
- This isn't ideal, but for now, especially in the early days, it means if something gets messed up it
is easy to revert.
- Grafana also has a built-in history of changes to dashboards it would also be possible to rely on that
