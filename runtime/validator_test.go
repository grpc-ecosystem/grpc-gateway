package runtime_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/grpc-ecosystem/grpc-gateway/v2/runtime/internal/examplepb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/protobuf/proto"
)

// testValidator rejects messages with empty Id.
func testValidator(msg proto.Message) error {
	if m, ok := msg.(*pb.SimpleMessage); ok && m.Id == "" {
		return errors.New("id is required")
	}
	return nil
}

func TestHandleWithMessage_StoresFactory(t *testing.T) {
	mux := runtime.NewServeMux()
	pat, err := runtime.NewPattern(1, []int{int(utilities.OpLitPush), 0}, []string{"test"}, "")
	if err != nil {
		t.Fatalf("NewPattern failed: %v", err)
	}

	factory := func() proto.Message { return &pb.SimpleMessage{} }
	mux.HandleWithMessage("POST", pat, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.WriteHeader(http.StatusOK)
	}, factory)

	key := "POST /test"
	f, ok := mux.LookupMessageFactory(key)
	if !ok {
		t.Fatalf("LookupMessageFactory(%q) not found", key)
	}
	msg := f()
	if _, ok := msg.(*pb.SimpleMessage); !ok {
		t.Fatalf("expected *SimpleMessage, got %T", msg)
	}
}

func TestHandleWithMessage_NilFactory(t *testing.T) {
	mux := runtime.NewServeMux()
	pat, err := runtime.NewPattern(1, []int{int(utilities.OpLitPush), 0}, []string{"test"}, "")
	if err != nil {
		t.Fatalf("NewPattern failed: %v", err)
	}

	mux.HandleWithMessage("GET", pat, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.WriteHeader(http.StatusOK)
	}, nil)

	key := "GET /test"
	_, ok := mux.LookupMessageFactory(key)
	if ok {
		t.Fatalf("LookupMessageFactory(%q) should not be found for nil factory", key)
	}
}

func TestWithValidator_ValidRequest(t *testing.T) {
	mux := runtime.NewServeMux(
		runtime.WithValidator(testValidator),
	)
	pat, err := runtime.NewPattern(1, []int{int(utilities.OpLitPush), 0}, []string{"test"}, "")
	if err != nil {
		t.Fatalf("NewPattern failed: %v", err)
	}

	handlerCalled := false
	mux.HandleWithMessage("POST", pat, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}, func() proto.Message { return &pb.SimpleMessage{} })

	body := bytes.NewReader([]byte(`{"id":"123"}`))
	r := httptest.NewRequest("POST", "/test", body)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("w.Code = %d; want %d", w.Code, http.StatusOK)
	}
	if !handlerCalled {
		t.Error("handler was not called")
	}
}

func TestWithValidator_InvalidRequest(t *testing.T) {
	mux := runtime.NewServeMux(
		runtime.WithValidator(testValidator),
	)
	pat, err := runtime.NewPattern(1, []int{int(utilities.OpLitPush), 0}, []string{"test"}, "")
	if err != nil {
		t.Fatalf("NewPattern failed: %v", err)
	}

	handlerCalled := false
	mux.HandleWithMessage("POST", pat, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}, func() proto.Message { return &pb.SimpleMessage{} })

	body := bytes.NewReader([]byte(`{"id":""}`))
	r := httptest.NewRequest("POST", "/test", body)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("w.Code = %d; want %d", w.Code, http.StatusBadRequest)
	}
	if handlerCalled {
		t.Error("handler should not be called for invalid request")
	}
}

func TestWithValidator_EmptyBody(t *testing.T) {
	mux := runtime.NewServeMux(
		runtime.WithValidator(testValidator),
	)
	pat, err := runtime.NewPattern(1, []int{int(utilities.OpLitPush), 0}, []string{"test"}, "")
	if err != nil {
		t.Fatalf("NewPattern failed: %v", err)
	}

	handlerCalled := false
	mux.HandleWithMessage("POST", pat, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}, func() proto.Message { return &pb.SimpleMessage{} })

	r := httptest.NewRequest("POST", "/test", nil)
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, r)

	// Empty body should skip validation and call handler
	if w.Code != http.StatusOK {
		t.Errorf("w.Code = %d; want %d", w.Code, http.StatusOK)
	}
	if !handlerCalled {
		t.Error("handler should be called for empty body")
	}
}

func TestWithValidator_NoFactoryRegistered(t *testing.T) {
	mux := runtime.NewServeMux(
		runtime.WithValidator(testValidator),
	)
	pat, err := runtime.NewPattern(1, []int{int(utilities.OpLitPush), 0}, []string{"test"}, "")
	if err != nil {
		t.Fatalf("NewPattern failed: %v", err)
	}

	// Use regular Handle (not HandleWithMessage) — no factory registered
	handlerCalled := false
	mux.Handle("GET", pat, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	r := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, r)

	// No factory → validation skipped → handler called
	if w.Code != http.StatusOK {
		t.Errorf("w.Code = %d; want %d", w.Code, http.StatusOK)
	}
	if !handlerCalled {
		t.Error("handler should be called when no factory is registered")
	}
}

func TestWithValidator_BodyReplayed(t *testing.T) {
	mux := runtime.NewServeMux(
		runtime.WithValidator(testValidator),
	)
	pat, err := runtime.NewPattern(1, []int{int(utilities.OpLitPush), 0}, []string{"test"}, "")
	if err != nil {
		t.Fatalf("NewPattern failed: %v", err)
	}

	var receivedBody string
	mux.HandleWithMessage("POST", pat, func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		receivedBody = buf.String()
		w.WriteHeader(http.StatusOK)
	}, func() proto.Message { return &pb.SimpleMessage{} })

	bodyStr := `{"id":"123"}`
	r := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte(bodyStr)))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, r)

	if receivedBody != bodyStr {
		t.Errorf("handler received body = %q; want %q", receivedBody, bodyStr)
	}
}
