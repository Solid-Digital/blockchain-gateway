#!/bin/bash
set -e
secret_name="$(kubectl get sa tbg-nodes-access-${DEPLOYMENT_TARGET} --namespace=${NAMESPACE} -o jsonpath='{.secrets[0].name}')"
ca="$(kubectl get secret/${secret_name} --namespace=${NAMESPACE} -o jsonpath='{.data.ca\.crt}')"
token="$(kubectl get secret/${secret_name} --namespace=${NAMESPACE} -o jsonpath='{.data.token}' | base64 --decode)"
kubeconfig_replace="s/{{secret_name}}/${secret_name}/g;s/{{ca}}/${ca}/g;s/{{token}}/${token}/g"

cat tbg_nodes_access.kubeconfig.template | sed "${kubeconfig_replace}" > tbg_nodes_access.kubeconfig
