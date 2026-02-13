.PHONY: deps proto gen

# Install protoc Go plugins (ensure $GOPATH/bin or $GOBIN is in PATH)
deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go mod tidy

# Generate Go code from .proto files. Requires protoc and the Go plugins.
proto: gen
gen:
	./scripts/genproto.sh
