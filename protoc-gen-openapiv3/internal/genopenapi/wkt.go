package genopenapi

// wktSchemas maps fully-qualified protobuf well-known type names (with leading
// dot, as proto descriptors return them) to the inline OpenAPI schema we emit
// when a field references that type.
//
// All WKTs are inlined; this generator does not emit reusable component
// schemas for them. Behavior matches protojson's serialization model:
//   - Timestamp / Duration / FieldMask: string forms
//   - Wrapper types: their underlying primitive
//   - Empty / Struct: object
//   - Value: any (no type constraint)
//   - ListValue: array of any
//   - NullValue: literal null
//   - Any: object with @type plus open additional properties
//
// Wrapper types are intentionally not marked nullable. The strictly-correct
// JSON Schema 2020-12 form is `type: [<primitive>, "null"]`, but no Go
// OpenAPI generator in the ecosystem handles that form yet: openapi-
// generator-cli writes a literal `anyOf<string,integer>` placeholder for
// it, and oapi-codegen errors with "unhandled Schema type: &[string null]"
// on request bodies that reference it. Until tooling catches up we
// describe wrappers as their underlying primitive; the gateway runtime
// still accepts a JSON null on the wire because protojson treats wrappers
// as optional regardless.
//
// Each call should return a fresh schema so callers can mutate it without
// risk of cross-contamination.
func wellKnownTypeSchema(typeName string) *Schema {
	switch typeName {
	case ".google.protobuf.Timestamp":
		return &Schema{Type: SchemaType{"string"}, Format: "date-time"}
	case ".google.protobuf.Duration":
		return &Schema{Type: SchemaType{"string"}}
	case ".google.protobuf.FieldMask":
		return &Schema{Type: SchemaType{"string"}}

	case ".google.protobuf.StringValue":
		return &Schema{Type: SchemaType{"string"}}
	case ".google.protobuf.BytesValue":
		return &Schema{Type: SchemaType{"string"}, Format: "byte"}
	case ".google.protobuf.Int32Value":
		return &Schema{Type: SchemaType{"integer"}, Format: "int32"}
	case ".google.protobuf.UInt32Value":
		// uint32's range (0..4294967295) exceeds int32; OpenAPI has no
		// "uint32" format, so int64 is the narrowest that fits.
		return &Schema{Type: SchemaType{"integer"}, Format: "int64"}
	case ".google.protobuf.Int64Value":
		return &Schema{Type: SchemaType{"string"}, Format: "int64"}
	case ".google.protobuf.UInt64Value":
		return &Schema{Type: SchemaType{"string"}, Format: "uint64"}
	case ".google.protobuf.FloatValue":
		return &Schema{Type: SchemaType{"number"}, Format: "float"}
	case ".google.protobuf.DoubleValue":
		return &Schema{Type: SchemaType{"number"}, Format: "double"}
	case ".google.protobuf.BoolValue":
		return &Schema{Type: SchemaType{"boolean"}}

	case ".google.protobuf.Empty":
		return &Schema{Type: SchemaType{"object"}}
	case ".google.protobuf.Struct":
		return &Schema{Type: SchemaType{"object"}}
	case ".google.protobuf.Value":
		return &Schema{}
	case ".google.protobuf.ListValue":
		return &Schema{
			Type:  SchemaType{"array"},
			Items: &SchemaOrRef{Value: &Schema{}},
		}
	case ".google.protobuf.NullValue":
		return &Schema{Type: SchemaType{"null"}}

	case ".google.protobuf.Any":
		t := true
		return &Schema{
			Type: SchemaType{"object"},
			Properties: map[string]*SchemaOrRef{
				"@type": {Value: &Schema{
					Type:        SchemaType{"string"},
					Description: "A URL/resource name that uniquely identifies the type of the serialized protocol buffer message.",
				}},
			},
			AdditionalProperties: &AdditionalProperties{Bool: &t},
		}
	}
	return nil
}
