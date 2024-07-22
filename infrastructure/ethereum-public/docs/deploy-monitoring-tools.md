# Deploy monitoring tools

The monitoring tools are used to monitor all Ethereum networks, both public and private. This implies we deploy them only once.

For monitoring we use Grafana and Prometheus. Both of them are deployed in the `monitoring` namespace as follows:

## Create monitoring namespace

```bash
kubectl apply -f monitoring/namespace.yaml
```

## Deploy Grafana

```bash
kubectl apply -f monitoring/grafana-configmap.yaml
kubectl apply -f monitoring/grafana-deployment.yaml
kubectl apply -f monitoring/grafana-service.yaml
```

Grafana is available on port `:30030`, the inititial credentials are `admin:admin`. To expose it to the outside world you need to configure a loadbalancer or ingress.

## Deploy Prometheus

```bash
kubectl apply -f monitoring/prometheus-configmap.yaml
kubectl apply -f monitoring/prometheus-deployment.yaml
kubectl apply -f monitoring/prometheus-service.yaml
kubectl apply -f monitoring/prometheus-rbac.yaml
```

Prometheus is available on port `:30090`. To expose it to the outside world you need to configure a loadbalancer or ingress.

