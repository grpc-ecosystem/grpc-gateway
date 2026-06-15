---
layout: default
title: Merging OpenAPI 3.1 Output
nav_order: 8
parent: Mapping
---

# Merging OpenAPI 3.1 Output

{: .note }
This page describes `openapiv3-merge`, the post-processing tool that
combines `protoc-gen-openapiv3` output into a single OpenAPI document. The
generator itself is documented under
[OpenAPI 3.1 Output](./openapi_v3.md).

[`protoc-gen-openapiv3`](./openapi_v3.md) emits **one OpenAPI document per
proto file**, matching the protobuf convention of 1:1 input-to-output
files. Many OpenAPI consumers — UI viewers, client generators, gateway
configurations — expect a single document instead. `openapiv3-merge` is the
post-processing step that combines the per-file outputs into one.

## Installation

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
document, so put the file with the broadest scope (the one carrying any
shared [`openapiv3_document`](./openapi_v3.md#example) annotation) first.

### With buf

Generate per-proto-file documents the usual way, then merge them:

```sh
#!/usr/bin/env bash
set -euo pipefail

buf generate

# Put the file with the shared openapiv3_document annotation first so
# its info/servers/externalDocs are kept, then everything else in
# sorted order for determinism.
root=gen/api/v1/api.openapi.json
mapfile -t rest < <(find gen -name '*.openapi.json' ! -path "$root" | sort)
openapiv3-merge "$root" "${rest[@]}" > gen/api.openapi.json
```

If you don't have a designated root file, just sort all inputs and accept
whichever sorts first — but set the `openapiv3_document` annotation on
that proto so the merged output has a sensible `title`, `version`, and
`servers`.

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

## What gets merged, and how

| Field | Rule |
| --- | --- |
| `openapi` | Must match across all inputs. |
| `info`, `servers`, `externalDocs` | Taken from the first input. Later inputs' values are dropped. |
| `paths`, `webhooks` | Unioned in input order. Path collisions with non-identical content are errors. |
| `components/*` | Unioned with sorted keys. Same name with non-identical content is an error. |
| `tags` | Deduplicated by `name`. Same name with non-identical metadata is an error. |
| `security` | First input to declare a non-empty array wins. Later inputs that declare a different non-empty array are an error. |
| Unknown top-level keys (extensions, `x-*`) | Taken from the first input. |

"Non-identical content" is determined by canonical comparison: two values
that differ only in key order compare equal. Since `protoc-gen-openapiv3`
names component schemas with their fully qualified proto name (e.g.
`example.v1.User`), the same message reused across packages contributes a
single, identical component to every output, and the merger combines them
cleanly.

The first input controls the merged document's `info` block. If you want a
shared title, version, contact, license, or set of `servers` across all
merged output, set them with an
[`openapiv3_document`](./openapi_v3.md#example) annotation on the
proto file that appears first in your merge invocation.

## Output

- Top-level fields are emitted in OpenAPI 3.1.0 declaration order
  (`openapi`, `info`, `servers`, `paths`, `webhooks`, `components`,
  `security`, `tags`, `externalDocs`, then any extensions).
- `paths` and `webhooks` are emitted in input order.
- `components/*` sub-maps are sorted lexicographically, matching what
  `protoc-gen-openapiv3` itself emits per file.
- `tags` are emitted in first-occurrence order across inputs.

## Example merge errors

When a merge fails, the error message identifies the conflicting field:

```
openapiv3-merge: paths."/v1/echo": echo.openapi.json redefines an entry with a different value
```

Common causes:

- **Conflicting path definitions.** Two protos declared the same HTTP path
  with different bindings. Decide which one is authoritative.
- **Conflicting component definitions.** Two messages with the same fully
  qualified proto name but different shapes — usually means a duplicate
  `package` + message declaration somewhere in your proto tree.
- **Conflicting tag metadata.** Two `openapiv3_document.tags` annotations
  described the same tag with different descriptions. Make them match or
  move the metadata to one place.
- **Mismatched OpenAPI versions.** All inputs must declare the same
  `openapi` value.

## See also

- [OpenAPI 3.1 Output](./openapi_v3.md) — the generator that produces the
  per-file inputs to `openapiv3-merge`.
- [openapiv3-merge README](https://github.com/grpc-ecosystem/grpc-gateway/tree/main/openapiv3-merge) — quick reference and source.
