VAULT_ADDR=https://vault.tools.dgo.unchain.io
VAULT_USERNAME ?= $(shell read -p "Enter Vault username: " username; echo $$username)
VAULT_SECRET ?= passdb/velero

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

sed_replace = s/{{DIGITALOCEAN_API_TOKEN}}/${DIGITALOCEAN_API_TOKEN}/g;s/{{SPACES_ACCESS_KEY_ID}}/${SPACES_ACCESS_KEY_ID}/g;s/{{SPACES_SECRET_ACCESS_KEY}}/${SPACES_SECRET_ACCESS_KEY}/g

_deploy:
	velero install \
		--provider velero.io/aws \
		--bucket unchain-velero-backups \
		--plugins velero/velero-plugin-for-aws:v1.0.0,digitalocean/velero-plugin:v1.0.0 \
		--backup-location-config s3Url=https://ams3.digitaloceanspaces.com,region=ams3 \
		--use-volume-snapshots=false \
		--no-secret
	
	velero snapshot-location create default --provider digitalocean.com/velero
	cat ./01-velero-secret.yaml | sed "${sed_replace}" | kubectl apply -f -
	kubectl patch deployment/velero -p "$$(cat ./02-velero-deployment.patch.yaml)" --namespace velero

_check:
	velero install \
		--provider velero.io/aws \
		--bucket unchain-velero-backups \
		--plugins velero/velero-plugin-for-aws:v1.0.0,digitalocean/velero-plugin:v1.0.0 \
		--backup-location-config s3Url=https://ams3.digitaloceanspaces.com,region=ams3 \
		--use-volume-snapshots=false \
		--no-secret --dry-run -o yaml
		
	cat ./01-velero-secret.yaml | sed "${sed_replace}"
	kubectl patch deployment/velero -p "$$(cat ./02-velero-deployment.patch.yaml)" --namespace velero --dry-run -o yaml

test:
	velero version
