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
// Uses version-agnostic format: Nullable bool for nullable types.
// Adapters convert to version-specific format (3.0.x: nullable: true, 3.1.0: type array).
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

	// Wrapper types - unwrap to primitives, nullable
	".google.protobuf.StringValue": {
		Type:     SchemaType{"string"},
		Nullable: true,
	},
	".google.protobuf.BytesValue": {
		Type:     SchemaType{"string"},
		Format:   "byte",
		Nullable: true,
	},
	".google.protobuf.Int32Value": {
		Type:     SchemaType{"integer"},
		Format:   "int32",
		Nullable: true,
	},
	".google.protobuf.UInt32Value": {
		Type:     SchemaType{"integer"},
		Format:   "int64",
		Nullable: true,
	},
	".google.protobuf.Int64Value": {
		Type:     SchemaType{"string"},
		Format:   "int64",
		Nullable: true,
	},
	".google.protobuf.UInt64Value": {
		Type:     SchemaType{"string"},
		Format:   "uint64",
		Nullable: true,
	},
	".google.protobuf.FloatValue": {
		Type:     SchemaType{"number"},
		Format:   "float",
		Nullable: true,
	},
	".google.protobuf.DoubleValue": {
		Type:     SchemaType{"number"},
		Format:   "double",
		Nullable: true,
	},
	".google.protobuf.BoolValue": {
		Type:     SchemaType{"boolean"},
		Nullable: true,
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
