
# Important!
# Version is either set to the the TAG env var. (when running in Docker), or the git tag.
version := $(shell if [ ! -z "${TAG}" ]; then echo "${TAG}"; else git describe --tags --always; fi)
date := $(shell date)
export env := ${env}
tag := ${TAG}
branch := $(shell if [ ! -z "${BRANCH}" ]; then echo "${BRANCH}"; else git rev-parse --abbrev-ref HEAD; fi)
builder := $(shell if [ ! -z "${BUILDER}" ]; then echo "${BUILDER}"; elif [ ! -z ${BITBUCKET_BUILD_NUMBER} ]; then git log -1 --pretty=format:'%an' | xargs ; else git config user.name; fi)
ldflags := "-X 'main.version=${version}' -X 'main.branch=${branch}' -X 'main.builder=${builder}' -X 'main.buildDate=${date}'"

.PHONY: install db db-start test testall clean up up-bitbucket db-migrate k3s-clean gen req restart

version:
	@echo ${version}

req:
	go get -v -u github.com/oxequa/realize/...
	go get -v -u github.com/google/wire/...
	go get -v -u github.com/go-swagger/go-swagger/...
	go get -v -u github.com/volatiletech/sqlboiler/...
	go get -v -u github.com/golang/mock/mockgen/...

up:
	docker-compose up -d

reset: clean db-fixtures

testall: clean test

clean: install-migrate install-fixtures install-bootstrap restart wait db-migrate bootstrap

restart:
	docker-compose rm --stop --force -v
#	docker run --rm --volume "${PWD}:/repo" alpine rm -rf /repo/data
#	sudo rm -rf data
	docker-compose up -d

test:
	go test -count=1 -cover -p=1 ./...

test-short:
	go test -count=1 -cover -p=1 -test.short ./...

test-bitbucket:
	go test -cover -p=1 --short ./...

up-bitbucket:
	./tools/docker-compose up -d postgres mail redis s3

db-reset: db-clean gen-orm db-fixtures

db-clean: install-migrate install-fixtures install-bootstrap db-recreate wait db-migrate bootstrap

install-bootstrap:
	go install -ldflags ${ldflags} ./cmd/ares-bootstrap/...

bootstrap:
	ares-bootstrap

db-migrate:
	ares-migrate up

db-fixtures:
	ares-fixtures up

db-dump:
	PGPASSWORD="mysecretpassword" pg_dump --dbname="ares_dev" --schema='public' --table=public.action --table=public.pipeline --table=public."organization_environment" --table=public."trigger_version_supported_bases" --table=public."trigger_draft_configuration" --table=public."action_version" --table=public.organization --table=public."base_version" --table=public."trigger_configuration" --table=public."draft_configuration" --table=public."base_configuration" --table=public.environment --table=public."trigger_version" --table=public."base_draft_configuration" --table=public."action_configuration" --table=public.users --table=public."user_organization" --table=public.configuration --table=public."action_draft_configuration" --table=public.base --table=public.deployment --table=public."action_version_supported_bases" --table=public.trigger --file="${PWD}/dump.sql" --data-only --column-inserts --format=p --username=postgres --host=localhost --port=5432

db-recreate:
	docker-compose rm --stop --force -v postgres
	docker-compose up -d postgres

gen: gen-api gen-client gen-orm gen-wire gen-mocks

gen-go:
	go generate ./...

gen-wire:
	wire ./gen/wire

gen-orm:
	eval "$$(go env)"; \
	 sqlboiler psql --wipe \
     		--struct-tag-casing camel \
     		--templates "./templates/sqlboiler" \
     		--no-tests

gen-api:
	$(MAKE) gen-api-spec
	mkdir -p gen
	swagger generate server -f ./api/swagger.yaml -t ./gen -C go-swagger.yaml --exclude-main -P dto.User --server-package api --model-package dto
	$(MAKE) sync-handlers

sync-handlers:
## assuming that the current stuff are already synced
# replace old stuff with current stuff
# re-generate current stuff
# sync
	sync-handlers ./gen/api/handlers ./pkg/http

gen-api-spec:
	-rm ./api/swagger/swagger.yaml
	swagger mixin -o ./api/swagger/swagger.yaml --format yaml ./api/swagger/*.yaml -c 5
	swagger flatten -o ./api/swagger.yaml --format=yaml ./api/swagger/swagger.yaml

gen-client:
	$(MAKE) gen-api-spec
	mkdir -p gen
	swagger generate client -f ./api/swagger.yaml -t ./gen --model-package dto

gen-mocks:
	mockgen -source=./pkg/ares/auth.go -destination=./gen/mocks/auth.go
	mockgen -source=./pkg/ares/component.go -destination=./gen/mocks/component.go
	mockgen -source=./pkg/ares/organization.go -destination=./gen/mocks/organization.go
	mockgen -source=./pkg/ares/pipeline.go -destination=./gen/mocks/pipeline.go
	mockgen -source=./pkg/ares/deployment.go -destination=./gen/mocks/deployment.go
	mockgen -source=./pkg/ares/aws.go -destination=./gen/mocks/aws.go
	mockgen -source=./pkg/ares/subscription.go -destination=./gen/mocks/subscription.go

wait:
	sleep 3

wait-more:
	sleep 5

k3s-clean:
	docker-compose rm --stop --force -v k3s-server k3s-node registry elastic
	docker-compose up -d k3s-server k3s-node registry elastic

## Build steps
run: 
	docker run -P -it ares 

docker-login:
	docker login --username ${DOCKER_USERNAME} --password ${DOCKER_PASSWORD} ${DOCKER_REGISTRY}

build:
	CGO_ENABLED=0 go build -ldflags ${ldflags} -a -installsuffix cgo -o bin/ares ./cmd/ares/...
	CGO_ENABLED=0 go build -ldflags ${ldflags} -a -installsuffix cgo -o bin/ares-migrate ./cmd/ares-migrate/...
	CGO_ENABLED=0 go build -ldflags ${ldflags} -a -installsuffix cgo -o bin/ares-fixtures ./cmd/ares-fixtures/...
	CGO_ENABLED=0 go build -ldflags ${ldflags} -a -installsuffix cgo -o bin/ares-bootstrap ./cmd/ares-bootstrap/...

install: install-ares install-fixtures install-migrate install-bootstrap

install-migrate:
	go install -ldflags ${ldflags} ./cmd/ares-migrate/...

install-fixtures:
	go install -ldflags ${ldflags} ./cmd/ares-fixtures/...

install-ares:
	go install -ldflags ${ldflags} ./cmd/ares/...

image:
	docker build -t registry.unchain.io/unchain/ares:${version} .

full-image:
	docker build -f Dockerfile.complete --build-arg TAG="${version}" --build-arg  BRANCH="${branch}" --build-arg BUILDER="${builder}" -t registry.unchain.io/unchain/ares:${version} .

push:
	docker push registry.unchain.io/unchain/ares:${version}

deployment:
	cd deployments && $(MAKE) deployment

kubelogs:
	kubectl logs -c ares -l app.kubernetes.io/name=ares

exec:
	@>&2 echo "creating exec command.."
	@kubectl get pods --no-headers -o custom-columns=":metadata.name" -l app.kubernetes.io/name=ares \
	| head -n 1 \
	| xargs -I % echo -n kubectl exec -it -c ares % bash
