kubectl apply -f namespace/

cat <<EOF | kubectl apply -f - 
{
  "apiVersion": "v1",
  "kind": "Secret",
  "metadata": {
    "name": "basic-auth",
    "labels": {
      "app": "basic-auth"
    },
    "namespace": "dan2"
  },
  "type": "Opaque",
  "stringData": $(vault kv get -format=json passdb/dan-nodes/basic-auth | jq ".data.data")
}
EOF

cat <<EOF | kubectl apply -f - 
{
  "apiVersion": "v1",
  "kind": "Secret",
  "metadata": {
    "name": "besu-validator1-key",
    "labels": {
      "app": "besu-validator1-key"
    },
    "namespace": "dan2"
  },
  "type": "Opaque",
  "stringData": $(vault kv get -format=json passdb/dan-nodes/validator1 | jq ".data.data")
}
EOF

cat <<EOF | kubectl apply -f - 
{
  "apiVersion": "v1",
  "kind": "Secret",
  "metadata": {
    "name": "besu-validator2-key",
    "labels": {
      "app": "besu-validator2-key"
    },
    "namespace": "dan2"
  },
  "type": "Opaque",
  "stringData": $(vault kv get -format=json passdb/dan-nodes/validator2 | jq ".data.data")
}
EOF

cat <<EOF | kubectl apply -f - 
{
  "apiVersion": "v1",
  "kind": "Secret",
  "metadata": {
    "name": "besu-validator3-key",
    "labels": {
      "app": "besu-validator3-key"
    },
    "namespace": "dan2"
  },
  "type": "Opaque",
  "stringData": $(vault kv get -format=json passdb/dan-nodes/validator3 | jq ".data.data")
}
EOF


kubectl apply -f configmap/
kubectl apply -f services/
kubectl apply -f statefulsets/
kubectl apply -f ingresses/