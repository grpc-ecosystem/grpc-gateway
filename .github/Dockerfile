FROM golang:1.25.0

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get -y install --no-install-recommends \
        wget \
        unzip \
        openjdk-17-jre \
        bzip2 \
        patch && \
    apt-get clean -y && \
    rm -rf /var/lib/apt/lists/*

# Install swagger-codegen
ENV SWAGGER_CODEGEN_VERSION=2.4.8
RUN wget https://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/${SWAGGER_CODEGEN_VERSION}/swagger-codegen-cli-${SWAGGER_CODEGEN_VERSION}.jar \
    -O /usr/local/bin/swagger-codegen-cli.jar && \
    echo '#!/bin/bash\njava -jar /usr/local/bin/swagger-codegen-cli.jar "$@"' > /usr/local/bin/swagger-codegen && \
	chmod +x /usr/local/bin/swagger-codegen

# Install Bazelisk as bazel to manage Bazel
RUN go install github.com/bazelbuild/bazelisk@latest && \
    mv $(which bazelisk) /usr/local/bin/bazel

# Install buildifier for bazel formatting
RUN go install github.com/bazelbuild/buildtools/buildifier@latest
