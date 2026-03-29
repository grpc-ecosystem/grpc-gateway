package genopenapiv3

// wellKnownTypeSchema returns OpenAPI schema for protobuf well-known types.
// The schemas match the behavior of the JSON unmarshaler in protobuf-go.
// See: https://github.com/protocolbuffers/protobuf-go/blob/main/encoding/protojson/well_known_types.go
func wellKnownTypeSchema(typeName string) *Schema {
	schema, ok := wktSchemas[typeName]
	if !ok {
		return nil
	}
	// Return a copy to prevent mutations
	schemaCopy := *schema
	return &schemaCopy
}

// wktSchemas maps well-known type names to their OpenAPI schema representations.
// Uses OpenAPI 3.1.0 style: type arrays for nullable (e.g., ["string", "null"]).
var wktSchemas = map[string]*Schema{
	// Timestamp -> RFC 3339 string
	".google.protobuf.Timestamp": {
		Type:   SchemaType{"string"},
		Format: "date-time",
	},

	// Duration -> string like "3.5s"
	".google.protobuf.Duration": {
		Type: SchemaType{"string"},
	},

	// FieldMask -> comma-separated string
	".google.protobuf.FieldMask": {
		Type: SchemaType{"string"},
	},

	// Wrapper types - unwrap to primitives, nullable via type array (3.1.0 style)
	".google.protobuf.StringValue": {
		Type: SchemaType{"string", "null"},
	},
	".google.protobuf.BytesValue": {
		Type:   SchemaType{"string", "null"},
		Format: "byte",
	},
	".google.protobuf.Int32Value": {
		Type:   SchemaType{"integer", "null"},
		Format: "int32",
	},
	".google.protobuf.UInt32Value": {
		Type:   SchemaType{"integer", "null"},
		Format: "int64",
	},
	".google.protobuf.Int64Value": {
		Type:   SchemaType{"string", "null"},
		Format: "int64",
	},
	".google.protobuf.UInt64Value": {
		Type:   SchemaType{"string", "null"},
		Format: "uint64",
	},
	".google.protobuf.FloatValue": {
		Type:   SchemaType{"number", "null"},
		Format: "float",
	},
	".google.protobuf.DoubleValue": {
		Type:   SchemaType{"number", "null"},
		Format: "double",
	},
	".google.protobuf.BoolValue": {
		Type: SchemaType{"boolean", "null"},
	},

	// Empty -> empty object
	".google.protobuf.Empty": {
		Type: SchemaType{"object"},
	},

	// Struct -> arbitrary JSON object
	".google.protobuf.Struct": {
		Type: SchemaType{"object"},
	},

	// Value -> any JSON value (no type constraint)
	".google.protobuf.Value": {
		// Empty schema allows any type
	},

	// ListValue -> JSON array of any values
	".google.protobuf.ListValue": {
		Type: SchemaType{"array"},
		Items: &SchemaOrReference{
			Schema: &Schema{
				// Empty schema for items allows any type
			},
		},
	},

	// NullValue -> represents JSON null
	".google.protobuf.NullValue": {
		Type: SchemaType{"null"},
	},

	// Any -> object with @type field
	".google.protobuf.Any": {
		Type: SchemaType{"object"},
		Properties: map[string]*SchemaOrReference{
			"@type": {
				Schema: &Schema{
					Type:        SchemaType{"string"},
					Description: "A URL/resource name that uniquely identifies the type of the serialized protocol buffer message.",
				},
			},
		},
		AdditionalProperties: &SchemaOrReference{Schema: &Schema{}}, // Allow any additional properties
	},
}
