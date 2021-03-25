---
layout: default
title: Binary file uploads
nav_order: 2
parent: Mapping
---

# Binary file uploads

If you need to do a binary file upload, e.g. via;

```sh
curl -X POST -F "attachment=@/tmp/somefile.txt" http://localhost:9090/v1/files
```

then your request will contain the binary data directly and there is no way to model this using gRPC.

What you can do instead is to add a custom route directly on the `mux` instance.

## Custom route on a mux instance

Here we'll setup a handler (`handleBinaryFileUpload`) for `POST` requests: 

```go
// Create a mux instance
mux := runtime.NewServeMux()

// Attachment upload from http/s handled manually
mux.HandlePath("POST", "/v1/files", handleBinaryFileUpload)
```

And then in your handler you can do something like:

```go
func handleBinaryFileUpload(w http.ResponseWriter, rq *http.Request, params map[string]string) {
  
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	f, header, err := rq.FormFile("attachment")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// header contains filename etc if you need this

	var buffer bytes.Buffer
	_, err = io.Copy(buffer, f)
	if err != nil {
		panic(err)
	}

	//
	// Now do something with the bytes in the `buffer`
	//

	w.WriteHeader(http.StatusOK)
}
```
