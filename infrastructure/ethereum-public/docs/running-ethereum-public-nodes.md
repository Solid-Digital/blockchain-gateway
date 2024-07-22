# Running Ethereum public nodes

This guide applies to Ethereum mainnet and ropsten. Nodes for other public Ethereum networks can be deployed as well, by making some minor adjustments to the statefulsets.

## Deploy Ethereum node

We have kube specs for both Ethereum mainnet and ropsten. They are deployed in the `mainnet-nodes` and `ropsten-nodes-fullsync` namespaces. Nodes for other public networks can be created by changing the command line arguments for the besu client in the statefulset.

### Ethereum mainnet

Run the following to deploy the Ethereum mainnet node:

```bash
kubectl apply -f mainnet/namespace.yaml
kubectl apply -f mainnet/configmap.yaml
kubectl apply -f mainnet/statefulset.yaml
kubectl apply -f mainnet/service.yaml
```

### Ethereum ropsten

Run the following to deploy the Ethereum ropsten node:

```bash
kubectl apply -f ropsten/namespace.yaml
kubectl apply -f ropsten/configmap.yaml
kubectl apply -f ropsten/statefulset.yaml
kubectl apply -f ropsten/service.yaml
```

## Persisten volume retention policy

When creating persistent volumes using volume claim templates, which we do in our statefulsets, the persistent volumes are created with reclaim policy `Delete` by default. This means that the persistent volume is deleted, when removing the persistent volume claim. Since we absolutely don't want to risk losing our blockchain database, we have to change the reclaim policy to `Retain`. Note that you have to remove persistent volumes manually when you don't need them anymore. Which is important, because they are costly.

To change the reclaim policy we have to fetch the persistent volumes first:

```bash
kubectl get pv
```

Then change the reclaim policy to `Retain`:

```bash
kubectl patch pv <persistent-volume-name> -p '{"spec": "persistentVolumeReclaimPolicy":"Retain"}}
```

## Deploy backup solution

Each node uses a persistent volume to store the blockchain database. We use k8s-snapshots to create snapshots from persistent volumes claims. Run the following to deploy k8s-snapshots in the `kube-system` namespace:

```bash
kubectl apply -f node-backup/deployment.yaml
kubectl apply -f node-backup/rbac.yaml
```

## Enable backups

We need to enable backups for the persistent volume claims of our nodes. Go the `mainnet-nodes` or `ropsten-nodes-fullsync` namespace and fetch a list of persisten volume claims as follows:

```bash
kubectl get pvc
```

And enable backups for this persistent volume claim:

```bash
kubectl patch pv <persistent-volume-claim-name> -p '{"metadata": {"annotations": {"backup.kubernetes.io/deltas": "P7D P7D"}}}'
```

This will create a snapsot every 7 days, and will keep each snapshot for 7 days after creating the next one. The snapshots are created through the Digitalocean API. Therefore, they are only visible on the Digitalocean dashboard under Images \ Snapshots \ Volumes.
