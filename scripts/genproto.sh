#!/usr/bin/env bash
# Generate Go code from .proto files. Run from repo root.
#
# Prerequisites:
#   - protoc: https://grpc.io/docs/protoc-installation/
#   - Go plugins: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#                 go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
#   - GOPATH/bin or GOBIN in PATH so protoc can find the plugins

set -e
cd "$(dirname "$0")/.."
# Use Go bin for protoc-gen-go and protoc-gen-go-grpc
export PATH="$(go env GOPATH)/bin:$PATH"

PROTO_ROOT=api/proto
GEN_DIR=api/gen

mkdir -p "$GEN_DIR"

protoc -I . \
  --go_out=. --go_opt=module=openswarm \
  --go-grpc_out=. --go-grpc_opt=module=openswarm \
  "$PROTO_ROOT"/types.proto \
  "$PROTO_ROOT"/worker.proto \
  "$PROTO_ROOT"/coordinator.proto

echo "Generated Go code in $GEN_DIR/"
