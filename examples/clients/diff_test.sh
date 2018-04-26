#!/bin/sh -e

while test $# -gt 0; do
  file=$1
  golden=$2
  shift; shift
  echo $file >&2
  diff -qN "$file" "$golden"
done
