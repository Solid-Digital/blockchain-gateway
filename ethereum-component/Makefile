start:
	docker-compose up -d

stop:
	docker-compose rm --stop --force -v

test: stop start sleep run-tests

run-tests:
	go test -p 1 ./... -count=1

sleep:
	sleep 2

build:
	mkdir -p bin/
	mkdir -p cache/go-build/

	docker run --rm --user ${UID}:${GID} -i -v "$$(pwd)/cache:/root/.cache" -v ${PWD}:/go/src/bitbucket.org/unchain/ethereum2 -w '/go/src/bitbucket.org/unchain/ethereum2' unchain/golang-libtdl:latest go build -buildmode=plugin -o bin/action.so ./cmd/action/
	docker run --rm --user ${UID}:${GID} -i -v "$$(pwd)/cache:/root/.cache" -v ${PWD}:/go/src/bitbucket.org/unchain/ethereum2 -w '/go/src/bitbucket.org/unchain/ethereum2' unchain/golang-libtdl:latest go build -buildmode=plugin -o bin/trigger.so ./cmd/trigger/
