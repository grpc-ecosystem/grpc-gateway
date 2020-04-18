## gRPC-Gateway CI testing setup

Contained within is the CI test setup for the Gateway. It runs on Circle CI.

### Whats up with the Dockerfile?

The `Dockerfile` in this folder is used as the build environment when regenerating the files (see above).
The canonical repository for this Dockerfile is `jfbrandhorst/grpc-gateway-build-env`. Please request access
before attempting to make any changes to the Dockerfile.
