---
layout: default
title: Generating stubs using buf
parent: Generating stubs
grand_parent: Tutorials
nav_order: 1
---

## Generating stubs using buf

Buf is configured through a `buf.yaml` file that should be checked in to the root of your repository. Buf will automatically read this file if present. Configuration can also be provided via the command-line flag `--config`, which accepts a path to a `.json` or `.yaml` file, or direct JSON or YAML data.

All Buf operations that use your local `.proto` files as input rely on a valid build configuration. This configuration tells Buf where to search for `.proto` files, and how to handle imports. As opposed to protoc, where all `.proto` files are manually specified on the command-line, buf operates by recursively discovering all `.proto` files under configuration and building them.

The following is an example of all configuration options for the build.

```yml
version: v1beta1
build:
  roots:
    - proto
    - vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis
  excludes:
    - proto/foo/bar
```

To generate stubs for for multiple languages. Create the file `buf.gen.yaml` at the root of the repository:

```yml
version: v1beta1
plugins:
  - name: java
    out: java
  - name: cpp
    out: cpp
```

Then run

```sh
buf generate --file ./proto/example.proto
```
