kubectl delete secret grafana-creds
kubectl delete -f ./configmap.yaml
kubectl delete -f ./deployment.yaml
kubectl delete -f ./service.yaml