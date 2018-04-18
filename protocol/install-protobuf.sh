#!/bin/sh
set -e
# check to see if protobuf folder is empty
if [ ! -d "$HOME/protobuf/lib" ]; then
  wget https://github.com/google/protobuf/releases/download/v3.5.1/protobuf-all-3.5.1.tar.gz
  tar -xzvf protobuf-all-3.5.1.tar.gz
  cd protobuf-3.5.1 && ./configure --prefix=$HOME/protobuf && make && make install
else
  echo "Using cached directory."
fi
