cat <<EOF | kubectl apply -f - 
{
  "apiVersion": "v1",
  "kind": "Secret",
  "metadata": {
    "name": "influxdb-secrets",
    "namespace": "bravenfelt-metrics"
  },
  "type": "Opaque",
  "stringData": $(vault kv get -format=json passdb/bravenfelt-metrics/influxdb | jq ".data.data")
}
EOF

kubectl apply -f configmap.yaml
kubectl apply -f pvc.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml