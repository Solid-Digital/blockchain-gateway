#!/bin/bash

cat <<EOF | kubectl apply -f - 
{
  "apiVersion": "v1",
  "kind": "Secret",
  "metadata": {
    "name": "adapter-prod",
    "labels": {
      "app": "adapter-prod"
    },
    "namespace": "dan2"
  },
  "type": "Opaque",
  "stringData": $(vault kv get -format=json passdb/dan-do-adapter | jq ".data.data")
}
EOF

kubectl apply -f configmap.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml
