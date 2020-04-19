# How to contribute

Thank you for your contribution to grpc-gateway.
Here's the recommended process of contribution.

1. `go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway`
1. `cd $GOPATH/src/github.com/grpc-ecosystem/grpc-gateway`
1. hack, hack, hack...
1. Make sure that your change follows best practices in Go
   - [Effective Go](https://golang.org/doc/effective_go.html)
   - [Go Code Review Comments](https://golang.org/wiki/CodeReviewComments)
1. Make sure that `go test ./...` passes.
1. Sign [a Contributor License Agreement](https://cla.developers.google.com/clas)
1. Open a pull request in Github

When you work on a larger contribution, it is also recommended that you get in touch
with us through the issue tracker.

### Code reviews

All submissions, including submissions by project members, require review.

### I want to regenerate the files after making changes!

Great, it should be as simple as thus (run from the root of the directory):

```bash
docker run -v $(pwd):/src/grpc-gateway --rm jfbrandhorst/grpc-gateway-build-env:1.14 \
    /bin/bash -c 'cd /src/grpc-gateway && \
        make realclean && \
        make examples && \
        make testproto'
docker run -itv $(pwd):/grpc-gateway -w /grpc-gateway --entrypoint /bin/bash --rm \
    l.gcr.io/google/bazel -c 'bazel run :gazelle -- update-repos -from_file=go.mod -to_macro=repositories.bzl%go_repositories; bazel run :buildifier'
docker run -itv $(pwd):/grpc-gateway -w /grpc-gateway --entrypoint /bin/bash --rm \
    l.gcr.io/google/bazel -c 'bazel run :gazelle'
```

If this has resulted in some file changes in the repo, please ensure you check those in with your merge request.

### Making a release

To make a release, follow these steps:

1. Decide on a release version. The `gorelease` job can
    recommend whether the new release should be a patch or minor release.
    See [CircleCI](https://app.circleci.com/pipelines/github/grpc-ecosystem/grpc-gateway/126/workflows/255a8a04-de9c-46a9-a66b-f107d2b39439/jobs/6428)
    for an example.
1. Generate a Github token with `repo` access.
1. Create a new branch and edit the Makefile `changelog` job, settings
    the `future-release=` variable to the name of the version you plan to release
1. Run `CHANGELOG_GITHUB_TOKEN=<yourtoken> make changelog`
1. Commit the `Makefile` and `CHANGELOG.md` changes.
1. Open a PR and check that everything looks right.
1. Merge the PR.
1. Tag the release on `master`, the tag should be made against the commit you just merged.
1. (Optional) Delete your Github token again.
1. (Required) Sit back and pat yourself on the back for a job well done :clap:.
