package runtime

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime/internal/examplepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

func TestFieldMaskDynamicEmptyObjects(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		msg   proto.Message
		paths []string
	}{
		{
			name: "struct", input: `{"structField":{}}`,
			msg: &examplepb.NonStandardMessage{}, paths: []string{"struct_field"},
		},
		{
			name: "value", input: `{"valueField":{}}`,
			msg: &examplepb.NonStandardMessage{}, paths: []string{"value_field"},
		},
		{
			name: "nested empty struct", input: `{"structField":{"child":{}}}`,
			msg: &examplepb.NonStandardMessage{}, paths: []string{"struct_field.child"},
		},
		{
			name: "deep empty value", input: `{"valueField":{"parent":{"child":{}}}}`,
			msg: &examplepb.NonStandardMessage{}, paths: []string{"value_field.parent.child"},
		},
		{
			name: "siblings", input: `{"structField":{"empty":{},"populated":{"leaf":1}},"valueField":{}}`,
			msg: &examplepb.NonStandardMessage{}, paths: []string{"struct_field.empty", "struct_field.populated.leaf", "value_field"},
		},
		{
			name: "parent prefix", input: `{"body":{"structField":{},"valueField":{"child":{}}}}`,
			msg: &examplepb.NonStandardUpdateRequest{}, paths: []string{"body.struct_field", "body.value_field.child"},
		},
		{
			name: "custom JSON names", input: `{"StructField":{},"ValueField":{"child":{}}}`,
			msg: &examplepb.NonStandardMessageWithJSONNames{}, paths: []string{"struct_field", "value_field.child"},
		},
		{
			name: "custom JSON names with parent", input: `{"body":{"StructField":{"child":{}},"ValueField":{}}}`,
			msg: &examplepb.NonStandardWithJSONNamesUpdateRequest{}, paths: []string{"body.struct_field.child", "body.value_field"},
		},
		{
			name: "unrelated field", input: `{"id":"unchanged"}`,
			msg: &examplepb.NonStandardMessage{}, paths: []string{"id"},
		},
		{
			name: "ordinary empty message", input: `{"thing":{}}`,
			msg: &examplepb.NonStandardMessage{}, paths: []string{"thing"},
		},
		{
			name: "null", input: `{"structField":null,"valueField":null}`,
			msg: &examplepb.NonStandardMessage{}, paths: []string{"struct_field", "value_field"},
		},
		{
			name: "array and scalar", input: `{"structField":{"items":[]},"valueField":false}`,
			msg: &examplepb.NonStandardMessage{}, paths: []string{"struct_field.items", "value_field"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := (&JSONPb{}).Unmarshal([]byte(tc.input), tc.msg); err != nil {
				t.Fatalf("invalid protobuf JSON fixture: %v", err)
			}
			before := proto.Clone(tc.msg)
			mask, err := FieldMaskFromRequestBody(strings.NewReader(tc.input), tc.msg)
			if err != nil {
				t.Fatalf("FieldMaskFromRequestBody: %v", err)
			}
			if !reflect.DeepEqual(mask.GetPaths(), tc.paths) {
				t.Errorf("paths = %v, want %v", mask.GetPaths(), tc.paths)
			}
			if !proto.Equal(tc.msg, before) {
				t.Error("field mask extraction changed the request message")
			}
		})
	}
}

type fieldMaskEmptyObjectServer struct {
	examplepb.UnimplementedNonStandardServiceServer
	requests chan *examplepb.NonStandardUpdateRequest
}

func (s *fieldMaskEmptyObjectServer) Update(_ context.Context, req *examplepb.NonStandardUpdateRequest) (*examplepb.NonStandardMessage, error) {
	s.requests <- req
	return req.GetBody(), nil
}

func TestFieldMaskDynamicEmptyObjectsHTTPGRPC(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	requests := make(chan *examplepb.NonStandardUpdateRequest, 1)
	server := grpc.NewServer()
	examplepb.RegisterNonStandardServiceServer(server, &fieldMaskEmptyObjectServer{requests: requests})
	go func() { _ = server.Serve(lis) }()
	t.Cleanup(server.Stop)
	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	client := examplepb.NewNonStandardServiceClient(conn)
	marshaler := &JSONPb{}
	mux := NewServeMux()
	err = mux.HandlePath(http.MethodPatch, "/dynamic", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req := &examplepb.NonStandardUpdateRequest{Body: &examplepb.NonStandardMessage{}}
		if err := marshaler.Unmarshal(body, req.Body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.UpdateMask, err = FieldMaskFromRequestBody(bytes.NewReader(body), req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result, err := client.Update(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		if err := marshaler.NewEncoder(w).Encode(result); err != nil {
			t.Errorf("encode response: %v", err)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	httpServer := httptest.NewServer(mux)
	t.Cleanup(httpServer.Close)
	for _, tc := range []struct {
		name  string
		input string
		paths []string
	}{
		{name: "struct", input: `{"structField":{}}`, paths: []string{"struct_field"}},
		{name: "value", input: `{"valueField":{}}`, paths: []string{"value_field"}},
		{name: "nested", input: `{"structField":{"child":{}},"valueField":{"parent":{"child":{}}}}`, paths: []string{"struct_field.child", "value_field.parent.child"}},
		{name: "siblings", input: `{"structField":{"empty":{},"populated":{"leaf":1}},"valueField":{}}`, paths: []string{"struct_field.empty", "struct_field.populated.leaf", "value_field"}},
		{name: "unrelated field", input: `{"id":"unchanged"}`, paths: []string{"id"}},
		{name: "null and scalar", input: `{"structField":null,"valueField":false}`, paths: []string{"struct_field", "value_field"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, http.MethodPatch, httpServer.URL+"/dynamic", strings.NewReader(tc.input))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := httpServer.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = resp.Body.Close() }()
			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				t.Fatalf("HTTP status = %d: %s", resp.StatusCode, body)
			}
			var returned examplepb.NonStandardMessage
			if err := marshaler.NewDecoder(resp.Body).Decode(&returned); err != nil {
				t.Fatal(err)
			}
			select {
			case received := <-requests:
				if !reflect.DeepEqual(received.GetUpdateMask().GetPaths(), tc.paths) {
					t.Errorf("gRPC paths = %v, want %v", received.GetUpdateMask().GetPaths(), tc.paths)
				}
				var expected examplepb.NonStandardMessage
				if err := marshaler.Unmarshal([]byte(tc.input), &expected); err != nil {
					t.Fatal(err)
				}
				if !proto.Equal(received.GetBody(), &expected) || !proto.Equal(&returned, &expected) {
					t.Error("HTTP/gRPC round trip changed the dynamic field values")
				}
			case <-ctx.Done():
				t.Fatal("gRPC service did not receive the request")
			}
		})
	}
}
