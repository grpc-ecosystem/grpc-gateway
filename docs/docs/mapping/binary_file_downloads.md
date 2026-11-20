---
layout: default
title: Binary file downloads
nav_order: 4
parent: Mapping
permalink: /docs/mapping/binary_file_downloads/
---

# Binary file downloads

If you need to return a file download from your gRPC-Gateway endpoint (e.g. a CSV
export, PDF, or any binary file), the recommended approach is to use
`google.api.HttpBody` as your response type. This ensures the route is registered
natively in gRPC-Gateway and works reliably in all environments — including those
behind a reverse proxy or API gateway.

> **Why not use a custom `mux.HandlePath` route for downloads?**
>
> For file **uploads**, a custom route on the mux is necessary because multipart
> form data cannot be modelled in gRPC (see [Binary file uploads](binary_file_uploads.md)).
>
> For file **downloads**, using a custom route (e.g. a Gin catch-all handler) can
> work in local development but will silently fail in environments behind a reverse
> proxy or API gateway — because the proxy only knows about routes registered
> through gRPC-Gateway. Using `google.api.HttpBody` registers the route properly
> and avoids this inconsistency.

## Define your proto

Add `google/api/httpbody.proto` to your imports and use `google.api.HttpBody` as
the return type:

```protobuf
syntax = "proto3";

import "google/api/annotations.proto";
import "google/api/httpbody.proto";

service ReportService {
  rpc ExportReport(ExportReportRequest) returns (google.api.HttpBody) {
    option (google.api.http) = {
      get: "/v1/reports/export"
    };
  }
}

message ExportReportRequest {
  int64 date_from = 1;
  int64 date_to   = 2;
}
```

## Implement the handler in Go

In your gRPC server implementation, build your file content, set the
`Content-Disposition` header via `grpc.SetHeader`, and return an `HttpBody`:

```go
import (
    "bytes"
    "context"
    "fmt"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/status"
    httpbody "google.golang.org/genproto/googleapis/api/httpbody"

    pb "your/generated/proto/package"
)

func (s *Server) ExportReport(
    ctx context.Context,
    req *pb.ExportReportRequest,
) (*httpbody.HttpBody, error) {

    // Build your file content — CSV, PDF, etc.
    var buf bytes.Buffer
    buf.WriteString("col1,col2,col3\n")
    buf.WriteString("val1,val2,val3\n")

    // Set Content-Disposition so the browser triggers a file save dialog.
    filename := fmt.Sprintf("report-%s.csv", time.Now().Format("20060102"))
    grpc.SetHeader(ctx, metadata.Pairs(
        "Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename),
    ))

    return &httpbody.HttpBody{
        ContentType: "text/csv; charset=utf-8",
        Data:        buf.Bytes(),
    }, nil
}
```

## How the browser receives it

gRPC-Gateway forwards the `Content-Disposition` metadata header to the HTTP
response. The browser sees a standard file download response:

```
HTTP/1.1 200 OK
Content-Type: text/csv; charset=utf-8
Content-Disposition: attachment; filename="report-20240315.csv"
```

You can trigger the download from a frontend using the
[File System Access API](https://developer.mozilla.org/en-US/docs/Web/API/File_System_Access_API)
for a native Save As dialog (with fallback for unsupported browsers):

```js
async function downloadReport(dateFrom, dateTo) {
  const url = `/v1/reports/export?date_from=${dateFrom}&date_to=${dateTo}`;
  const response = await fetch(url);

  if (!response.ok) {
    throw new Error(`Export failed: ${response.status}`);
  }

  const blob = await response.blob();

  // Use File System Access API if available (Chrome, Edge)
  if (window.showSaveFilePicker) {
    try {
      const handle = await window.showSaveFilePicker({
        suggestedName: "report.csv",
        types: [{ description: "CSV file", accept: { "text/csv": [".csv"] } }],
      });
      const writable = await handle.createWritable();
      await writable.write(blob);
      await writable.close();
      return;
    } catch (err) {
      if (err.name === "AbortError") return; // user cancelled
    }
  }

  // Fallback: trigger automatic download
  const blobUrl = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = blobUrl;
  a.download = "report.csv";
  a.click();
  URL.revokeObjectURL(blobUrl);
}
```

## Important: reverse proxy environments

If your service runs behind a reverse proxy (Nginx, Envoy, etc.), the proxy
routes requests based on the paths registered with gRPC-Gateway. A custom
handler added outside of gRPC-Gateway (e.g. via a separate Gin router) will
**not** be known to the proxy and requests will fail with 404 or 502 in staging
and production environments, even if it works on localhost.

Using `google.api.HttpBody` as shown above registers the download route through
gRPC-Gateway itself, so it is visible to the proxy and behaves consistently
across all environments.

## Filenames with non-ASCII characters

If your filename contains non-ASCII characters (e.g. accented characters or
CJK characters), the `Content-Disposition` header will be rejected by the Go
HTTP layer with:

```
header key "content-disposition" contains value with non-printable ASCII characters
```

Use RFC 5987 encoding for such filenames:

```go
import "net/url"

encodedName := url.PathEscape(filename)
grpc.SetHeader(ctx, metadata.Pairs(
    "Content-Disposition",
    fmt.Sprintf("attachment; filename*=UTF-8''%s", encodedName),
))
```