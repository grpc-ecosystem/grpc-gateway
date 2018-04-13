#!/bin/sh
set -e
wget http://central.maven.org/maven2/com/google/api/grpc/googleapis-common-protos/0.0.3/googleapis-common-protos-0.0.3.jar
jar xvf googleapis-common-protos-0.0.3.jar
cp -r google/ $HOME/protobuf/include/
ls -l



