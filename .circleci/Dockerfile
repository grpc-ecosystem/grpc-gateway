FROM golang:1.13.4

# Warm apt cache and install dependencies
# bzip2 is required by the node_tests (to extract its dependencies).
RUN apt-get update && \
    apt-get install -y wget unzip \
    openjdk-11-jre \
    bzip2

# Install swagger-codegen
ENV SWAGGER_CODEGEN_VERSION=2.4.8
RUN wget http://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/${SWAGGER_CODEGEN_VERSION}/swagger-codegen-cli-${SWAGGER_CODEGEN_VERSION}.jar \
    -O /usr/local/bin/swagger-codegen-cli.jar

# Wrap the jar for swagger-codgen
RUN echo -e '#!/bin/bash\njava -jar /usr/local/bin/swagger-codegen-cli.jar "$@"' > /usr/local/bin/swagger-codegen && \
	chmod +x /usr/local/bin/swagger-codegen

# Install protoc
ENV PROTOC_VERSION=3.10.1
RUN wget https://github.com/google/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip \
    -O /protoc-${PROTOC_VERSION}-linux-x86_64.zip && \
    unzip /protoc-${PROTOC_VERSION}-linux-x86_64.zip -d /usr/local/ && \
    rm -f /protoc-${PROTOC_VERSION}-linux-x86_64.zip

# Install node, used by NVM
ENV NODE_VERSION=v10.16.3
ENV NVM_VERSION=v0.35.0
RUN wget -qO- https://raw.githubusercontent.com/creationix/nvm/${NVM_VERSION}/install.sh | bash

# Clean up
RUN apt-get autoremove -y && \
    apt-get remove -y wget \
    unzip && \
    rm -rf /var/lib/apt/lists/*
