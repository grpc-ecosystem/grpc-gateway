---
layout: default
title: Security considerations
nav_order: 6
parent: Operations
---

# Security considerations

Things worth knowing about the defaults before you put a gateway in front of a
service that matters.

## Unknown JSON fields are dropped, and nothing says so

By default the gateway unmarshals JSON with `protojson` configured as
`DiscardUnknown: true`. A field that the target message does not define is
removed during transcoding and the request continues as if it had never been
sent.

Given this service:

```proto
message Command {
  string action = 1;
  int32  amount = 2;
}

service Gateway {
  rpc Send(Command) returns (Ack) {
    option (google.api.http) = {
      post: "/v1/send"
      body: "*"
    };
  }
}
```

a request carrying a third field:

```
POST /v1/send  {"action":"transfer","amount":1,"auth":"..."}
```

reaches the backend as `action` and `amount` only. The `auth` field is gone,
the backend sees a well formed message, and the client gets `200 OK`. Neither
end is told that anything was removed.

This is the behaviour the proto3 JSON mapping describes, and other JSON to gRPC
transcoders default the same way, so a client that works against one of them
will behave the same here. It is worth stating plainly because it is more
lenient than the `protojson` library the gateway wraps, which rejects unknown
fields by default.

### When this matters

The case to watch for is a client that puts something in the request expecting
the server to act on it, where the field is not in the schema. A token added by
an intermediary, a tenant identifier, a signature over the body, a flag a newer
client sends to an older backend. All of these disappear quietly, and the
request still succeeds, so the failure does not look like a failure.

It also matters when the schema changes. Removing or renaming a field does not
break clients that still send it. That is often the point, but it means a
rename can silently stop delivering data that the backend used to receive.

### Rejecting unknown fields instead

Set `DiscardUnknown` to false on the marshaler:

```go
mux := runtime.NewServeMux(
	runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: false,
		},
	}),
)
```

The same request now returns `400` with `unknown field "auth"`.

This is a good default for an internal API where every client is yours and a
stray field means someone made a mistake. It is a poor default for a public API
where clients you do not control may send extra fields for reasons of their
own, since it turns forward compatibility into an error.

### Noticing without rejecting

If rejecting is too strict but silence is too quiet, you can unmarshal twice and
log the difference: once with `DiscardUnknown: false` to see whether the strict
parse would have failed, then with the permissive setting for the request that
actually proceeds. That costs an extra parse per request, so it suits a
canary or a sampled fraction of traffic rather than everything.

## Request headers become gRPC metadata, under a prefix

The default incoming header matcher forwards the permanent HTTP headers
(`Authorization`, `Cookie`, `User-Agent` and so on) to the backend as gRPC
metadata under the `grpcgateway-` prefix, plus `Authorization` under the bare
`authorization` key. Everything else is dropped. A backend can therefore tell a
value that came from the HTTP client (`grpcgateway-cookie`) apart from one it
set itself.

Before this default changed, any header of the form `Grpc-Metadata-<key>` was
also forwarded as the bare key `<key>`. That let a client choose the metadata
key, so a backend that read, say, `md["user"]` expecting an interceptor to have
put it there could be handed a value from `Grpc-Metadata-User` instead.
[CVE-2026-18972](https://docs.velociraptor.app/announcements/advisories/cve-2026-18972/)
was exactly this.

If you rely on `Grpc-Metadata-` forwarding, `WithLegacyHeaderMatching` restores
it:

```go
mux := runtime.NewServeMux(
	runtime.WithLegacyHeaderMatching(),
)
```

Before doing so, check that nothing on the gRPC side trusts a metadata key the
client could now set. The safer option is a custom matcher that forwards the
specific headers you need, see
[Mapping from HTTP request headers to gRPC client metadata](../mapping/customizing_your_gateway.md#mapping-from-http-request-headers-to-grpc-client-metadata).
A matcher that keeps the old default alongside your own rules should fall back
to `runtime.LegacyHeaderMatcher` rather than `runtime.DefaultHeaderMatcher`.

## See also

- [Customizing your gateway](../mapping/customizing_your_gateway.md), which shows
  where the marshaler options go.
- [Custom marshalers](../mapping/custom_marshalers.md).
