package utilities

import "strings"

// KVVar is a command line flag.Value implementation accepting 0 or
// more instances on the command line with a value like
// "param1=some-str-val1,param2=some-str-val2". Multiple occurrences
// are concatenated and the last instance of a duplicate key wins.
type KVVar map[string]string

// NewKVVar return an empty KVar ready for Set() & String() calls
func NewKVVar() KVVar {
	return make(map[string]string)
}

// String returns the serialized value, e.g. "k=v,k2=v2" with keys in
// unspecified order across calls.
func (k KVVar) String() string {
	var b []byte
	i := 1
	for key, val := range k {
		b = append(b, key...)
		b = append(b, byte('='))
		b = append(b, val...)
		if i < len(k) {
			b = append(b, byte(','))
		}
		i++
	}
	return string(b)
}

// Set adds the given key value pairs to the KVVar in the order
// given. Syntax of given string is "k=v,k2=v2,k=v3,k3" where last
// duplicated key wins (k=v3) and keys without a value are set to ""
// (k3="").  Set may be called more than once, e.g. when a KVVar
// command line flag is given more than once, in which case they are
// added as described above. Leading and trailing whitespace around keys
// and values is removed.
func (k KVVar) Set(s string) error {
	kvs := strings.Split(s, ",")
	for _, kv := range kvs {
		p := strings.SplitN(kv, "=", 2)
		if len(p) == 1 {
			k[strings.TrimSpace(p[0])] = ""
			continue
		}
		k[strings.TrimSpace(p[0])] = strings.TrimSpace(p[1])
	}
	return nil
}

// Get makes this flag.Getter compatible
func (k KVVar) Get() interface{} {
	return k
}
