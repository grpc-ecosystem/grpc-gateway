#!/bin/sh
set -e

go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go

wget http://central.maven.org/maven2/com/google/api/grpc/googleapis-common-protos/0.0.3/googleapis-common-protos-0.0.3.jar
jar xvf googleapis-common-protos-0.0.3.jar
cp -r google/ $HOME/protobuf/include/
ls -l



