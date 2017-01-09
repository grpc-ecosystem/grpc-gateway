package main

import (
	"flag"
	"reflect"
	"testing"
)

func TestParseReqParam(t *testing.T) {

	f := flag.CommandLine

	// this one must be first - with no leading clearFlags call it
	// verifies our expectation of default values as we reset by
	// clearFlags
	err := parseReqParam("", f)
	if err != nil {
		t.Errorf("Test 0: unexpected parse error '%v'", err)
	}
	checkFlags(false, "stdin", "", map[string]string{}, t, 0)

	clearFlags()
	err = parseReqParam("allow_delete_body,file=./foo.pb,import_prefix=/bar/baz,Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api", f)
	if err != nil {
		t.Errorf("Test 1: unexpected parse error '%v'", err)
	}
	checkFlags(true, "./foo.pb", "/bar/baz",
		map[string]string{"google/api/annotations.proto": "github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api"}, t, 1)

	clearFlags()
	err = parseReqParam("allow_delete_body=true,file=./foo.pb,import_prefix=/bar/baz,Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api", f)
	if err != nil {
		t.Errorf("Test 2: unexpected parse error '%v'", err)
	}
	checkFlags(true, "./foo.pb", "/bar/baz",
		map[string]string{"google/api/annotations.proto": "github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api"}, t, 2)

	clearFlags()
	err = parseReqParam("allow_delete_body=false,Ma/b/c.proto=github.com/x/y/z,Mf/g/h.proto=github.com/1/2/3/", f)
	if err != nil {
		t.Errorf("Test 3: unexpected parse error '%v'", err)
	}
	checkFlags(false, "stdin", "", map[string]string{"a/b/c.proto": "github.com/x/y/z", "f/g/h.proto": "github.com/1/2/3/"}, t, 3)

	clearFlags()
	err = parseReqParam("", f)
	if err != nil {
		t.Errorf("Test 4: unexpected parse error '%v'", err)
	}
	checkFlags(false, "stdin", "", map[string]string{}, t, 4)

	clearFlags()
	err = parseReqParam("unknown_param=17", f)
	if err == nil {
		t.Error("Test 5: expected parse error not returned")
	}
	checkFlags(false, "stdin", "", map[string]string{}, t, 5)

	clearFlags()
	err = parseReqParam("Mfoo", f)
	if err == nil {
		t.Error("Test 6: expected parse error not returned")
	}
	checkFlags(false, "stdin", "", map[string]string{}, t, 6)

	clearFlags()
	err = parseReqParam("allow_delete_body,file,import_prefix", f)
	if err != nil {
		t.Errorf("Test 7: unexpected parse error '%v'", err)
	}
	checkFlags(true, "", "", map[string]string{}, t, 7)

}

func checkFlags(allowDeleteV bool, fileV, importPathV string, pkgMapV map[string]string, t *testing.T, tid int) {
	if *importPrefix != importPathV {
		t.Errorf("Test %v: import_prefix misparsed, expected '%v', got '%v'", tid, importPathV, *importPrefix)
	}
	if *file != fileV {
		t.Errorf("Test %v: file misparsed, expected '%v', got '%v'", tid, fileV, *file)
	}
	if *allowDeleteBody != allowDeleteV {
		t.Errorf("Test %v: allow_delete_body misparsed, expected '%v', got '%v'", tid, allowDeleteV, *allowDeleteBody)
	}
	if !reflect.DeepEqual(map[string]string(pkgMap), pkgMapV) {
		t.Errorf("Test %v: pkg_map misparsed, expected '%v', got '%v'", tid, pkgMapV, (map[string]string)(pkgMap))
	}
}

func clearFlags() {
	*importPrefix = ""
	*file = "stdin"
	*allowDeleteBody = false
	for k := range pkgMap {
		delete(pkgMap, k)
	}
}
