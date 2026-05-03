FROM golang:1.26.2

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get -y install --no-install-recommends \
        jq && \
    apt-get clean -y && \
    rm -rf /var/lib/apt/lists/*

# Install Bazelisk as bazel to manage Bazel
RUN go install github.com/bazelbuild/bazelisk@latest && \
    mv $(which bazelisk) /usr/local/bin/bazel

# Install buildifier for bazel formatting
RUN go install github.com/bazelbuild/buildtools/buildifier@latest
