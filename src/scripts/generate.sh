#!/bin/bash

PROTO_DIR=./common/protobuf
GENERATED_DIR=./common/generated

mkdir -p $GENERATED_DIR

protoc \
    --go_out=$GENERATED_DIR \
    --go-grpc_out=$GENERATED_DIR \
    --go_opt=paths=import \
    --go-grpc_opt=paths=import \
    $PROTO_DIR/*.proto