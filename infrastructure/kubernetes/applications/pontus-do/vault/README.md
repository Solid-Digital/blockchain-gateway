# Vault

# Prerequisites
1. kubectl - https://kubernetes.io/docs/tasks/tools/install-kubectl/
2. jq - https://stedolan.github.io/jq/download/
3. doctl - https://github.com/digitalocean/doctl#installing-doctl
   - run `doctl auth init` to login
4. vault - https://learn.hashicorp.com/vault/getting-started/install
5. Access to the unseal keys which is provided on a need to have basis

# Apply

Apply via `make apply`

# Delete

Delete via `make delete`

# Test

Test via: 
 - `make login VAULT_USER=<ldap-username-without-the-@unchain.io-part>`
 - `make test`
# Upgrade

1. Create a fork of the vault storage database via the DigitalOcean interface - https://cloud.digitalocean.com/databases/vault-postgres-db?i=977411
2. Test if the fork works
   - run `make apply VAULT_DB_URL='<connection-url-of-the-fork>` or `make apply VAULT_DB_ID='<database-id-of-the-fork>'`
   - run `make unseal` for each unseal key
   - run `make login VAULT_USER=<ldap-username-without-the-@unchain.io-part>`
   - run `make test`, which should result in a success
3. Update the vault version
   - update the image version in deployment.yaml
   - run `make apply`
   - run `make unseal` for each unseal key
   - run `make login VAULT_USER=<ldap-username-without-the-@unchain.io-part>`
   - run `make test` again to verify if the new version of vault is working
4. If the upgrade was successful, delete the fork of the database