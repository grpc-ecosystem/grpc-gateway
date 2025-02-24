package runtime

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime/internal/examplepb"
	"google.golang.org/protobuf/proto"
)

func TestUrlEncodedDecoder_Decode(t *testing.T) {
	tests := []struct {
		name    string
		values  url.Values
		want    proto.Message
		wantErr bool
	}{
		{
			name: "simple form fields",
			values: url.Values{
				"single_nested.name":   {"test"},
				"single_nested.amount": {"42"},
			},
			want: &examplepb.ABitOfEverything{
				SingleNested: &examplepb.ABitOfEverything_Nested{
					Name:   "test",
					Amount: 42,
				},
			},
			wantErr: false,
		},
		{
			name: "fields with special characters",
			values: url.Values{
				"single_nested.name":   {"Hello World!"},
				"single_nested.amount": {"123"},
			},
			want: &examplepb.ABitOfEverything{
				SingleNested: &examplepb.ABitOfEverything_Nested{
					Name:   "Hello World!",
					Amount: 123,
				},
			},
			wantErr: false,
		},
		{
			name:    "empty input",
			values:  url.Values{},
			want:    &examplepb.ABitOfEverything{},
			wantErr: false,
		},
		{
			name: "repeated field",
			values: url.Values{
				"repeated_string": {"one", "two", "three"},
			},
			want: &examplepb.ABitOfEverything{
				RepeatedStringValue: []string{"one", "two", "three"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "http://example.com", strings.NewReader(tt.values.Encode()))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			decoder := NewUrlEncodedDecoder(req.Body)
			msg := &examplepb.ABitOfEverything{}

			err = decoder.Decode(msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("UrlEncodedDecoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !proto.Equal(msg, tt.want) {
				t.Errorf("UrlEncodedDecoder.Decode() = %v, want %v", msg, tt.want)
			}
		})
	}
}

func TestUrlEncodedDecoder_DecodeNonProto(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://example.com", strings.NewReader(""))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	decoder := NewUrlEncodedDecoder(req.Body)
	var nonProto struct{}

	err = decoder.Decode(&nonProto)
	if err == nil {
		t.Error("UrlEncodedDecoder.Decode() expected error for non-proto message")
	}
}

func TestUrlEncodeMarshal_ContentType(t *testing.T) {
	m := &UrlEncodeMarshal{}
	if got := m.ContentType(nil); got != "application/x-www-form-urlencoded" {
		t.Errorf("UrlEncodeMarshal.ContentType() = %v, want application/x-www-form-urlencoded", got)
	}
}

func TestUrlEncodeMarshal_Marshal(t *testing.T) {
	msg := &examplepb.ABitOfEverything{
		SingleNested: &examplepb.ABitOfEverything_Nested{
			Name:   "test",
			Amount: 42,
		},
	}

	marshaler := &UrlEncodeMarshal{
		Marshaler: &JSONPb{},
	}

	got, err := marshaler.Marshal(msg)
	if err != nil {
		t.Fatalf("UrlEncodeMarshal.Marshal() error = %v", err)
	}

	want := []byte(`{"single_nested":{"name":"test","amount":42}}`)
	if !bytes.Equal(got, want) {
		t.Errorf("UrlEncodeMarshal.Marshal() = %s, want %s", got, want)
	}
}

func TestUrlEncodeMarshal_NewDecoder(t *testing.T) {
	m := &UrlEncodeMarshal{}
	req, err := http.NewRequest(http.MethodPost, "http://example.com", strings.NewReader(""))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	decoder := m.NewDecoder(req.Body)

	if _, ok := decoder.(*UrlEncodedDecoder); !ok {
		t.Error("UrlEncodeMarshal.NewDecoder() did not return *UrlEncodedDecoder")
	}
}
