# Important!
# Version is either set to the the TAG env var. (when running in Docker), or the git tag.
version := $(shell if [ ! -z ${TAG} ]; then echo "${TAG}"; else git describe --tags --always; fi)
date := $(shell date)
tag := ${TAG}

info:
	@echo Make for deployments. Try \"make deployment env=staging\"

version:
	echo ${version}

check-env:
	@if [ -z ${env} ]; then echo "env [dev,staging,prod] is not set!"; exit 2; fi 

deployment: check-env
	sed "s/{{tag}}/${version}/g;s/{{date}}/${date}/g" ${env}/ares-deployment.template.yaml > ${env}/ares-deployment.yaml
	kubectl apply -f ${env}/ares-deployment.yaml -n ares-${env}
	@if [ ${env} = production ]; then \
		echo "succesfully deployed to ares.unchain.io"; \
	else \
		echo "succesfully deployed to ares.${env}.unchain.io"; \
	fi

service: check-env
	kubectl apply -f deployments/${env}/ares-service.yaml -n ares-${env}

logs:
	kubectl get pods -l app.kubernetes.io/name=ares -o=name | xargs kubectl logs


#####
## Scripts for running in bitbucket, this is not for you.
#####

# This is for running in the bitbucket pipeline only. Use your own Pontus credentials
context:
	kubectl config set-context staging --user=builduser --namespace=ares-staging --cluster=pontus \
	&& kubectl config use-context staging

# This is for running in the bitbucket pipeline only. Use your own Pontus credentials
cluster:
	kubectl config set-cluster pontus --server=https://rancher.aws.unchain.io/k8s/clusters/c-dbjpm

# This is for running in the bitbucket pipeline only. Use your own Pontus credentials
kubeconfig:
	@[ "${namespace}" ] || ( echo ">> namespace is not set"; exit 1 )
	kubectl config set-cluster pontus --server=https://rancher.aws.unchain.io/k8s/clusters/c-dbjpm
	kubectl config set-credentials builduser --username=${K8S_USER} --password=${K8S_PASSWORD}
	kubectl config set-context ares --user=builduser --namespace=$(namespace) --cluster=pontus
	kubectl config use-context ares