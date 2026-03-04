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
var wktSchemas = map[string]*Schema{
	// Timestamp -> RFC 3339 string
	".google.protobuf.Timestamp": {
		Type:   "string",
		Format: "date-time",
	},

	// Duration -> string like "3.5s"
	".google.protobuf.Duration": {
		Type: "string",
	},

	// FieldMask -> comma-separated string
	".google.protobuf.FieldMask": {
		Type: "string",
	},

	// Wrapper types - unwrap to primitives with nullable
	".google.protobuf.StringValue": {
		Type:     "string",
		Nullable: true,
	},
	".google.protobuf.BytesValue": {
		Type:     "string",
		Format:   "byte",
		Nullable: true,
	},
	".google.protobuf.Int32Value": {
		Type:     "integer",
		Format:   "int32",
		Nullable: true,
	},
	".google.protobuf.UInt32Value": {
		Type:     "integer",
		Format:   "int64",
		Nullable: true,
	},
	".google.protobuf.Int64Value": {
		Type:     "string",
		Format:   "int64",
		Nullable: true,
	},
	".google.protobuf.UInt64Value": {
		Type:     "string",
		Format:   "uint64",
		Nullable: true,
	},
	".google.protobuf.FloatValue": {
		Type:     "number",
		Format:   "float",
		Nullable: true,
	},
	".google.protobuf.DoubleValue": {
		Type:     "number",
		Format:   "double",
		Nullable: true,
	},
	".google.protobuf.BoolValue": {
		Type:     "boolean",
		Nullable: true,
	},

	// Empty -> empty object
	".google.protobuf.Empty": {
		Type: "object",
	},

	// Struct -> arbitrary JSON object
	".google.protobuf.Struct": {
		Type: "object",
	},

	// Value -> any JSON value (no type constraint)
	".google.protobuf.Value": {
		// Empty schema allows any type
	},

	// ListValue -> JSON array of any values
	".google.protobuf.ListValue": {
		Type: "array",
		Items: &SchemaRef{
			Value: &Schema{
				// Empty schema for items allows any type
			},
		},
	},

	// NullValue -> represents JSON null
	".google.protobuf.NullValue": {
		Type: "string",
	},

	// Any -> object with @type field
	".google.protobuf.Any": {
		Type: "object",
		Properties: map[string]*SchemaRef{
			"@type": {
				Value: &Schema{
					Type:        "string",
					Description: "A URL/resource name that uniquely identifies the type of the serialized protocol buffer message.",
				},
			},
		},
		AdditionalProperties: &SchemaRef{Value: &Schema{}}, // Allow any additional properties
	},
}
