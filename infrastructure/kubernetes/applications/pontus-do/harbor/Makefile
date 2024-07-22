SHELL=bash

# In order to fetch the secrets from vault; install the vault cli and issue the 
# following command: make login

vault ?= true

VAULT_USERNAME ?= $(shell read -p "Enter Vault username: " username; echo $$username)
VAULT_ADDR ?= https://vault.tools.dgo.unchain.io
VAULT_SECRET ?= passdb/accounts/harbor

HARBOR_CHART ?= bitnami/harbor
HARBOR_CHART_VERSION ?= 7.0.3

HELM_VALUES_FILE ?= values.yaml
HELM_SECRETS ?= --set harborAdminPassword=${harborAdminPassword} \
				--set core.secret=${core.secret} \
				--set core.secretKey=${secretKey} \
				--set postgresql.postgresqlPassword=${postgresqlPassword} \
				--set persistence.imageChartStorage.s3.accesskey=${S3_ACCESS_KEY} \
				--set persistence.imageChartStorage.s3.secretkey=${S3_SECRET_KEY}

init: login helm-add-repo

login:
	vault login -method=ldap -address=${VAULT_ADDR} username=${VAULT_USERNAME}

helm-add-repo:
	helm repo add bitnami https://charts.bitnami.com/bitnami

# target to capture all undefined targets starting with _ in order to avoid infinite recursion
_%:
	@echo Undefined target \"$*\"

# if vault is true, then load the environment variables from vault before executing the same target prefixed with underscore
%:
ifeq ($(vault),true)
	envconsul -vault-addr ${VAULT_ADDR} -no-prefix -secret ${VAULT_SECRET} $(MAKE) _$*
else
	$(MAKE) _$*
endif

_secrets:
	@echo '$$core.secret' = $(core.secret)
	@echo '$$harborAdminPassword' = $(harborAdminPassword)
	@echo '$$secretKey' = $(secretKey)
	@echo '$$postgresqlPassword' = $(postgresqlPassword)
	@echo '$$S3_ACCESS_KEY' = $(S3_ACCESS_KEY)
	@echo '$$S3_SECRET_KEY' = $(S3_SECRET_KEY)

env:
	envconsul -vault-addr ${VAULT_ADDR} -no-prefix -pristine -secret ${VAULT_SECRET} env

_deploy:
	helm upgrade \
		--install harbor \
		--namespace harbor --create-namespace \
		-f ${HELM_VALUES_FILE} \
		${HELM_SECRETS} \
		${HARBOR_CHART} --version ${HARBOR_CHART_VERSION} --debug

_check:
	helm upgrade --install harbor -f ${HELM_VALUES_FILE} \
	${HELM_SECRETS} \
	${HARBOR_CHART} --version ${HARBOR_CHART_VERSION} --debug --dry-run

test:
	docker pull registry.dgo.unchain.io/unchainio/whoami

backup-image:
	docker build -t registry.dgo.unchain.io/unchain/backup-tool ./backup-tool/
	docker push registry.dgo.unchain.io/unchain/backup-tool

deploy-backup-tool:
	kubectl apply -f ./backup-tool/harbor_pg_backup_cronjob.yaml

backup-job:
	kubectl create job --from=cronjob/harbor-db-backup-cronjob manual2
