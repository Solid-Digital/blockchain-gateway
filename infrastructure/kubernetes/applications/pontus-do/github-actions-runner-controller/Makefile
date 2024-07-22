VAULT_ADDR=https://vault.tools.unchain.io
VAULT_USERNAME ?= $(shell read -p "Enter Vault username: " username; echo $$username)
VAULT_SECRET ?= passdb/actions-runner-controller

DOCKER_IMAGE ?= registry.unchain.io/unchainio/github-action-runner:v0.0.4

init: login 

login:
	vault login -method=ldap -address=${VAULT_ADDR} username=${VAULT_USERNAME}

env:
	envconsul -vault-addr ${VAULT_ADDR} -no-prefix -pristine -secret ${VAULT_SECRET} env

# load the environment variables from vault before executing the same target prefixed with underscore
%:
	envconsul -vault-addr ${VAULT_ADDR} -no-prefix -secret ${VAULT_SECRET} $(MAKE) _$*

# target to capture all undefined targets starting with _ in order to avoid infinite recursion
_%:
	@echo Undefined target \"$*\"

EXTRA_FLAGS ?= 

sed_replace=s@{{DOCKER_IMAGE}}@${DOCKER_IMAGE}@g

_deploy:
	kubectl apply ${EXTRA_FLAGS} -f k8s/actions-runner-controller.yaml
	
	kubectl create secret generic controller-manager \
		--dry-run=client -o yaml \
		-n actions-runner-system \
		--from-literal=github_app_id="$${APP_ID}" \
		--from-literal=github_app_installation_id="$${INSTALLATION_ID}" \
		--from-literal=github_app_private_key="$${PRIVATE_KEY}" | kubectl apply ${EXTRA_FLAGS} -f -

	cat k8s/tbg-nodes-runner.yaml | sed "${sed_replace}" | kubectl apply ${EXTRA_FLAGS} -f -

_dry-run:
	make deploy EXTRA_FLAGS='--dry-run=client'

image:
	docker build ./docker --tag ${DOCKER_IMAGE}
	docker push ${DOCKER_IMAGE}
