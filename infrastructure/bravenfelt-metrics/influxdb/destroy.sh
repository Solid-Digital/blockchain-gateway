kubectl delete secret influxdb-secrets
kubectl delete -f configmap.yaml
kubectl delete -f deployment.yaml
kubectl delete -f service.yaml
kubectl delete -f pvc.yaml