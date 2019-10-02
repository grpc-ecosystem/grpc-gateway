FROM golang:1.13.1

# Warm apt cache
RUN apt-get update

# Install swagger-codegen
ENV SWAGGER_CODEGEN_VERSION=2.4.8
RUN apt-get install -y openjdk-11-jre wget && \
    wget http://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/${SWAGGER_CODEGEN_VERSION}/swagger-codegen-cli-${SWAGGER_CODEGEN_VERSION}.jar \
    -O /usr/local/bin/swagger-codegen-cli.jar && \
    apt-get remove -y wget

# Wrap the jar for swagger-codgen
RUN echo '#!/bin/bash\n java -jar /usr/local/bin/swagger-codegen-cli.jar "$@"' > /usr/local/bin/swagger-codegen && \
	chmod +x /usr/local/bin/swagger-codegen

# Install protoc
ENV PROTOC_VERSION=3.9.2
RUN apt-get install -y wget unzip && \
    wget https://github.com/google/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip \
    -O /protoc-${PROTOC_VERSION}-linux-x86_64.zip && \
    unzip /protoc-${PROTOC_VERSION}-linux-x86_64.zip -d /usr/local/ && \
    rm -f /protoc-${PROTOC_VERSION}-linux-x86_64.zip && \
    apt-get remove -y unzip wget

# Install node, used by NVM
ENV NODE_VERSION=v10.16.3
ENV NVM_VERSION=v0.35.0
RUN apt-get install -y wget bzip2 && \
    wget -qO- https://raw.githubusercontent.com/creationix/nvm/${NVM_VERSION}/install.sh | bash && \
    apt-get remove -y wget

# Clean up
RUN apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/*
