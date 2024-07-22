VAULT_ADDR=https://vault.tools.dgo.unchain.io
VAULT_DB_ID=13415e2e-cda7-448a-86ed-ada12dace1bd
VAULT_DB_URL=$(shell doctl databases get ${VAULT_DB_ID} -o json | jq ".[0].connection.uri")
sed_replace=s/{{date}}/${date}/g

ldapconfig:
	vault write auth/ldap/config \
		url="ldaps://ldaps.unchain.io" \
		userdn="dc=unchain,dc=io" \
		groupdn="dc=unchain,dc=io" \
		groupfilter="(&(objectClass=group)(member:1.2.840.113556.1.4.1941:={{.UserDN}}))" \
		groupattr="cn" \
		upndomain="unchain.io" \
		certificate=@ldap/ldap.unchain.io.crt \
		insecure_tls=false \
		starttls=true

unseal:
	VAULT_ADDR=${VAULT_ADDR} vault operator unseal 

env:
	@echo VAULT_PG_CONNECTION_URL=${VAULT_DB_URL} > vault.env

deploy: env
	kubectl kustomize | sed "${sed_replace}" | kubectl apply -f -
	$(MAKE) clean

delete: env
	kubectl kustomize | kubectl delete -f -
	$(MAKE) clean

test:
	VAULT_ADDR=${VAULT_ADDR} vault kv get ops/test

login:
	VAULT_ADDR=${VAULT_ADDR} vault login -method ldap username=${VAULT_USER}

clean: 
	rm -rf vault.env
