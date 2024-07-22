## Bravenfelt Metrics

This directory contains the manifests files (including configurations) for Grafana and InfluxDB. 

### Init

For setup, run the deploy scripts inside both sub directories. The secrets are loaded from vault and require the following env variables to be set in the terminal:

```
VAULT_ADDR=https://vault.tools.unchain.io
VAULT_TOKEN="your-vault-token"
```

### Ingress

The ingresses for both services are as follows:

```
bravenfelt.grafana.unchain.io
bravenfelt.influxdb.unchain.io
```

### Destroy

Both deployments can either be destroyed separately through the `destroy.sh` script inside the subdirectories, or the `bravenfelt-metrics` namespace can be deleted via `kubectl delete ns bravenfelt-metrics` to remove both.