package runtime_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/gengo/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func TestAnnotateContext(t *testing.T) {
	ctx := context.Background()

	request, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(%q, %q, nil) failed with %v; want success", "GET", "http://localhost", err)
	}

	request.RemoteAddr = "127.0.0.1"
	request.Header.Add("Some-Irrelevant-Header", "some value")
	annotated := runtime.AnnotateContext(ctx, request)
	ctx = metadata.NewContext(ctx, metadata.Pairs("RemoteAddr", "127.0.0.1"))
	if annotated != ctx {
		t.Errorf("AnnotateContext(ctx, request) = %v; want %v", annotated, ctx)
	}

	request.Header.Add("Grpc-Metadata-FooBar", "Value1")
	request.Header.Add("Grpc-Metadata-Foo-BAZ", "Value2")
	request.Header.Add("Grpc-Metadata-foo-bAz", "Value3")
	request.Header.Add("Authorization", "Token 1234567890")
	annotated = runtime.AnnotateContext(ctx, request)
	md, ok := metadata.FromContext(annotated)
	if !ok || len(md) != 4 {
		t.Errorf("Expected 4 metadata items in context; got %v", md)
	}
	if got, want := md["foobar"], []string{"Value1"}; !reflect.DeepEqual(got, want) {
		t.Errorf(`md["foobar"] = %q; want %q`, got, want)
	}
	if got, want := md["foo-baz"], []string{"Value2", "Value3"}; !reflect.DeepEqual(got, want) {
		t.Errorf(`md["foo-baz"] = %q want %q`, got, want)
	}
	if got, want := md["authorization"], []string{"Token 1234567890"}; !reflect.DeepEqual(got, want) {
		t.Errorf(`md["authorization"] = %q want %q`, got, want)
	}
}
