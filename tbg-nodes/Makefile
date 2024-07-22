.SHELLFLAGS := -eu -o pipefail -c
SHELL := bash

version = $(shell cat VERSION)
date := $(shell date -u '+%Y%m%d%H%M%S')
tag_version_date_replace = s/{{tag}}/${version}/g;s/{{date}}/${date}/g
remote_auth_service_image=registry.unchain.io/unchain/tbg-nodes-auth:v0.0.4
local_auth_service_image=localhost:5000/tbg-nodes-auth
deployment_target ?= staging

# Set to localhost:5000 to deploy locally, \
or run make commands from shell with regisry=localhost:5000 \
eg: make image-local registry=localhost:5000
registry = registry.unchain.io/unchain

new_image = $(registry)/tbg-nodes:${version}
docker_ip_cmd = docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'

sed_replace_cmd = sed "${tag_version_date_replace}"

tbg_nodes_access_kubeconfig_folder=./deployments_kustomize/base/step_2/secrets
local_tbg_nodes_config_folder=./deployments_kustomize/local/step_2/secrets
local_tbg_nodes_auth_config_folder=./deployments_kustomize/local/auth/secrets

kustomize_folder = ./deployments_kustomize/${deployment_target}

k8s_apply_dry_run := kubectl apply --dry-run=client -f -
k8s_apply := tee | kubectl apply -f -

tbg-nodes-access-kubeconfig:
	export NAMESPACE=tbg-nodes-${deployment_target} && \
		export DEPLOYMENT_TARGET=${deployment_target} && \
		mkdir -p ${tbg_nodes_access_kubeconfig_folder} && \
		cd ${tbg_nodes_access_kubeconfig_folder}  && \
		bash generate_tbg_nodes_access_kubeconfig.sh

tbg-nodes-access-kubeconfig-dry-run:
	mkdir -p ${tbg_nodes_access_kubeconfig_folder} && \
		cd ${tbg_nodes_access_kubeconfig_folder}  && \
		cp tbg_nodes_access.kubeconfig.template tbg_nodes_access.kubeconfig

deploy-template: write-version
	kubectl kustomize ${kustomize_folder}/step_1 | ${sed_replace_cmd} | ${k8s_apply_cmd}
	$(MAKE) ${create_tbg_hodes_access_kubeconfig}
	kubectl kustomize ${kustomize_folder}/step_2 | ${sed_replace_cmd} | ${k8s_apply_cmd}

build-local: deployment_target=local
build-local: k8s_apply_cmd=cat
build-local: create_tbg_hodes_access_kubeconfig=tbg-nodes-access-kubeconfig-dry-run
build-local: local-tbg-nodes-config-dry-run deploy-template

deploy-local-dry-run:
	$(MAKE) build-local \
		k8s_apply_cmd='${k8s_apply_dry_run}'

deploy-local: local-tbg-nodes-config
# Current deployment must be destroyed because new the pod \
	cannot startup as it needs port 4000, which it can't get \
	because the pod is using HostNetwork
	kubectl delete deployment tbg-nodes --ignore-not-found -n tbg-nodes-local
	docker pull ${remote_auth_service_image}
	docker tag ${remote_auth_service_image} ${local_auth_service_image}
	docker push ${local_auth_service_image}
	$(MAKE) deploy-template \
		deployment_target=local \
		create_tbg_hodes_access_kubeconfig=tbg-nodes-access-kubeconfig \
		k8s_apply_cmd='${k8s_apply}'

local-tbg-nodes-config-dry-run:
	mkdir -p ${local_tbg_nodes_config_folder}
	mkdir -p ${local_tbg_nodes_auth_config_folder}
	cd ${local_tbg_nodes_config_folder} && \
	touch local-tbg-nodes-config
	cd ${local_tbg_nodes_auth_config_folder} && \
	touch tbg-nodes-auth-dev.toml && \
	touch tbg-nodes-auth-test.toml

local-tbg-nodes-config:
	mkdir -p ${local_tbg_nodes_config_folder}
	mkdir -p ${local_tbg_nodes_auth_config_folder}
	$(shell cd ${local_tbg_nodes_config_folder} && [ -f local-tbg-nodes-config ] && rm local-tbg-nodes-config)
	cd ${local_tbg_nodes_config_folder} && \
		POSTGRES_HOST=$(shell ${docker_ip_cmd} tbg-nodes_postgres_1) REDIS_URL=redis://$(shell ${docker_ip_cmd} tbg-nodes_redis_1):6379 ./generate_tbg_nodes_local_config.sh
	cd ${local_tbg_nodes_auth_config_folder} && \
		POSTGRES_HOST=$(shell ${docker_ip_cmd} tbg-nodes_postgres_1) ./generate_tbg_nodes_auth_config.sh

build-staging: deployment_target=staging
build-staging: k8s_apply_cmd=cat
build-staging: create_tbg_hodes_access_kubeconfig=tbg-nodes-access-kubeconfig-dry-run
build-staging: deploy-template

build-prod:
	$(MAKE) build-staging \
		deployment_target=prod

deploy-staging-dry-run:
	$(MAKE) build-staging \
		k8s_apply_cmd='${k8s_apply_dry_run}'

deploy-prod-dry-run:
	$(MAKE) deploy-staging-dry-run \
		deployment_target=prod

deploy-%:
	$(MAKE) build-staging \
		k8s_apply_cmd='${k8s_apply}' \
		deployment_target=$* \
		create_tbg_hodes_access_kubeconfig=tbg-nodes-access-kubeconfig \

start:
	docker-compose up -d postgres redis pgadmin
	mix phx.server

stop:
	docker-compose stop

db-reset:
	docker-compose up -d
	docker-compose stop pgadmin
	mix ecto.reset
	docker-compose start pgadmin

db-migrate:
	docker-compose up -d
	docker-compose stop pgadmin
	mix ecto.migrate
	docker-compose start pgadmin

docker-login:
	docker login --username ${DOCKER_USERNAME} --password ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}

image: write-version
	echo Building $${MIX_ENV} image: ${new_image}
	docker build \
		--build-arg MIX_ENV=$${MIX_ENV} \
		-t ${new_image} .

push: write-version
	docker push ${new_image}

deployment-check-rollout-%:
	kubectl rollout status deployment/tbg-nodes -n tbg-nodes-$*

quick-full-deploy-local:
	KUBECONFIG=${PWD}/kubeconfig.local.yaml MIX_ENV=dev $(MAKE) full-deploy-local registry='localhost:5000'

deploy-local-auth: deployment_target=local
deploy-local-auth: KUBECONFIG=${PWD}/kubeconfig.local.yaml
deploy-local-auth: k8s_apply_cmd=${k8s_apply}
deploy-local-auth: local-tbg-nodes-config
	docker pull ${remote_auth_service_image}
	docker tag ${remote_auth_service_image} ${local_auth_service_image}
	docker push ${local_auth_service_image}
	export KUBECONFIG=${PWD}/kubeconfig.local.yaml && kubectl kustomize ${kustomize_folder}/auth | ${sed_replace_cmd} | ${k8s_apply_cmd}

quick-full-deploy-%:
	MIX_ENV=prod $(MAKE) full-deploy-$*

full-deploy-%:
	$(MAKE) image date=${date}
	$(MAKE) push date=${date}
	$(MAKE) deploy-$* date=${date}
	$(MAKE) deployment-check-rollout-$* date=${date}
	@echo Successfully deployed to $*!

write-version:
# Setting export TAG={value} will overwrite the version.
	echo $(shell if [ ! -z "${TAG}" ]; then echo "${TAG}"; else git describe --tags --always --dirty=-${date}; fi) > VERSION

local-k8s-port-forward:
	kubectl port-forward deployment/tbg-nodes -n tbg-nodes-local 4000

admin_user?=$(shell read -p "Enter admin email: " username; echo $$username)

set-admin-role:
	kubectl exec -n tbg-nodes-${deployment_target} -it deployment/tbg-nodes -- ./bin/tbg_nodes eval 'TbgNodes.Release.set_admin_role("${admin_user}")'

migrate:
	kubectl exec -n tbg-nodes-${deployment_target} -it deployment/tbg-nodes -- ./bin/tbg_nodes eval 'TbgNodes.Release.migrate()'

test: deploy-local-auth
	mix quality.ci
	mix test --include k8s

logs:
	kubectl logs -f -n tbg-nodes-${deployment_target} -c tbg-nodes deployment/tbg-nodes
