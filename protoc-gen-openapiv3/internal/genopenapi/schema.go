package genopenapi

import (
	"log"
	"slices"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// schemaBuilder owns the in-progress component schema set for one document.
// Construction is intentionally lazy: schemas are only added when something
// references them, transitively from RPC request/response types. We never
// emit unreferenced messages.
type schemaBuilder struct {
	reg *descriptor.Registry
	doc *Document
}

func newSchemaBuilder(reg *descriptor.Registry, doc *Document) *schemaBuilder {
	return &schemaBuilder{
		reg: reg,
		doc: doc,
	}
}

// fieldSchema returns the schema (or $ref) describing the given proto field's
// type. For repeated fields it produces an array; for map entries it produces
// an object with additionalProperties; for messages and enums it produces a
// $ref into components and ensures the referenced schema is generated.
func (b *schemaBuilder) fieldSchema(field *descriptor.Field) *SchemaOrRef {
	if field.GetLabel() != descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		return b.scalarOrRef(field)
	}
	// Map fields are repeated synthetic messages with map_entry=true. Resolve
	// the message once: a hit on map_entry takes the map path, a hit on a
	// regular message hands the resolved descriptor to scalarOrRef so it does
	// not have to look it up again, and a miss falls through to the
	// unresolved branch which logs and emits an open object.
	if field.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		msg, err := b.reg.LookupMsg("", field.GetTypeName())
		if err == nil {
			if msg.GetOptions().GetMapEntry() {
				return b.mapSchema(msg)
			}
			return &SchemaOrRef{
				Value: &Schema{
					Type:  SchemaType{"array"},
					Items: b.messageRef(msg),
				},
			}
		}
	}
	return &SchemaOrRef{
		Value: &Schema{
			Type:  SchemaType{"array"},
			Items: b.scalarOrRef(field),
		},
	}
}

// messageRef returns a $ref to the given message's component schema, ensuring
// the component is generated. Well-known types are inlined instead.
func (b *schemaBuilder) messageRef(msg *descriptor.Message) *SchemaOrRef {
	if wkt := wellKnownTypeSchema(msg.FQMN()); wkt != nil {
		return &SchemaOrRef{Value: wkt}
	}
	b.ensureMessageSchema(msg)
	return NewSchemaRef(schemaName(msg.FQMN()))
}

// mapSchema returns the additionalProperties schema for a proto map field.
// Map entries are messages with two fields:
// key (1) and value (2); the value field's type drives additionalProperties.
func (b *schemaBuilder) mapSchema(msg *descriptor.Message) *SchemaOrRef {
	var valueField *descriptor.Field
	for _, f := range msg.Fields {
		if f.GetNumber() == 2 {
			valueField = f
			break
		}
	}
	if valueField == nil {
		return nil
	}
	return &SchemaOrRef{
		Value: &Schema{
			Type: SchemaType{"object"},
			AdditionalProperties: &AdditionalProperties{
				Schema: b.scalarOrRef(valueField),
			},
		},
	}
}

// scalarOrRef produces the schema for a non-repeated field type. Message and
// enum types resolve to $refs into components; everything else is an inline
// scalar schema.
func (b *schemaBuilder) scalarOrRef(field *descriptor.Field) *SchemaOrRef {
	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return &SchemaOrRef{Value: &Schema{Type: SchemaType{"string"}}}
	case descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		return &SchemaOrRef{Value: &Schema{Type: SchemaType{"integer"}, Format: "int32"}}
	case descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		// 64-bit ints are JSON strings per protojson.
		return &SchemaOrRef{Value: &Schema{Type: SchemaType{"string"}, Format: "int64"}}
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		// OpenAPI has no "uint32" format; int64 is the narrowest that
		// covers the full uint32 range.
		return &SchemaOrRef{Value: &Schema{Type: SchemaType{"integer"}, Format: "int64"}}
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		return &SchemaOrRef{Value: &Schema{Type: SchemaType{"string"}, Format: "uint64"}}
	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		return &SchemaOrRef{Value: &Schema{Type: SchemaType{"number"}, Format: "float"}}
	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return &SchemaOrRef{Value: &Schema{Type: SchemaType{"number"}, Format: "double"}}
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return &SchemaOrRef{Value: &Schema{Type: SchemaType{"boolean"}}}
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return &SchemaOrRef{Value: &Schema{Type: SchemaType{"string"}, Format: "byte"}}
	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		// Well-known types are inlined; everything else is referenced.
		if wkt := wellKnownTypeSchema(field.GetTypeName()); wkt != nil {
			return &SchemaOrRef{Value: wkt}
		}
		msg, err := b.reg.LookupMsg("", field.GetTypeName())
		if err != nil {
			log.Printf("protoc-gen-openapiv3: cannot resolve message %q for field %q: %v; emitting open object",
				field.GetTypeName(), field.GetName(), err)
			return &SchemaOrRef{Value: &Schema{Type: SchemaType{"object"}}}
		}
		b.ensureMessageSchema(msg)
		return NewSchemaRef(schemaName(msg.FQMN()))
	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		enum, err := b.reg.LookupEnum("", field.GetTypeName())
		if err != nil {
			log.Printf("protoc-gen-openapiv3: cannot resolve enum %q for field %q: %v; emitting plain string",
				field.GetTypeName(), field.GetName(), err)
			return &SchemaOrRef{Value: &Schema{Type: SchemaType{"string"}}}
		}
		b.ensureEnumSchema(enum)
		return NewSchemaRef(schemaName(enum.FQEN()))
	default:
		// Fall back to a string value on unhandled field types
		return &SchemaOrRef{Value: &Schema{Type: SchemaType{"string"}}}
	}
}

// ensureMessageSchema generates a component schema for the given message if
// it has not already been emitted. Map-entry messages are skipped.
func (b *schemaBuilder) ensureMessageSchema(msg *descriptor.Message) {
	if msg.GetOptions().GetMapEntry() {
		// Map entries are handled separately
		return
	}
	name := schemaName(msg.FQMN())
	if _, exists := b.doc.Components.Schemas[name]; exists {
		// Message already in schema, skip
		return
	}

	// Reserve the slot up front so cycles terminate: when a nested field
	// recurses into a message we're already building, the exists check
	// above short-circuits instead of looping.
	schema := &Schema{
		Type:       SchemaType{"object"},
		Deprecated: messageDeprecated(msg),
	}
	if desc := messageComments(b.reg, msg); desc != "" {
		schema.Description = desc
	}
	b.doc.Components.Schemas[name] = &SchemaOrRef{Value: schema}

	// Partition fields into oneof groups and regular fields. Synthetic
	// proto3-optional oneofs are treated as regular fields.
	regular, oneofs := splitOneofs(msg)

	total := len(regular)
	for _, g := range oneofs {
		total += len(g.fields)
	}
	if total > 0 {
		schema.Properties = make(map[string]*SchemaOrRef, total)
	}
	for _, field := range regular {
		b.addProperty(schema, field)
	}
	for _, g := range oneofs {
		for _, field := range g.fields {
			b.addProperty(schema, field)
		}
	}

	// Each proto oneof group constrains its fields to "at most one set"
	// (proto3 oneof allows zero or one). Encode that as a JSON Schema
	// `oneOf` whose options are:
	//
	//   1. a "none of the fields are set" guard, expressed as
	//      {type: object, not: {anyOf: [{required: F1}, ..., {required: Fn}]}}
	//   2. one option per field carrying that field's own typed schema:
	//      {type: object, properties: {Fi: <fieldSchema>}, required: [Fi]}
	//
	// `oneOf` requires exactly one sub-schema to match, so:
	//   - zero set      → only the guard matches            → passes
	//   - exactly one Fi → only the Fi option matches        → passes
	//   - two or more   → multiple Fi options match         → fails
	//
	// Each per-field option embeds the field's full schema (a $ref for
	// nested message types, the inlined schema for scalars and WKTs),
	// not just `required`.
	//
	// The same property schemas also live on the parent's top-level
	// `properties` map, so consumers that don't walk into `oneOf` still
	// see the full set of fields. The duplication is intentional —
	// `properties`, `required`, and `oneOf` are orthogonal in JSON
	// Schema and tooling tends to consult them independently.
	//
	// Single-field oneof groups produce a trivially-true constraint
	// ("Fi is either present or not present"); we skip those for output
	// hygiene. proto3-optional fields use synthetic single-field oneofs
	// and are already filtered by splitOneofs.
	//
	// Multiple groups are independent constraints, encoded as an `allOf`
	// of one `oneOf` per group. A single non-trivial group is hoisted
	// directly onto the schema's `oneOf` for less verbose output.
	groupConstraints := make([]*SchemaOrRef, 0, len(oneofs))
	for _, g := range oneofs {
		if len(g.fields) < 2 {
			// Single-field oneof groups are a no-op constraint, so skip them.
			continue
		}
		anyRequired := make([]*SchemaOrRef, 0, len(g.fields))
		options := make([]*SchemaOrRef, 0, len(g.fields)+1)
		for _, field := range g.fields {
			anyRequired = append(anyRequired, &SchemaOrRef{Value: &Schema{
				Required: []string{jsonName(field)},
			}})
		}
		options = append(options, &SchemaOrRef{Value: &Schema{
			Type: SchemaType{"object"},
			Not:  &SchemaOrRef{Value: &Schema{AnyOf: anyRequired}},
		}})
		for _, field := range g.fields {
			fname := jsonName(field)
			options = append(options, &SchemaOrRef{Value: &Schema{
				Type: SchemaType{"object"},
				Properties: map[string]*SchemaOrRef{
					fname: b.fieldSchema(field),
				},
				Required: []string{fname},
			}})
		}
		groupConstraints = append(groupConstraints, &SchemaOrRef{Value: &Schema{OneOf: options}})
	}

	switch len(groupConstraints) {
	case 0:
		// No oneof groups, nothing to do.
	case 1:
		// A single oneof group can be hoisted directly onto the schema for less verbose output.
		schema.OneOf = groupConstraints[0].Value.OneOf
	default:
		// Multiple groups are independent constraints, so combine them with allOf.
		schema.AllOf = groupConstraints
	}
}

// addProperty inserts a single field as a property of schema, attaching
// description (via $ref sibling for refs, or directly for inline schemas)
// and updating Required from field_behavior.
func (b *schemaBuilder) addProperty(schema *Schema, field *descriptor.Field) {
	if schema.Properties == nil {
		schema.Properties = make(map[string]*SchemaOrRef)
	}
	name := jsonName(field)
	prop := b.propertySchema(field)
	schema.Properties[name] = prop

	if hasFieldBehavior(field, annotations.FieldBehavior_REQUIRED) {
		schema.Required = append(schema.Required, name)
	}
}

// propertySchema produces a property schema for the given field, with
// description and field-behavior flags applied. Description goes on the
// $ref sibling for referenced types, or directly on inline schemas.
func (b *schemaBuilder) propertySchema(field *descriptor.Field) *SchemaOrRef {
	prop := b.fieldSchema(field)
	desc := fieldComments(b.reg, field)

	if prop.Ref != "" {
		// 3.1.0 allows description as a sibling of $ref.
		if desc != "" {
			prop.Description = desc
		}
		// readOnly/writeOnly/deprecated next to a $ref also work as
		// siblings in 3.1.0; for those we wrap in allOf so the consumer
		// sees a real schema body to attach the flag to.
		if needsAllOfWrap(field) {
			wrapped := &Schema{AllOf: []*SchemaOrRef{{Ref: prop.Ref, Description: prop.Description}}}
			applyFieldFlags(wrapped, field)
			return &SchemaOrRef{Value: wrapped}
		}
		return prop
	}

	if prop.Value != nil {
		if desc != "" {
			prop.Value.Description = desc
		}
		applyFieldFlags(prop.Value, field)
	}
	return prop
}

// needsAllOfWrap reports whether a referenced field needs an allOf wrapper to
// carry per-occurrence flags (read-only / write-only / deprecated). A pure
// description does not.
func needsAllOfWrap(field *descriptor.Field) bool {
	if fieldDeprecated(field) {
		return true
	}
	for _, fb := range fieldBehaviors(field) {
		switch fb {
		case annotations.FieldBehavior_OUTPUT_ONLY, annotations.FieldBehavior_INPUT_ONLY:
			return true
		}
	}
	return false
}

// applyFieldFlags writes per-field schema-level flags (deprecated, readOnly,
// writeOnly) onto an inline schema. Required is handled at the parent.
func applyFieldFlags(s *Schema, field *descriptor.Field) {
	if fieldDeprecated(field) {
		s.Deprecated = true
	}
	for _, fb := range fieldBehaviors(field) {
		switch fb {
		case annotations.FieldBehavior_OUTPUT_ONLY:
			s.ReadOnly = true
		case annotations.FieldBehavior_INPUT_ONLY:
			s.WriteOnly = true
		}
	}
}

// fieldDeprecated reports whether a field should be flagged deprecated in
// the emitted schema. The flag cascades from the enclosing message and
// file: a field is deprecated if the field itself, its containing message,
// or the proto file that declares it is marked with `deprecated = true`.
func fieldDeprecated(field *descriptor.Field) bool {
	return field.GetOptions().GetDeprecated() ||
		field.Message.GetOptions().GetDeprecated() ||
		field.Message.File.GetOptions().GetDeprecated()
}

// methodDeprecated reports whether an RPC method should be flagged
// deprecated. The flag cascades from the service and file: a method is
// deprecated if the method itself, its service, or the proto file is
// marked deprecated. OpenAPI 3.1 has no `deprecated` flag on tags, so
// cascading service-level deprecation into every method of the service
// is the only way to reflect it.
func methodDeprecated(m *descriptor.Method) bool {
	return m.GetOptions().GetDeprecated() ||
		m.Service.GetOptions().GetDeprecated() ||
		m.Service.File.GetOptions().GetDeprecated()
}

// messageDeprecated reports whether a message component schema should
// be flagged deprecated. The flag cascades from the file: a message is
// deprecated if the message itself or the proto file is marked
// deprecated. Nested types are left alone — they're separately-named
// component schemas and deprecate independently of any outer message.
func messageDeprecated(msg *descriptor.Message) bool {
	return msg.GetOptions().GetDeprecated() ||
		msg.File.GetOptions().GetDeprecated()
}

// enumDeprecated reports whether an enum component schema should be
// flagged deprecated. The flag cascades from the file.
func enumDeprecated(enum *descriptor.Enum) bool {
	return enum.GetOptions().GetDeprecated() ||
		enum.File.GetOptions().GetDeprecated()
}

// ensureEnumSchema generates a component schema for the given enum if it
// has not been emitted yet. Enum values are rendered as their string names.
//
// The grpc-gateway runtime (via protojson) actually accepts enum values on
// the wire as either string names or integer numbers, but encoding that
// dual acceptance in the schema breaks consumer tooling: openapi-generator-
// cli's Go target produces an unbuildable wrapper type for both the
// `oneOf`-of-homogeneous-enums form (pulls in gopkg.in/validator.v2 plus a
// painful discriminator wrapper) and the `type: [string, integer]`
// mixed-type form (literally writes `anyOf<string,integer>` as the Go type
// name — an unimplemented codegen branch). Until tooling catches up, we
// document the string form in the spec and rely on protojson's leniency
// to accept the integer form at the gateway boundary.
func (b *schemaBuilder) ensureEnumSchema(enum *descriptor.Enum) {
	name := schemaName(enum.FQEN())
	if _, exists := b.doc.Components.Schemas[name]; exists {
		return
	}
	values := make([]any, 0, len(enum.Value))
	for _, v := range enum.Value {
		values = append(values, v.GetName())
	}
	schema := &Schema{
		Type:       SchemaType{"string"},
		Enum:       values,
		Deprecated: enumDeprecated(enum),
	}
	if desc := enumComments(b.reg, enum); desc != "" {
		schema.Description = desc
	}
	b.doc.Components.Schemas[name] = &SchemaOrRef{Value: schema}
}

// oneofGroup is a oneof declaration plus its constituent fields.
type oneofGroup struct {
	name   string
	fields []*descriptor.Field
}

// splitOneofs separates regular fields from real (non-synthetic) oneof groups
// in declaration order. proto3-optional fields use synthetic single-field
// oneofs which we treat as regular optional fields.
func splitOneofs(msg *descriptor.Message) (regular []*descriptor.Field, groups []oneofGroup) {
	byIndex := make(map[int32]*oneofGroup)
	for _, field := range msg.Fields {
		if field.OneofIndex != nil && !field.GetProto3Optional() {
			idx := field.GetOneofIndex()
			grp, ok := byIndex[idx]
			if !ok {
				grp = &oneofGroup{name: msg.GetOneofDecl()[idx].GetName()}
				byIndex[idx] = grp
			}
			grp.fields = append(grp.fields, field)
			continue
		}
		regular = append(regular, field)
	}
	for i := range msg.GetOneofDecl() {
		if grp, ok := byIndex[int32(i)]; ok {
			groups = append(groups, *grp)
		}
	}
	return regular, groups
}

// fieldBehaviors returns the [(google.api.field_behavior) = ...] entries on a
// field, or nil if there are none.
func fieldBehaviors(field *descriptor.Field) []annotations.FieldBehavior {
	if field.Options == nil {
		return nil
	}
	if !proto.HasExtension(field.Options, annotations.E_FieldBehavior) {
		return nil
	}
	out, _ := proto.GetExtension(field.Options, annotations.E_FieldBehavior).([]annotations.FieldBehavior)
	return out
}

func hasFieldBehavior(field *descriptor.Field, target annotations.FieldBehavior) bool {
	return slices.Contains(fieldBehaviors(field), target)
}

// jsonName returns the JSON name we use for a proto field. We always honor
// the proto json_name option since the runtime uses
// JSON names by default and we'd rather match what the wire produces.
func jsonName(field *descriptor.Field) string {
	if name := field.GetJsonName(); name != "" {
		return name
	}
	return field.GetName()
}
