#!/bin/sh

mkdir -p internal/proto/
protoc --go_out=internal --go_opt=paths=source_relative --go-grpc_out=internal --go-grpc_opt=paths=source_relative proto/orion_rpc.proto
