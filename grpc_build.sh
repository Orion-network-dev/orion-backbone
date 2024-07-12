#!/bin/sh

for f in $(find proto -name 'orion_*.proto'); do
    protoc --go_out=internal --go_opt=paths=source_relative --go-grpc_out=internal --go-grpc_opt=paths=source_relative "$f"
done
