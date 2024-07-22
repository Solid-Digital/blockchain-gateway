## This makefile manages configuration files
date := $(shell date)
tag := ${TAG}

info:
	@echo "Try: \n * make get-staging \n * make update "

check-env:
	@if [ -z ${env} ]; then echo "env [dev,staging,prod] is not set!"; exit 2; fi 

get-config: check-env
	kubectl get secret ares-config --namespace ares-${env} -o yaml --export -o json | jq -r '.data["config.toml"]' |base64 --decode > ${env}/config.toml

get-kubeconfig: check-env
	kubectl get secret kubeconfig --namespace ares-${env} -o yaml --export -o json | jq -r '.data["kubeconfig.yaml"]' |base64 --decode

set-kubeconfig: check-env
	kubectl create secret generic kubeconfig --namespace ares-${env} --from-file ${env}/kubeconfig.yaml

# Fetch config from server, and compare. Then save new version.
update: check-env
	@echo "getting secret from server"
	@kubectl get secret ares-config -o yaml -n ares-${env} --export -o json | jq -r '.data["config.toml"]' |base64 --decode > ${env}/remote_config.toml
	
	@if ! diff ${env}/remote_config.toml ${env}/config.toml; \
		then read -p "Do you want to continue? [y/n] " yn; \
		[ $$yn = "y" ] && make force-update env=${env}; \
		echo "all done"; \
	else echo "files are the same, no update required"; \
	fi
	@rm -f ${env}/remote_config.toml

# force update a config; takes env parameter
force-update: check-env
	@echo updating secret
	kubectl create secret generic ares-config --from-file ${env}/config.toml -n ares-${env} --dry-run -o yaml | kubectl apply -f -

rolling-update: check-env
	kubectl patch deployment ares -n ares-${env} --type='json' -p='[{"op": "replace", "path": "/spec/template/metadata/annotations/config-version", "value":"${date}"}]'
