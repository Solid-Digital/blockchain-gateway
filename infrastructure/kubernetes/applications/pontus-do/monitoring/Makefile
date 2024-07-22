SHELL=bash

PRUNE_ARGS:=--prune -l=app.kubernetes.io/managed-by=kustomize,infrastructure=monitoring
NS_ARG:=--namespace=monitoring
DRY_RUN_ARG=--dry-run=server

VAULT_ADDR:=https://vault.tools.unchain.io

secrets:=$(subst generate_,, $(wildcard secrets/generate_*))

deploy: $(secrets)
	kubectl kustomize ./ | kubectl apply $(PRUNE_ARGS) $(NS_ARG) -f -

check: $(secrets)
	kubectl kustomize ./ | kubectl apply $(PRUNE_ARGS) $(NS_ARG) $(DRY_RUN_ARG) -f -

build: $(secrets)
	kubectl kustomize ./

# Generate secrets. https://www.gnu.org/software/make/manual/html_node/Automatic-Variables.html
%.secret.env:
	export VAULT_ADDR=$(VAULT_ADDR) && \
		./$(@D)/generate_$(@F) $@

%.secret.yaml:
	export VAULT_ADDR=$(VAULT_ADDR) && \
		./$(@D)/generate_$(@F) $@

clean:
	rm $(secrets)
