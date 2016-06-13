#!/bin/sh -eu
bazel_version=$1
if test -z "${bazel_version}"; then
	echo "Usage: .travis/install-bazel.sh bazel-version"
	exit 1
fi
if ! $HOME/local/bazel/bin/bazel version 2>/dev/null | grep "Build label ${bazel_version}"; then
	rm -rf $HOME/local/bazel
	wget "https://github.com/bazelbuild/bazel/releases/download/${bazel_version}/bazel-${bazel_version}-installer-linux-x86_64.sh"
	bash ./bazel-${bazel_version}-installer-linux-x86_64.sh \
		--prefix=$HOME/local/bazel \
		--bazelrc=$HOME/local/bazel/etc/bazelrc
fi

ln -s $HOME/local/bazel/etc/bazelrc $HOME/.bazelrc
echo \$ $HOME/local/bazel/bin/bazel version
$HOME/local/bazel/bin/bazel version

