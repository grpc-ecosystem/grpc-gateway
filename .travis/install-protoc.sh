#!/bin/sh -eu
protoc_version=$1
if test -z "${protoc_version}"; then
	echo "Usage: .travis/install-protoc.sh protoc-version"
	exit 1
fi

protoc_path="/tmp/proto"
protoc_binary="${protoc_path}/bin/protoc-${protoc_version}"

if [ "$("${protoc_binary}" --version 2>/dev/null | cut -d' ' -f 2)" != "${protoc_version}" ]; then
	rm -rf "${protoc_path:?}/bin" "${protoc_path:?}/include"

  tempdir=$(mktemp -d 2>/dev/null || mktemp -d -t 'protoc')

	mkdir -p "${protoc_path}"
  cd "${tempdir}"

	wget "https://github.com/google/protobuf/releases/download/v${protoc_version}/protoc-${protoc_version}-linux-x86_64.zip"
	unzip "protoc-${protoc_version}-linux-x86_64.zip"

	mv bin "${protoc_path}/bin"
	mv include "${protoc_path}/include"
fi

echo "\$ ${protoc_path}/bin/protoc --version"
"${protoc_path}/bin/protoc" --version
