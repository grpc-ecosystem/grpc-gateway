#!/bin/bash
set -xe

# Go-fuzz doesn't support modules yet, so ensure we do everything in the old style GOPATH way
export GO111MODULE="off"

# Install go-fuzz
go get -u github.com/dvyukov/go-fuzz/go-fuzz github.com/dvyukov/go-fuzz/go-fuzz-build

# Compiling fuzz targets in fuzz.go with go-fuzz (https://github.com/dvyukov/go-fuzz) and libFuzzer support
# This is a workaround until go-fuzz has gomodules support https://github.com/dvyukov/go-fuzz/issues/195
BRANCH=$(git rev-parse --abbrev-ref HEAD)
git branch --set-upstream-to=origin/master $BRANCH

go get -v -u ./protoc-gen-grpc-gateway/httprule
go-fuzz-build -libfuzzer -o parse-http-rule.a ./protoc-gen-grpc-gateway/httprule
clang-9 -fsanitize=fuzzer parse-http-rule.a -o parse-http-rule

wget -q -O fuzzit https://github.com/fuzzitdev/fuzzit/releases/download/v2.4.29/fuzzit_Linux_x86_64
chmod a+x fuzzit

if [ -z "CIRCLE_PULL_REQUEST" ]; then
    TYPE="fuzzing"
else
    TYPE="local-regression"
fi
./fuzzit create job --type ${TYPE} grpc-gateway/parse-http-rule parse-http-rule
