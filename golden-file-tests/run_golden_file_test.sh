#!/usr/bin/env bash
cd cx-api && bash build_facade_files.sh
cd ..

if diff <(jq --sort-keys . apidocs_golden.swagger.json) <(jq --sort-keys . apidocs.swagger.json) > /dev/null; then
  echo "Plugin output matches golden file."
else
  echo "Plugin output does not match golden file."
  exit 1
fi

rm -r apidocs.swagger.json