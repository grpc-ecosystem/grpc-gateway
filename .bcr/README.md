# Configuration to use Publish to BCR

This directory contains a set of templates required by the [publish-to-bcr](https://github.com/bazel-contrib/publish-to-bcr/tree/main) plugin, which will publish a released version of **gRPC-Gateway** as a new module in the [Bazel Central Registry](https://registry.bazel.build/).

The plugin aims to eliminate a manual publishing process for the module and initiate the process when a release is created.

## Publish to BCR template files

The configuration consists of three files placed in the `.bcr` directory:

* `.bcr/metadata.template.json`: that describes the repository and maintainers' information.
* `.bcr/presubmit.yml`: describes the targets that will be built and tested on specific platforms and bazel versions to test the module.
* `.bcr/source.template.json`: that will automatically substitute values for the repository, owner, and tag based on the repository and release data.

_For more information regarding the files that form a BCR entry, check the following references:_

* [Bazel registries](https://bazel.build/external/registry).
* [External dependencies overview](https://bazel.build/external/overview).

## Result

The final result of this process is the creation of a PR in the BCR repository to publish the released version.

Once these templates are populated, the `publish-to-bcr` app should be configured as described [here](https://github.com/bazel-contrib/publish-to-bcr/tree/main?tab=readme-ov-file#how-it-works).
