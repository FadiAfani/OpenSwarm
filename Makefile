.PHONY: deps proto gen build run coordinator cli

COORDINATOR_CMD := ./cmd/coordinator
COORDINATOR_BIN := ./bin/coordinator
CLI_CMD := ./cli
CLI_BIN := ./bin/openswarm-cli
CLI_ROOT_BIN := ./openswarm-cli

# Install protoc Go plugins (ensure $GOPATH/bin or $GOBIN is in PATH)
deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go mod tidy

# Generate Go code from .proto files. Requires protoc and the Go plugins.
proto: gen
gen:
	./scripts/genproto.sh

# Build the coordinator binary.
build: coordinator

coordinator:
	go build -o $(COORDINATOR_BIN) $(COORDINATOR_CMD)

# Build the CLI binary.
cli:
	go build -o $(CLI_BIN) $(CLI_CMD)
	cp $(CLI_BIN) $(CLI_ROOT_BIN)

# Run the coordinator service directly.
run:
	go run $(COORDINATOR_CMD)
