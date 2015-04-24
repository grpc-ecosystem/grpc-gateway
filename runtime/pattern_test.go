package runtime

import (
	"reflect"
	"strings"
	"testing"
)

const (
	validVersion = 1
	anything     = 0
)

func TestNewPattern(t *testing.T) {
	for _, spec := range []struct {
		ops  []int
		pool []string

		stackSizeWant int
	}{
		{},
		{
			ops:           []int{int(opNop), anything},
			stackSizeWant: 0,
		},
		{
			ops:           []int{int(opPush), anything},
			stackSizeWant: 1,
		},
		{
			ops:           []int{int(opLitPush), 0},
			pool:          []string{"abc"},
			stackSizeWant: 1,
		},
		{
			ops:           []int{int(opPushM), anything},
			stackSizeWant: 1,
		},
		{
			ops: []int{
				int(opPush), anything,
				int(opConcatN), 1,
			},
			stackSizeWant: 1,
		},
		{
			ops: []int{
				int(opPush), anything,
				int(opConcatN), 1,
				int(opCapture), 0,
			},
			pool:          []string{"abc"},
			stackSizeWant: 1,
		},
		{
			ops: []int{
				int(opPush), anything,
				int(opLitPush), 0,
				int(opLitPush), 1,
				int(opPushM), anything,
				int(opConcatN), 2,
				int(opCapture), 2,
			},
			pool:          []string{"lit1", "lit2", "var1"},
			stackSizeWant: 4,
		},
	} {
		pat, err := NewPattern(validVersion, spec.ops, spec.pool)
		if err != nil {
			t.Errorf("NewPattern(%d, %v, %q) failed with %v; want success", validVersion, spec.ops, spec.pool, err)
			continue
		}
		if got, want := pat.stacksize, spec.stackSizeWant; got != want {
			t.Errorf("pat.stacksize = %d; want %d", got, want)
		}
	}
}

func TestNewPatternWithWrongOp(t *testing.T) {
	for _, spec := range []struct {
		ops  []int
		pool []string
	}{
		{
			// op code out of bound
			ops: []int{-1, anything},
		},
		{
			// op code out of bound
			ops: []int{int(opEnd), 0},
		},
		{
			// odd number of items
			ops: []int{int(opPush)},
		},
		{
			// negative index
			ops:  []int{int(opLitPush), -1},
			pool: []string{"abc"},
		},
		{
			// index out of bound
			ops:  []int{int(opLitPush), 1},
			pool: []string{"abc"},
		},
		{
			// negative # of segments
			ops:  []int{int(opConcatN), -1},
			pool: []string{"abc"},
		},
		{
			// negative index
			ops:  []int{int(opCapture), -1},
			pool: []string{"abc"},
		},
		{
			// index out of bound
			ops:  []int{int(opCapture), 1},
			pool: []string{"abc"},
		},
	} {
		_, err := NewPattern(validVersion, spec.ops, spec.pool)
		if err == nil {
			t.Errorf("NewPattern(%d, %v, %q) succeeded; want failure with %v", validVersion, spec.ops, spec.pool, ErrInvalidPattern)
			continue
		}
		if err != ErrInvalidPattern {
			t.Errorf("NewPattern(%d, %v, %q) failed with %v; want failure with %v", validVersion, spec.ops, spec.pool, err, ErrInvalidPattern)
			continue
		}
	}
}

func TestNewPatternWithStackUnderflow(t *testing.T) {
	for _, spec := range []struct {
		ops  []int
		pool []string
	}{
		{
			ops: []int{int(opConcatN), 1},
		},
		{
			ops:  []int{int(opCapture), 0},
			pool: []string{"abc"},
		},
	} {
		_, err := NewPattern(validVersion, spec.ops, spec.pool)
		if err == nil {
			t.Errorf("NewPattern(%d, %v, %q) succeeded; want failure with %v", validVersion, spec.ops, spec.pool, ErrInvalidPattern)
			continue
		}
		if err != ErrInvalidPattern {
			t.Errorf("NewPattern(%d, %v, %q) failed with %v; want failure with %v", validVersion, spec.ops, spec.pool, err, ErrInvalidPattern)
			continue
		}
	}
}

func segments(path string) []string {
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}

func TestMatch(t *testing.T) {
	for _, spec := range []struct {
		ops  []int
		pool []string

		match    []string
		notMatch []string
	}{
		{
			match:    []string{""},
			notMatch: []string{"example"},
		},
		{
			ops:      []int{int(opNop), anything},
			match:    []string{""},
			notMatch: []string{"example", "path/to/example"},
		},
		{
			ops:      []int{int(opPush), anything},
			match:    []string{"abc", "def"},
			notMatch: []string{"", "abc/def"},
		},
		{
			ops:      []int{int(opLitPush), 0},
			pool:     []string{"v1"},
			match:    []string{"v1"},
			notMatch: []string{"", "v2"},
		},
		{
			ops:   []int{int(opPushM), anything},
			match: []string{"", "abc", "abc/def", "abc/def/ghi"},
		},
		{
			ops: []int{
				int(opLitPush), 0,
				int(opLitPush), 1,
				int(opPush), anything,
				int(opConcatN), 1,
				int(opCapture), 2,
			},
			pool:  []string{"v1", "bucket", "name"},
			match: []string{"v1/bucket/my-bucket", "v1/bucket/our-bucket"},
			notMatch: []string{
				"",
				"v1",
				"v1/bucket",
				"v2/bucket/my-bucket",
				"v1/pubsub/my-topic",
			},
		},
		{
			ops: []int{
				int(opLitPush), 0,
				int(opLitPush), 1,
				int(opPushM), anything,
				int(opConcatN), 2,
				int(opCapture), 2,
			},
			pool: []string{"v1", "o", "name"},
			match: []string{
				"v1/o",
				"v1/o/my-bucket",
				"v1/o/our-bucket",
				"v1/o/my-bucket/dir",
				"v1/o/my-bucket/dir/dir2",
				"v1/o/my-bucket/dir/dir2/obj",
			},
			notMatch: []string{
				"",
				"v1",
				"v2/o/my-bucket",
				"v1/b/my-bucket",
			},
		},
		{
			ops: []int{
				int(opLitPush), 0,
				int(opLitPush), 1,
				int(opPush), anything,
				int(opConcatN), 2,
				int(opCapture), 2,
				int(opLitPush), 3,
				int(opPush), anything,
				int(opConcatN), 1,
				int(opCapture), 4,
			},
			pool: []string{"v2", "b", "name", "o", "oname"},
			match: []string{
				"v2/b/my-bucket/o/obj",
				"v2/b/our-bucket/o/obj",
				"v2/b/my-bucket/o/dir",
			},
			notMatch: []string{
				"",
				"v2",
				"v2/b",
				"v2/b/my-bucket",
				"v2/b/my-bucket/o",
			},
		},
	} {
		pat, err := NewPattern(validVersion, spec.ops, spec.pool)
		if err != nil {
			t.Errorf("NewPattern(%d, %v, %q) failed with %v; want success", validVersion, spec.ops, spec.pool, err)
			continue
		}

		for _, path := range spec.match {
			_, err = pat.Match(segments(path))
			if err != nil {
				t.Errorf("pat.Match(%q) failed with %v; want success; pattern = (%v, %q)", path, err, spec.ops, spec.pool)
			}
		}

		for _, path := range spec.notMatch {
			_, err = pat.Match(segments(path))
			if err == nil {
				t.Errorf("pat.Match(%q) succeeded; want failure with %v; pattern = (%v, %q)", path, ErrNotMatch, spec.ops, spec.pool)
				continue
			}
			if err != ErrNotMatch {
				t.Errorf("pat.Match(%q) failed with %v; want failure with %v; pattern = (%v, %q)", spec.notMatch, err, ErrNotMatch, spec.ops, spec.pool)
			}
		}
	}
}

func TestMatchWithBinding(t *testing.T) {
	for _, spec := range []struct {
		ops  []int
		pool []string
		path string

		want map[string]string
	}{
		{
			want: make(map[string]string),
		},
		{
			ops:  []int{int(opNop), anything},
			want: make(map[string]string),
		},
		{
			ops:  []int{int(opPush), anything},
			path: "abc",
			want: make(map[string]string),
		},
		{
			ops:  []int{int(opLitPush), 0},
			pool: []string{"endpoint"},
			path: "endpoint",
			want: make(map[string]string),
		},
		{
			ops:  []int{int(opPushM), anything},
			path: "abc/def/ghi",
			want: make(map[string]string),
		},
		{
			ops: []int{
				int(opLitPush), 0,
				int(opLitPush), 1,
				int(opPush), anything,
				int(opConcatN), 1,
				int(opCapture), 2,
			},
			pool: []string{"v1", "bucket", "name"},
			path: "v1/bucket/my-bucket",
			want: map[string]string{
				"name": "my-bucket",
			},
		},
		{
			ops: []int{
				int(opLitPush), 0,
				int(opLitPush), 1,
				int(opPushM), anything,
				int(opConcatN), 2,
				int(opCapture), 2,
			},
			pool: []string{"v1", "o", "name"},
			path: "v1/o/my-bucket/dir/dir2/obj",
			want: map[string]string{
				"name": "o/my-bucket/dir/dir2/obj",
			},
		},
		{
			ops: []int{
				int(opLitPush), 0,
				int(opLitPush), 1,
				int(opPush), anything,
				int(opConcatN), 2,
				int(opCapture), 2,
				int(opLitPush), 3,
				int(opPush), anything,
				int(opConcatN), 1,
				int(opCapture), 4,
			},
			pool: []string{"v2", "b", "name", "o", "oname"},
			path: "v2/b/my-bucket/o/obj",
			want: map[string]string{
				"name":  "b/my-bucket",
				"oname": "obj",
			},
		},
	} {
		pat, err := NewPattern(validVersion, spec.ops, spec.pool)
		if err != nil {
			t.Errorf("NewPattern(%d, %v, %q) failed with %v; want success", validVersion, spec.ops, spec.pool, err)
			continue
		}

		got, err := pat.Match(segments(spec.path))
		if err != nil {
			t.Errorf("pat.Match(%q) failed with %v; want success; pattern = (%v, %q)", spec.path, err, spec.ops, spec.pool)
		}
		if !reflect.DeepEqual(got, spec.want) {
			t.Errorf("pat.Match(%q) = %q; want %q; pattern = (%v, %q)", spec.path, got, spec.want, spec.ops, spec.pool)
		}
	}
}
