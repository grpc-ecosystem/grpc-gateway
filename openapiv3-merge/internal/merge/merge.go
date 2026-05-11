// Package merge merges multiple OpenAPI 3.1 JSON documents into a single
// document.
//
// The merger is strict by design: it refuses to silently overwrite anything.
// Two inputs that declare the same path, the same component name, or the
// same tag name with conflicting content are rejected as errors. In
// practice, output from protoc-gen-openapiv3 merges cleanly across packages
// because component names are fully-qualified proto names (e.g.
// `example.v1.User`); conflicts surface only for genuine mistakes.
//
// Field-by-field rules:
//
//   - `openapi` must match across all inputs.
//   - `info`, `servers`, `externalDocs`, and any unknown top-level keys
//     (including OpenAPI extensions `x-*`) are taken from the first
//     input; the corresponding fields in later inputs are silently
//     discarded. This matches the practical case where
//     protoc-gen-openapiv3 derives a default `info.title` from each
//     file's name — without this lenience, the common case would never
//     merge.
//   - `paths` and `webhooks` are unioned in input order; key collisions
//     with non-identical values are an error.
//   - `components/*` sub-maps are unioned with sorted-key output;
//     same-named entries with non-identical values are an error.
//   - `tags` are deduplicated by `name`; same-named entries with
//     non-identical metadata are an error.
//   - `security` is first-wins: the first input to declare a non-empty
//     `security` array sets it; later inputs that declare a different
//     non-empty array are rejected.
package merge

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

// Input is a single OpenAPI 3.1 JSON document to merge.
type Input struct {
	// Name identifies the document in error messages, typically the path
	// the document was loaded from.
	Name string
	// Data is the raw JSON content of the document.
	Data []byte
}

// Merge combines inputs into a single OpenAPI 3.1 JSON document and returns
// its pretty-printed bytes (two-space indent, matching protoc-gen-openapiv3).
//
// At least one input is required. The first input's `info`, `servers`, and
// `externalDocs` are kept in the merged document.
func Merge(inputs []Input) ([]byte, error) {
	if len(inputs) == 0 {
		return nil, errors.New("at least one input is required")
	}
	docs := make([]*document, len(inputs))
	for i, in := range inputs {
		d, err := parse(in)
		if err != nil {
			return nil, err
		}
		docs[i] = d
	}
	merged, err := mergeAll(docs)
	if err != nil {
		return nil, err
	}
	// Drop empty container fields so encoding/json's omitempty omits them.
	if merged.Paths.len() == 0 {
		merged.Paths = nil
	}
	if merged.Webhooks.len() == 0 {
		merged.Webhooks = nil
	}
	if merged.Components.empty() {
		merged.Components = nil
	}
	if merged.extras.len() == 0 {
		merged.extras = nil
	}
	return json.MarshalIndent(merged, "", "  ")
}

// document is the in-memory and serialised shape of an OpenAPI 3.1
// document. Field order matches the spec, so encoding/json emits them in
// the canonical order; encoding/json also sorts map keys alphabetically,
// which is what protoc-gen-openapiv3 itself emits per file. paths and
// webhooks preserve input-file ordering via orderedObject's MarshalJSON.
type document struct {
	OpenAPI      string            `json:"openapi"`
	Info         json.RawMessage   `json:"info,omitempty"`
	Servers      json.RawMessage   `json:"servers,omitempty"`
	Paths        *orderedObject    `json:"paths,omitempty"`
	Webhooks     *orderedObject    `json:"webhooks,omitempty"`
	Components   *components       `json:"components,omitempty"`
	Security     []json.RawMessage `json:"security,omitempty"`
	Tags         []json.RawMessage `json:"tags,omitempty"`
	ExternalDocs json.RawMessage   `json:"externalDocs,omitempty"`

	// name identifies the input in error messages and is not serialised.
	name string
	// extras holds top-level keys outside the OpenAPI 3.1 defined set,
	// notably extensions (`x-*`). The json:"-" tag keeps them out of
	// the default marshaler; MarshalJSON splices them in after the
	// known fields, preserving first-occurrence order.
	extras *orderedObject `json:"-"`
}

// MarshalJSON marshals the struct fields via the standard encoder, then
// splices in extras before the closing brace. The alias type strips this
// MarshalJSON method so the inner Marshal call doesn't recurse.
func (d *document) MarshalJSON() ([]byte, error) {
	type alias document
	body, err := json.Marshal((*alias)(d))
	if err != nil {
		return nil, err
	}
	if d.extras == nil || d.extras.len() == 0 {
		return body, nil
	}
	// body ends with '}'. If the struct produced an empty object ("{}"),
	// the spliced extras must not be preceded by a comma.
	body = body[:len(body)-1]
	needComma := !bytes.Equal(body, []byte("{"))
	var buf bytes.Buffer
	buf.Write(body)
	for _, k := range d.extras.keys {
		if needComma {
			buf.WriteByte(',')
		}
		needComma = true
		keyJSON, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		buf.Write(keyJSON)
		buf.WriteByte(':')
		buf.Write(d.extras.vals[k])
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// components captures the components/* section. Field order matches
// OpenAPI 3.1.0 §4.8.7; encoding/json sorts the keys of each map.
type components struct {
	Schemas         map[string]json.RawMessage `json:"schemas,omitempty"`
	Responses       map[string]json.RawMessage `json:"responses,omitempty"`
	Parameters      map[string]json.RawMessage `json:"parameters,omitempty"`
	Examples        map[string]json.RawMessage `json:"examples,omitempty"`
	RequestBodies   map[string]json.RawMessage `json:"requestBodies,omitempty"`
	Headers         map[string]json.RawMessage `json:"headers,omitempty"`
	SecuritySchemes map[string]json.RawMessage `json:"securitySchemes,omitempty"`
	Links           map[string]json.RawMessage `json:"links,omitempty"`
	Callbacks       map[string]json.RawMessage `json:"callbacks,omitempty"`
	PathItems       map[string]json.RawMessage `json:"pathItems,omitempty"`
}

func (c *components) empty() bool {
	if c == nil {
		return true
	}
	return len(c.Schemas)+len(c.Responses)+len(c.Parameters)+len(c.Examples)+
		len(c.RequestBodies)+len(c.Headers)+len(c.SecuritySchemes)+
		len(c.Links)+len(c.Callbacks)+len(c.PathItems) == 0
}

// parse decodes one input into a document, separating known top-level
// fields from extras. The token-based parser is used instead of
// json.Unmarshal because we need to capture first-occurrence order of
// unknown keys and the insertion order of `paths`/`webhooks` entries.
func parse(in Input) (*document, error) {
	d := &document{
		name:       in.Name,
		Paths:      newOrderedObject(),
		Webhooks:   newOrderedObject(),
		Components: &components{},
		extras:     newOrderedObject(),
	}
	dec := json.NewDecoder(bytes.NewReader(in.Data))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", in.Name, err)
	}
	if delim, ok := tok.(json.Delim); !ok || delim != '{' {
		return nil, fmt.Errorf("%s: expected JSON object at top level", in.Name)
	}
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", in.Name, err)
		}
		key, ok := tok.(string)
		if !ok {
			return nil, fmt.Errorf("%s: unexpected token %v", in.Name, tok)
		}
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			return nil, fmt.Errorf("%s: %q: %w", in.Name, key, err)
		}
		switch key {
		case "openapi":
			if err := json.Unmarshal(raw, &d.OpenAPI); err != nil {
				return nil, fmt.Errorf("%s: openapi: %w", in.Name, err)
			}
		case "info":
			d.Info = raw
		case "servers":
			d.Servers = raw
		case "paths":
			obj, err := decodeOrderedObject(raw)
			if err != nil {
				return nil, fmt.Errorf("%s: paths: %w", in.Name, err)
			}
			d.Paths = obj
		case "webhooks":
			obj, err := decodeOrderedObject(raw)
			if err != nil {
				return nil, fmt.Errorf("%s: webhooks: %w", in.Name, err)
			}
			d.Webhooks = obj
		case "components":
			if !isJSONNull(raw) {
				if err := json.Unmarshal(raw, d.Components); err != nil {
					return nil, fmt.Errorf("%s: components: %w", in.Name, err)
				}
			}
		case "security":
			if err := json.Unmarshal(raw, &d.Security); err != nil {
				return nil, fmt.Errorf("%s: security: %w", in.Name, err)
			}
		case "tags":
			if err := json.Unmarshal(raw, &d.Tags); err != nil {
				return nil, fmt.Errorf("%s: tags: %w", in.Name, err)
			}
		case "externalDocs":
			d.ExternalDocs = raw
		default:
			d.extras.set(key, raw)
		}
	}
	if _, err := dec.Token(); err != nil {
		return nil, fmt.Errorf("%s: %w", in.Name, err)
	}
	if d.OpenAPI == "" {
		return nil, fmt.Errorf("%s: missing required field \"openapi\"", in.Name)
	}
	if isJSONNull(d.Info) {
		return nil, fmt.Errorf("%s: missing required field \"info\"", in.Name)
	}
	return d, nil
}

// mergeAll merges parsed documents in order. The first establishes the
// values that later inputs must not contradict.
func mergeAll(docs []*document) (*document, error) {
	first := docs[0]
	out := &document{
		name:         "merged",
		OpenAPI:      first.OpenAPI,
		Info:         first.Info,
		Servers:      first.Servers,
		Paths:        newOrderedObject(),
		Webhooks:     newOrderedObject(),
		Components:   &components{},
		ExternalDocs: first.ExternalDocs,
		extras:       newOrderedObject(),
	}
	seenTags := map[string]json.RawMessage{}

	for _, d := range docs {
		if d.OpenAPI != out.OpenAPI {
			return nil, fmt.Errorf("openapi: %s declares %q but %s declares %q",
				first.name, out.OpenAPI, d.name, d.OpenAPI)
		}
		// info/servers/externalDocs and unknown top-level keys are first-wins,
		// silently. The first input's value is already in `out`; nothing to
		// do for later inputs.
		if err := mergeOrdered("paths", out.Paths, d.Paths, d.name); err != nil {
			return nil, err
		}
		if err := mergeOrdered("webhooks", out.Webhooks, d.Webhooks, d.name); err != nil {
			return nil, err
		}
		if err := mergeComponents(out.Components, d.Components, d.name); err != nil {
			return nil, err
		}
		if err := mergeTags(out, seenTags, d); err != nil {
			return nil, err
		}
		if err := mergeSecurity(out, d); err != nil {
			return nil, err
		}
		mergeExtras(out.extras, d.extras)
	}
	return out, nil
}

// mergeOrdered appends entries from src into dst, rejecting key collisions
// whose values differ canonically.
func mergeOrdered(field string, dst, src *orderedObject, srcName string) error {
	for _, k := range src.keys {
		v := src.vals[k]
		if existing, ok := dst.get(k); ok {
			same, err := canonicalEqual(existing, v)
			if err != nil {
				return fmt.Errorf("%s.%s: %w", field, k, err)
			}
			if !same {
				return fmt.Errorf("%s.%q: %s redefines an entry with a different value", field, k, srcName)
			}
			continue
		}
		dst.set(k, v)
	}
	return nil
}

// mergeComponents unions each components/* sub-map. Same-named entries
// across inputs must be canonically identical. The sub-maps are walked in
// spec order so error messages are deterministic.
func mergeComponents(dst, src *components, srcName string) error {
	pairs := []struct {
		name string
		dst  *map[string]json.RawMessage
		src  *map[string]json.RawMessage
	}{
		{"schemas", &dst.Schemas, &src.Schemas},
		{"responses", &dst.Responses, &src.Responses},
		{"parameters", &dst.Parameters, &src.Parameters},
		{"examples", &dst.Examples, &src.Examples},
		{"requestBodies", &dst.RequestBodies, &src.RequestBodies},
		{"headers", &dst.Headers, &src.Headers},
		{"securitySchemes", &dst.SecuritySchemes, &src.SecuritySchemes},
		{"links", &dst.Links, &src.Links},
		{"callbacks", &dst.Callbacks, &src.Callbacks},
		{"pathItems", &dst.PathItems, &src.PathItems},
	}
	for _, p := range pairs {
		if len(*p.src) == 0 {
			continue
		}
		if *p.dst == nil {
			*p.dst = make(map[string]json.RawMessage, len(*p.src))
		}
		for k, v := range *p.src {
			if prev, dup := (*p.dst)[k]; dup {
				same, err := canonicalEqual(prev, v)
				if err != nil {
					return fmt.Errorf("components.%s.%s: %w", p.name, k, err)
				}
				if !same {
					return fmt.Errorf("components.%s.%q: %s redefines an entry with a different value", p.name, k, srcName)
				}
				continue
			}
			(*p.dst)[k] = v
		}
	}
	return nil
}

// mergeTags appends tags from src to out, deduplicating by `name`. Two tag
// entries sharing a name must declare identical metadata.
func mergeTags(out *document, seen map[string]json.RawMessage, src *document) error {
	for _, raw := range src.Tags {
		name, err := tagName(raw)
		if err != nil {
			return fmt.Errorf("%s: tags: %w", src.name, err)
		}
		if prev, ok := seen[name]; ok {
			same, err := canonicalEqual(prev, raw)
			if err != nil {
				return fmt.Errorf("tags[%q]: %w", name, err)
			}
			if !same {
				return fmt.Errorf("tags[%q]: %s redefines a tag with different metadata", name, src.name)
			}
			continue
		}
		seen[name] = raw
		out.Tags = append(out.Tags, raw)
	}
	return nil
}

// mergeSecurity applies first-wins to the root `security` array. The first
// input to declare a non-empty `security` value establishes it; later
// inputs that declare a different non-empty value are an error. (The root
// `security` is a list of alternatives that apply across the whole API; if
// two generators disagree, silently keeping one would change what callers
// are allowed to do.)
func mergeSecurity(out, src *document) error {
	if len(src.Security) == 0 {
		return nil
	}
	if len(out.Security) == 0 {
		out.Security = src.Security
		return nil
	}
	a, err := json.Marshal(out.Security)
	if err != nil {
		return fmt.Errorf("security: %w", err)
	}
	b, err := json.Marshal(src.Security)
	if err != nil {
		return fmt.Errorf("security: %w", err)
	}
	same, err := canonicalEqual(a, b)
	if err != nil {
		return fmt.Errorf("security: %w", err)
	}
	if !same {
		return fmt.Errorf("security: %s declares different requirements", src.name)
	}
	return nil
}

// mergeExtras applies first-wins to unknown top-level keys (notably
// extensions `x-*`). Conflicting redeclarations from later inputs are
// silently ignored, matching the policy for info/servers/etc.
func mergeExtras(dst, src *orderedObject) {
	for _, k := range src.keys {
		if _, ok := dst.get(k); ok {
			continue
		}
		dst.set(k, src.vals[k])
	}
}

// canonicalEqual reports whether a and b decode to the same JSON value
// under canonical (sorted-key) encoding.
func canonicalEqual(a, b json.RawMessage) (bool, error) {
	if bytes.Equal(a, b) {
		return true, nil
	}
	ca, err := canonicalize(a)
	if err != nil {
		return false, err
	}
	cb, err := canonicalize(b)
	if err != nil {
		return false, err
	}
	return bytes.Equal(ca, cb), nil
}

// canonicalize re-encodes raw with sorted map keys, so two semantically
// identical JSON documents that happen to differ in key order compare equal.
func canonicalize(raw json.RawMessage) ([]byte, error) {
	if len(raw) == 0 {
		return raw, nil
	}
	var v any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	return json.Marshal(v)
}

// tagName extracts the `name` field from a tag entry. Tag entries without
// a name violate the OpenAPI spec and are rejected.
func tagName(raw json.RawMessage) (string, error) {
	var t struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &t); err != nil {
		return "", fmt.Errorf("invalid tag entry: %w", err)
	}
	if t.Name == "" {
		return "", errors.New("tag entry missing required \"name\"")
	}
	return t.Name, nil
}

// isJSONNull reports whether raw is the JSON null literal (possibly
// surrounded by whitespace). Empty input also counts as null.
func isJSONNull(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return true
	}
	return bytes.Equal(bytes.TrimSpace(raw), []byte("null"))
}

// orderedObject is an insertion-ordered JSON object. It is used for
// `paths`, `webhooks`, and the bucket of unrecognised top-level keys —
// the three places where input-file order matters and encoding/json's
// alphabetical key sort would lose information.
type orderedObject struct {
	keys []string
	vals map[string]json.RawMessage
}

func newOrderedObject() *orderedObject {
	return &orderedObject{vals: map[string]json.RawMessage{}}
}

func (o *orderedObject) set(k string, v json.RawMessage) {
	if _, ok := o.vals[k]; !ok {
		o.keys = append(o.keys, k)
	}
	o.vals[k] = v
}

func (o *orderedObject) get(k string) (json.RawMessage, bool) {
	if o == nil {
		return nil, false
	}
	v, ok := o.vals[k]
	return v, ok
}

func (o *orderedObject) len() int {
	if o == nil {
		return 0
	}
	return len(o.keys)
}

// MarshalJSON emits the entries in insertion order.
func (o *orderedObject) MarshalJSON() ([]byte, error) {
	if o == nil || len(o.keys) == 0 {
		return []byte("{}"), nil
	}
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, k := range o.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		key, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteByte(':')
		buf.Write(o.vals[k])
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// decodeOrderedObject parses a JSON object into an orderedObject, preserving
// the input's key order.
func decodeOrderedObject(raw json.RawMessage) (*orderedObject, error) {
	if isJSONNull(raw) {
		return newOrderedObject(), nil
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if d, ok := tok.(json.Delim); !ok || d != '{' {
		return nil, fmt.Errorf("expected JSON object, got %v", tok)
	}
	out := newOrderedObject()
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		k, ok := tok.(string)
		if !ok {
			return nil, fmt.Errorf("expected string key, got %v", tok)
		}
		var v json.RawMessage
		if err := dec.Decode(&v); err != nil {
			return nil, err
		}
		out.set(k, v)
	}
	if _, err := dec.Token(); err != nil {
		return nil, err
	}
	return out, nil
}
