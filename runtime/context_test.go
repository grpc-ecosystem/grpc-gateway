package runtime_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/gengo/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const (
	emptyForwardMetaCount = 2
)

func TestAnnotateContext_WorksWithEmpty(t *testing.T) {
	ctx := context.Background()

	request, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(%q, %q, nil) failed with %v; want success", "GET", "http://www.example.com", err)
	}
	// Make sure we set a remote.
	request.RemoteAddr = "192.168.0.1:12345"

	request.Header.Add("Some-Irrelevant-Header", "some value")
	annotated := runtime.AnnotateContext(ctx, request)
	md, ok := metadata.FromContext(annotated)
	if !ok || len(md) != emptyForwardMetaCount {
		t.Errorf("Expected 2 metadata items in context; got %v", md)
	}
}

func TestAnnotateContext_ForwardsGrpcMetadata(t *testing.T) {
	ctx := context.Background()
	request, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(%q, %q, nil) failed with %v; want success", "GET", "http://www.example.com", err)
	}
	request.RemoteAddr = "192.168.0.1:12345"

	request.Header.Add("Some-Irrelevant-Header", "some value")
	request.Header.Add("Grpc-Metadata-FooBar", "Value1")
	request.Header.Add("Grpc-Metadata-Foo-BAZ", "Value2")
	request.Header.Add("Grpc-Metadata-foo-bAz", "Value3")
	request.Header.Add("Authorization", "Token 1234567890")
	annotated := runtime.AnnotateContext(ctx, request)
	md, ok := metadata.FromContext(annotated)
	if !ok || len(md) != emptyForwardMetaCount+3 {
		t.Errorf("Expected 5 metadata items in context; got %v", md)
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

func TestAnnotateContext_XForwardedFor(t *testing.T) {
	ctx := context.Background()
	request, err := http.NewRequest("GET", "http://bar.foo.example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(%q, %q, nil) failed with %v; want success", "GET", "http://bar.foo.example.com", err)
	}
	request.Header.Add("X-Forwarded-For", "192.168.0.100") // client
	request.RemoteAddr = "8.8.8.8:12345"                   // proxy

	annotated := runtime.AnnotateContext(ctx, request)
	md, ok := metadata.FromContext(annotated)
	if !ok || len(md) != emptyForwardMetaCount {
		t.Errorf("Expected 2 metadata items in context; got %v", md)
	}
	if got, want := md["x-forwarded-host"], []string{"bar.foo.example.com"}; !reflect.DeepEqual(got, want) {
		t.Errorf("md[\"host\"] = %v; want %v", got, want)
	}
	// Note: it must be in order client, proxy1, proxy2
	if got, want := md["x-forwarded-for"], []string{"192.168.0.100, 8.8.8.8"}; !reflect.DeepEqual(got, want) {
		t.Errorf("md[\"x-forwarded-for\"] = %v want %v", got, want)
	}
}
