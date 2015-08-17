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

	request, _ := http.NewRequest("GET", "http://localhost", nil)
	request.Header = http.Header{}
	annotated := runtime.AnnotateContext(ctx, request)
	if annotated != ctx {
		t.Errorf("AnnotateContext(ctx, request) = %v; want %v", annotated, ctx)
	}
	request.Header.Add("Grpc-Metadata-FooBar", "Value1")
	request.Header.Add("Grpc-Metadata-Foo-BAZ", "Value2")
	annotated = runtime.AnnotateContext(ctx, request)
	md, ok := metadata.FromContext(annotated)
	if !ok || len(md) != 2 {
		t.Errorf("Expected 2 metadata items in context; got %v", md)
	}
	if got, want := md["Foobar"], []string{"Value1"}; !reflect.DeepEqual(got, want) {
		t.Errorf("md[\"Foobar\"] = %v; want %v", got, want)
	}
	if got, want := md["Foo-Baz"], []string{"Value2"}; !reflect.DeepEqual(got, want) {
		t.Errorf("md[\"Foo-Baz\"] = %v want %v", got, want)
	}
}
