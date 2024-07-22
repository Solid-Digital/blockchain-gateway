version := $(shell if [ ! -z "${TAG}" ]; then echo "${TAG}"; else git describe --tags --always; fi)
date := $(shell date)
export env := ${env}
tag := ${TAG}
branch := $(shell if [ ! -z "${BRANCH}" ]; then echo "${BRANCH}"; else git rev-parse --abbrev-ref HEAD; fi)
builder := $(shell if [ ! -z "${BUILDER}" ]; then echo "${BUILDER}"; elif [ ! -z ${BITBUCKET_BUILD_NUMBER} ]; then git log -1 --pretty=format:'%an' | xargs ; else git config user.name; fi)
ldflags := "-X 'main.version=${version}' -X 'main.branch=${branch}' -X 'main.builder=${builder}' -X 'main.buildDate=${date}'"

build:
	CGO_ENABLED=0 go build -ldflags ${ldflags} -a -installsuffix cgo -o bin/tbg-nodes-auth ./cmd/tbg-nodes-auth/...

image:
	docker build -t registry.unchain.io/unchain/tbg-nodes-auth:${version} -t registry.unchain.io/unchain/tbg-nodes-auth:latest .

push:
	docker push registry.unchain.io/unchain/tbg-nodes-auth:${version}

check-env:
	@if [ -z ${env} ]; then echo "env [dev,staging,prod] is not set!"; exit 2; fi

deployment: check-env
	sed "s/{{tag}}/${version}/g;s/{{date}}/${date}/g" deployments/${env}/deployment.template.yaml > deployments/${env}/deployment.yaml
	kubectl apply -f deployments/${env}/deployment.yaml -n tbg-nodes-${env}
	@if [ ${env} = production ]; then \
		echo "succesfully deployed"; \
	else \
		echo "succesfully deployed"; \
	fi

quick-deployment:
	env=$(shell if [ ! -z "${env}" ]; then echo "${env}"; else echo "staging"; fi) TAG=${version}-$(shell date '+%s') make build image push deployment


install-tbg-nodes-auth:
	go install -ldflags ${ldflags} ./cmd/tbg-nodes-auth/...


gen-orm:
	eval "$$(go env)"; \
	 sqlboiler psql --wipe \
     		--struct-tag-casing camel \
     		--no-tests

gen-wire:
	wire ./gen/wire
