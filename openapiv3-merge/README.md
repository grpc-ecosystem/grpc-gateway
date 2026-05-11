# openapiv3-merge

`openapiv3-merge` combines multiple OpenAPI 3.1 JSON documents into a single
document. It is the companion tool for [`protoc-gen-openapiv3`](../protoc-gen-openapiv3/),
which emits one document per proto file. Pipe its output through
`openapiv3-merge` to produce a single combined spec.

## Why a separate binary

`protoc-gen-openapiv3` follows the protobuf convention of one output file per
input file. That keeps the plugin well-behaved under all protoc/buf
invocations and avoids forcing `buf` users onto `strategy: all`. The trade-off
is that consumers who want one document have to combine the per-file outputs
themselves. `openapiv3-merge` is that step.

## Install

```sh
go install github.com/grpc-ecosystem/grpc-gateway/v2/openapiv3-merge@latest
```

## Usage

```sh
openapiv3-merge FILE [FILE ...] > merged.openapi.json
```

The merged document is written to stdout. Errors go to stderr and the
process exits with a non-zero status. Input order matters: the first
input's `info`, `servers`, and `externalDocs` are kept in the merged
document, so put the file with the broadest scope first.

### With buf

```sh
#!/usr/bin/env bash
set -euo pipefail

buf generate

# Put the file with the shared openapiv3_document annotation first so
# its info/servers/externalDocs are kept; sort the rest for determinism.
root=gen/api/v1/api.openapi.json
mapfile -t rest < <(find gen -name '*.openapi.json' ! -path "$root" | sort)
openapiv3-merge "$root" "${rest[@]}" > gen/api.openapi.json
```

### With protoc

```sh
#!/usr/bin/env bash
set -euo pipefail

protoc -I. \
  --openapiv3_out=./gen \
  $(find . -name '*.proto')

root=gen/api/v1/api.openapi.json
mapfile -t rest < <(find gen -name '*.openapi.json' ! -path "$root" | sort)
openapiv3-merge "$root" "${rest[@]}" > gen/api.openapi.json
```

### One-shot

For trivial trees where any input may go first:

```sh
openapiv3-merge $(find . -name '*.openapi.json' | sort) > api.openapi.json
```

## Merge rules

The merger is strict: anything that could be a silent overwrite is rejected.

| Field | Rule |
| --- | --- |
| `openapi` | Must match across all inputs. |
| `info`, `servers`, `externalDocs` | Taken from the first input. Later inputs' values are dropped. |
| `paths`, `webhooks` | Unioned in input order. Same path with non-identical content → error. |
| `components/*` | Unioned with sorted keys. Same name with non-identical content → error. |
| `tags` | Deduplicated by `name`. Same name with non-identical metadata → error. |
| `security` | First input to declare a non-empty array wins. Later inputs that declare a different non-empty array → error. |
| Unknown top-level keys (extensions, `x-*`) | Taken from the first input. |

"Non-identical content" is compared canonically: two values that differ only
in key order are equal. This means component schemas with the same fully
qualified name (which is what `protoc-gen-openapiv3` emits — e.g.
`example.v1.User`) merge cleanly across packages.

## Output

- Top-level fields appear in OpenAPI 3.1.0 declaration order.
- `paths` and `webhooks` are emitted in input order.
- `components/*` sub-maps are sorted by key, matching `protoc-gen-openapiv3`.
- Tags appear in first-occurrence order.
