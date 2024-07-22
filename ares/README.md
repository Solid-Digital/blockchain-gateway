# README #

## Running

Directly:

`go run cmd/ares/main.go`

With auto-recompile upon changes:

`docker-compose up`

or 

`realize start`

## Dependency management:

We use Go modules for dependency management. Those are built into and turned on by default for go 1.13 and above.
`go run/build/install/test` will now automatically download all necessary dependencies and therefore the vendor directory is no longer necessary.
Additionally the GOPATH is no longer mandatory and go projects can be cloned anywhere on the filesystem.

Some commands:

`go mod tidy` - makes sure go.mod matches the source code in the module.

`go mod download` - explicitly downloads the necessary dependencies. Usually not useful because `go run/build/install/test` will already do this.

### Go modules integration in Goland

Enable the following settings when imports in Goland don't work (after downloading dependencies):

`Settings | Go | Go modules (vgo) | Enable Go Modules (vgo) integration`

## Get all required tools:

To download realize, wire, go-swagger, sqlboiler, mockery, run:

`make req`

## Generating dependency injection code:

We use wire to generate the dependency injection code. Get it via:

`go get github.com/google/wire/...`

To generate the code:

`wire ./pkg/ares`

## Generate the api from swagger:

`make gen-api` 

The minimum version of Go is 1.11 otherwise it won't work (you will get some 
error in a template complaining about the '=') sign.

## Running tests:

`make testall`

## Generating mocks

We use mockery to generate mocks. This will generate mocks for all interfaces
in the `./pkg/domain` directory:

`make gen-mocks`

## Troubleshooting:

1. If elasticsearch doesn't start, run `sudo sysctl -w vm.max_map_count=262144`