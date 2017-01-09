package utilities

import (
	"reflect"
	"testing"
)

func TestKVVar(t *testing.T) {
	// simple
	kvv := NewKVVar()
	kvv.Set("x=y")
	kvcheck("x", "y", kvv, t)
	kvv2 := NewKVVar()
	kvv2.Set(kvv.String()) // check equiv string serialization
	kveq(kvv, kvv2, t)

	// more than one kv
	kvv = NewKVVar()
	kvv.Set("a=2,b=1")
	kvcheck("a", "2", kvv, t)
	kvcheck("b", "1", kvv, t)
	kvv2 = NewKVVar()
	kvv2.Set(kvv.String())
	kveq(kvv, kvv2, t)

	// empty input
	kvv = NewKVVar()
	kvv.Set("")
	kvv.Set("    ")
	if len(kvv) != 0 {
		t.Errorf("KVVar should remain empty after empty input string: %v", kvv)
	}

	// multiple Set() calls with mutiple kv's per calls
	kvv = NewKVVar()
	kvv.Set(" c = 1 , d = 2 ")
	kvv.Set("mountain = top of the world ,time=")
	kvcheck("c", "1", kvv, t)
	kvcheck("d", "2", kvv, t)
	kvcheck("mountain", "top of the world", kvv, t)
	kvcheck("time", "", kvv, t)
	kvv2 = NewKVVar()
	kvv2.Set(kvv.String())
	kveq(kvv, kvv2, t)

	// last dup key wins
	kvv = NewKVVar()
	kvv.Set("f=1,f=2")
	kvcheck("f", "2", kvv, t)
	kvv2 = NewKVVar()
	kvv2.Set(kvv.String())
	kveq(kvv, kvv2, t)

	kvv = NewKVVar()
	kvv.Set("g=1,h=2")
	kvv.Set("g=3")
	kvcheck("g", "3", kvv, t)
	kvcheck("h", "2", kvv, t)
	kvv2 = NewKVVar()
	kvv2.Set(kvv.String())
	kveq(kvv, kvv2, t)

	kvv = NewKVVar()
	kvv.Set("i=1")
	kvv.Set("i=3")
	kvv.Set("i=4,j=3")
	kvv.Set("k=4,foo=3")
	kvcheck("i", "4", kvv, t)
	kvcheck("j", "3", kvv, t)
	kvcheck("k", "4", kvv, t)
	kvcheck("foo", "3", kvv, t)
	kvv2 = NewKVVar()
	kvv2.Set(kvv.String())
	kveq(kvv, kvv2, t)

	if k2, ok := (kvv.Get()).(KVVar); !ok || !reflect.DeepEqual(k2, kvv) {
		t.Errorf("KVVar.Get is broken")
	}
}

func kvcheck(key, value string, kvv KVVar, t *testing.T) {
	if _, ok := kvv[key]; value != "" && !ok {
		t.Errorf("KVVar key '%v' does not exist", key)
	}
	if kvv[key] != value {
		t.Errorf("got '%v' for KVVar key '%v', expected '%v'", kvv[key], key, value)
	}
}

func kveq(kvv1, kvv2 KVVar, t *testing.T) {
	if !reflect.DeepEqual(kvv1, kvv2) {
		t.Errorf("KVars not equal: '%v' vs '%v'", kvv1, kvv2)
	}
}
