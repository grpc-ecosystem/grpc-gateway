#!/bin/sh -eu
codegen_version=$1
if test -z "${codegen_version}"; then
	echo "Usage: .travis/install-swagger-codegen.sh codegen-version"
	exit 1
fi

swagger_codegen_cli="/tmp/swagger-codegen-cli.jar"

# Want to test with an unreleased version of swagger-codegne-cli? Try out the sonatype repo for the SNAPSHOT builds.
if false; then
  # Directory listing can be found at
  # https://oss.sonatype.org/content/repositories/snapshots/io/swagger/swagger-codegen-cli/2.4.0-SNAPSHOT/
  # Replace the version number with the appropriate version.
  wget https://oss.sonatype.org/content/repositories/snapshots/io/swagger/swagger-codegen-cli/2.4.0-SNAPSHOT/swagger-codegen-cli-2.4.0-20180407.135302-217.jar \
    -O "${swagger_codegen_cli}"
else
  wget "http://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/${codegen_version}/swagger-codegen-cli-${codegen_version}.jar" \
    -O "${swagger_codegen_cli}"
fi
