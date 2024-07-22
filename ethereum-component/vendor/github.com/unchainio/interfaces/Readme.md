### Update protobufs:

Install protoc from https://github.com/protocolbuffers/protobuf/releases

Install protoc-gen-go: go get -u github.com/golang/protobuf/protoc-gen-go

`protoc -I adapter/proto/ adapter/proto/*.proto --go_out=plugins=grpc:adapter/proto/`

