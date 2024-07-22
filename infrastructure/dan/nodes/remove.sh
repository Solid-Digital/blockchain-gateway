kubectl delete -f statefulsets/

kubectl delete secret basic-auth
kubectl delete secret besu-validator1-key
kubectl delete secret besu-validator2-key
kubectl delete secret besu-validator3-key

kubectl delete -f configmap/
kubectl delete -f services/
kubectl delete -f ingresses/
kubectl delete -f namespace/