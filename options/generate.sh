#!/usr/bin/env bash
protoc --go_out=Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:. middleware.proto