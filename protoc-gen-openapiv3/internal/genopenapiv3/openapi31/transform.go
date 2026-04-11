package openapi31

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapiv3/model"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapiv3/transform"
)

// Transformer embeds BaseTransformer to inherit shared transformation methods.
type Transformer struct {
	transform.BaseTransformer
}

// NewTransformer creates a new OpenAPI 3.1.0 transformer.
func NewTransformer() *Transformer {
	return &Transformer{}
}

// TransformDocument converts a canonical Document to OpenAPI 3.1.0 output format.
func TransformDocument(doc *model.Document) *Document {
	if doc == nil {
		return nil
	}

	t := NewTransformer()
	return &Document{
		OpenAPI:      "3.1.0",
		Info:         t.TransformInfo(doc.Info),
		Servers:      t.TransformServers(doc.Servers),
		Paths:        t.transformPaths(doc.Paths),
		Components:   t.transformComponents(doc.Components),
		Security:     t.TransformSecurityRequirements(doc.Security),
		Tags:         t.TransformTags(doc.Tags),
		ExternalDocs: t.TransformExternalDocs(doc.ExternalDocs),
	}
}

func (t *Transformer) transformPaths(paths *model.Paths) *Paths {
	if paths == nil || len(paths.Items) == 0 {
		return nil
	}
	result := &Paths{
		Items: make(map[string]*PathItem),
	}
	for path, item := range paths.Items {
		result.Items[path] = t.transformPathItem(item)
	}
	return result
}

func (t *Transformer) transformPathItem(item *model.PathItem) *PathItem {
	if item == nil {
		return nil
	}
	return &PathItem{
		Ref:         item.Ref,
		Summary:     item.Summary,
		Description: item.Description,
		Get:         t.transformOperation(item.Get),
		Put:         t.transformOperation(item.Put),
		Post:        t.transformOperation(item.Post),
		Delete:      t.transformOperation(item.Delete),
		Options:     t.transformOperation(item.Options),
		Head:        t.transformOperation(item.Head),
		Patch:       t.transformOperation(item.Patch),
		Trace:       t.transformOperation(item.Trace),
		Servers:     t.TransformServers(item.Servers),
		Parameters:  t.transformParameterOrRefs(item.Parameters),
	}
}

func (t *Transformer) transformOperation(op *model.Operation) *Operation {
	if op == nil {
		return nil
	}
	return &Operation{
		Tags:         op.Tags,
		Summary:      op.Summary,
		Description:  op.Description,
		ExternalDocs: t.TransformExternalDocs(op.ExternalDocs),
		OperationID:  op.OperationID,
		Parameters:   t.transformParameterOrRefs(op.Parameters),
		RequestBody:  t.transformRequestBodyOrRef(op.RequestBody),
		Responses:    t.transformResponses(op.Responses),
		Callbacks:    t.transformCallbacks(op.Callbacks),
		Deprecated:   op.Deprecated,
		Security:     t.TransformSecurityRequirements(op.Security),
		Servers:      t.TransformServers(op.Servers),
	}
}

func (t *Transformer) transformParameterOrRefs(params []*model.ParameterOrRef) []*ParameterOrRef {
	if len(params) == 0 {
		return nil
	}
	result := make([]*ParameterOrRef, len(params))
	for i, p := range params {
		result[i] = t.transformParameterOrRef(p)
	}
	return result
}

func (t *Transformer) transformParameterOrRef(p *model.ParameterOrRef) *ParameterOrRef {
	if p == nil {
		return nil
	}
	if p.Ref != "" {
		return &ParameterOrRef{Ref: p.Ref}
	}
	if p.Value == nil {
		return nil
	}
	return &ParameterOrRef{
		Value: &Parameter{
			Name:            p.Value.Name,
			In:              p.Value.In,
			Description:     p.Value.Description,
			Required:        p.Value.Required,
			Deprecated:      p.Value.Deprecated,
			AllowEmptyValue: p.Value.AllowEmptyValue,
			Style:           p.Value.Style,
			Explode:         p.Value.Explode,
			AllowReserved:   p.Value.AllowReserved,
			Schema:          t.transformSchemaOrRef(p.Value.Schema),
			Examples:        t.transformExamples(p.Value.Examples),
			Content:         t.transformContent(p.Value.Content),
		},
	}
}

func (t *Transformer) transformRequestBodyOrRef(rb *model.RequestBodyOrRef) *RequestBodyOrRef {
	if rb == nil {
		return nil
	}
	if rb.Ref != "" {
		return &RequestBodyOrRef{Ref: rb.Ref}
	}
	if rb.Value == nil {
		return nil
	}
	return &RequestBodyOrRef{
		Value: &RequestBody{
			Description: rb.Value.Description,
			Content:     t.transformContent(rb.Value.Content),
			Required:    rb.Value.Required,
		},
	}
}

func (t *Transformer) transformContent(content map[string]*model.MediaType) map[string]*MediaType {
	if len(content) == 0 {
		return nil
	}
	result := make(map[string]*MediaType)
	for mt, media := range content {
		result[mt] = t.transformMediaType(media)
	}
	return result
}

func (t *Transformer) transformMediaType(mt *model.MediaType) *MediaType {
	if mt == nil {
		return nil
	}
	return &MediaType{
		Schema:   t.transformSchemaOrRef(mt.Schema),
		Examples: t.transformExamples(mt.Examples),
		Encoding: t.transformEncoding(mt.Encoding),
	}
}

func (t *Transformer) transformEncoding(encoding map[string]*model.Encoding) map[string]*Encoding {
	if len(encoding) == 0 {
		return nil
	}
	result := make(map[string]*Encoding)
	for name, enc := range encoding {
		result[name] = &Encoding{
			ContentType:   enc.ContentType,
			Headers:       t.transformHeaderOrRefs(enc.Headers),
			Style:         enc.Style,
			Explode:       enc.Explode,
			AllowReserved: enc.AllowReserved,
		}
	}
	return result
}

func (t *Transformer) transformResponses(r *model.Responses) *Responses {
	if r == nil {
		return nil
	}
	result := &Responses{
		Default: t.transformResponseOrRef(r.Default),
		Codes:   make(map[string]*ResponseOrRef),
	}
	for code, resp := range r.Codes {
		result.Codes[code] = t.transformResponseOrRef(resp)
	}
	return result
}

func (t *Transformer) transformResponseOrRef(r *model.ResponseOrRef) *ResponseOrRef {
	if r == nil {
		return nil
	}
	if r.Ref != "" {
		return &ResponseOrRef{Ref: r.Ref}
	}
	if r.Value == nil {
		return nil
	}
	return &ResponseOrRef{
		Value: &Response{
			Description: r.Value.Description,
			Headers:     t.transformHeaderOrRefs(r.Value.Headers),
			Content:     t.transformContent(r.Value.Content),
			Links:       t.transformLinkOrRefs(r.Value.Links),
		},
	}
}

func (t *Transformer) transformHeaderOrRefs(headers map[string]*model.HeaderOrRef) map[string]*HeaderOrRef {
	if len(headers) == 0 {
		return nil
	}
	result := make(map[string]*HeaderOrRef)
	for name, h := range headers {
		result[name] = t.transformHeaderOrRef(h)
	}
	return result
}

func (t *Transformer) transformHeaderOrRef(h *model.HeaderOrRef) *HeaderOrRef {
	if h == nil {
		return nil
	}
	if h.Ref != "" {
		return &HeaderOrRef{Ref: h.Ref}
	}
	if h.Value == nil {
		return nil
	}
	return &HeaderOrRef{
		Value: &Header{
			Description:     h.Value.Description,
			Required:        h.Value.Required,
			Deprecated:      h.Value.Deprecated,
			AllowEmptyValue: h.Value.AllowEmptyValue,
			Style:           h.Value.Style,
			Explode:         h.Value.Explode,
			AllowReserved:   h.Value.AllowReserved,
			Schema:          t.transformSchemaOrRef(h.Value.Schema),
			Examples:        t.transformExamples(h.Value.Examples),
			Content:         t.transformContent(h.Value.Content),
		},
	}
}

func (t *Transformer) transformLinkOrRefs(links map[string]*model.LinkOrRef) map[string]*LinkOrRef {
	if len(links) == 0 {
		return nil
	}
	result := make(map[string]*LinkOrRef)
	for name, l := range links {
		result[name] = t.transformLinkOrRef(l)
	}
	return result
}

func (t *Transformer) transformLinkOrRef(l *model.LinkOrRef) *LinkOrRef {
	if l == nil {
		return nil
	}
	if l.Ref != "" {
		return &LinkOrRef{Ref: l.Ref}
	}
	if l.Value == nil {
		return nil
	}
	return &LinkOrRef{
		Value: &Link{
			OperationRef: l.Value.OperationRef,
			OperationID:  l.Value.OperationID,
			Parameters:   l.Value.Parameters,
			RequestBody:  l.Value.RequestBody,
			Description:  l.Value.Description,
			Server:       t.TransformServer(l.Value.Server),
		},
	}
}

func (t *Transformer) transformCallbacks(callbacks map[string]*model.CallbackOrRef) map[string]*CallbackOrRef {
	if len(callbacks) == 0 {
		return nil
	}
	result := make(map[string]*CallbackOrRef)
	for name, cb := range callbacks {
		result[name] = t.transformCallbackOrRef(cb)
	}
	return result
}

func (t *Transformer) transformCallbackOrRef(cb *model.CallbackOrRef) *CallbackOrRef {
	if cb == nil {
		return nil
	}
	if cb.Ref != "" {
		return &CallbackOrRef{Ref: cb.Ref}
	}
	if cb.Value == nil {
		return nil
	}
	pathItems := make(map[string]*PathItem)
	for expr, item := range cb.Value {
		pathItems[expr] = t.transformPathItem(item)
	}
	return &CallbackOrRef{Value: pathItems}
}

// transformSchemaOrRef transforms a canonical SchemaOrRef to 3.1.0 format.
// In OpenAPI 3.1.0, $ref can have sibling summary/description.
func (t *Transformer) transformSchemaOrRef(s *model.SchemaOrRef) *SchemaOrRef {
	if s == nil {
		return nil
	}
	if s.Ref != "" {
		// 3.1.0 allows $ref with sibling summary/description
		return &SchemaOrRef{
			Ref:         s.Ref,
			Summary:     s.Summary,
			Description: s.Description,
		}
	}
	if s.Value == nil {
		return nil
	}
	return &SchemaOrRef{
		Value: t.transformSchema(s.Value),
	}
}

// transformSchema transforms a canonical Schema to 3.1.0 format.
// Key difference: nullable is expressed as type array ["type", "null"] in 3.1.0.
func (t *Transformer) transformSchema(s *model.Schema) *Schema {
	if s == nil {
		return nil
	}

	// Handle nullable via type array in 3.1.0
	var typeVal any
	if s.Type != "" {
		if s.IsNullable {
			// 3.1.0 style: type array with null
			typeVal = []string{s.Type, "null"}
		} else {
			typeVal = s.Type
		}
	}

	return &Schema{
		Type:                 typeVal,
		Format:               s.Format,
		Title:                s.Title,
		Description:          s.Description,
		Default:              s.Default,
		Examples:             s.Examples,
		Deprecated:           s.Deprecated,
		ReadOnly:             s.ReadOnly,
		WriteOnly:            s.WriteOnly,
		ExternalDocs:         t.TransformExternalDocs(s.ExternalDocs),
		MultipleOf:           s.MultipleOf,
		Minimum:              s.Minimum,
		Maximum:              s.Maximum,
		ExclusiveMinimum:     s.ExclusiveMinimum,
		ExclusiveMaximum:     s.ExclusiveMaximum,
		MinLength:            s.MinLength,
		MaxLength:            s.MaxLength,
		Pattern:              s.Pattern,
		MinItems:             s.MinItems,
		MaxItems:             s.MaxItems,
		UniqueItems:          s.UniqueItems,
		Items:                t.transformSchemaOrRef(s.Items),
		MinProperties:        s.MinProperties,
		MaxProperties:        s.MaxProperties,
		Required:             s.Required,
		Properties:           t.transformSchemaProperties(s.Properties),
		AdditionalProperties: t.transformAdditionalProperties(s.AdditionalProperties),
		AllOf:                t.transformSchemaOrRefs(s.AllOf),
		AnyOf:                t.transformSchemaOrRefs(s.AnyOf),
		OneOf:                t.transformSchemaOrRefs(s.OneOf),
		Not:                  t.transformSchemaOrRef(s.Not),
		Discriminator:        t.TransformDiscriminator(s.Discriminator),
		Enum:                 s.Enum,
	}
}

func (t *Transformer) transformSchemaProperties(props map[string]*model.SchemaOrRef) map[string]*SchemaOrRef {
	if len(props) == 0 {
		return nil
	}
	result := make(map[string]*SchemaOrRef)
	for name, prop := range props {
		result[name] = t.transformSchemaOrRef(prop)
	}
	return result
}

func (t *Transformer) transformSchemaOrRefs(schemas []*model.SchemaOrRef) []*SchemaOrRef {
	if len(schemas) == 0 {
		return nil
	}
	result := make([]*SchemaOrRef, len(schemas))
	for i, s := range schemas {
		result[i] = t.transformSchemaOrRef(s)
	}
	return result
}

func (t *Transformer) transformAdditionalProperties(ap *model.AdditionalProperties) *AdditionalProperties {
	if ap == nil {
		return nil
	}
	return &AdditionalProperties{
		Allowed: ap.Allowed,
		Schema:  t.transformSchemaOrRef(ap.Schema),
	}
}

func (t *Transformer) transformExamples(examples map[string]*model.Example) map[string]*Example {
	if len(examples) == 0 {
		return nil
	}
	result := make(map[string]*Example)
	for name, ex := range examples {
		if ex == nil {
			continue
		}
		result[name] = &Example{
			Ref:           ex.Ref,
			Summary:       ex.Summary,
			Description:   ex.Description,
			Value:         ex.Value,
			ExternalValue: ex.ExternalValue,
		}
	}
	return result
}

func (t *Transformer) transformComponents(c *model.Components) *Components {
	if c == nil {
		return nil
	}
	return &Components{
		Schemas:         t.transformComponentSchemas(c.Schemas),
		Responses:       t.transformComponentResponses(c.Responses),
		Parameters:      t.transformComponentParameters(c.Parameters),
		Examples:        t.transformComponentExamples(c.Examples),
		RequestBodies:   t.transformComponentRequestBodies(c.RequestBodies),
		Headers:         t.transformComponentHeaders(c.Headers),
		SecuritySchemes: t.transformComponentSecuritySchemes(c.SecuritySchemes),
		Links:           t.transformComponentLinks(c.Links),
		Callbacks:       t.transformComponentCallbacks(c.Callbacks),
	}
}

func (t *Transformer) transformComponentSchemas(schemas map[string]*model.SchemaOrRef) map[string]*SchemaOrRef {
	if len(schemas) == 0 {
		return nil
	}
	result := make(map[string]*SchemaOrRef)
	for name, s := range schemas {
		result[name] = t.transformSchemaOrRef(s)
	}
	return result
}

func (t *Transformer) transformComponentResponses(responses map[string]*model.ResponseOrRef) map[string]*ResponseOrRef {
	if len(responses) == 0 {
		return nil
	}
	result := make(map[string]*ResponseOrRef)
	for name, r := range responses {
		result[name] = t.transformResponseOrRef(r)
	}
	return result
}

func (t *Transformer) transformComponentParameters(params map[string]*model.ParameterOrRef) map[string]*ParameterOrRef {
	if len(params) == 0 {
		return nil
	}
	result := make(map[string]*ParameterOrRef)
	for name, p := range params {
		result[name] = t.transformParameterOrRef(p)
	}
	return result
}

func (t *Transformer) transformComponentExamples(examples map[string]*model.Example) map[string]*Example {
	if len(examples) == 0 {
		return nil
	}
	result := make(map[string]*Example)
	for name, ex := range examples {
		if ex == nil {
			continue
		}
		result[name] = &Example{
			Ref:           ex.Ref,
			Summary:       ex.Summary,
			Description:   ex.Description,
			Value:         ex.Value,
			ExternalValue: ex.ExternalValue,
		}
	}
	return result
}

func (t *Transformer) transformComponentRequestBodies(bodies map[string]*model.RequestBodyOrRef) map[string]*RequestBodyOrRef {
	if len(bodies) == 0 {
		return nil
	}
	result := make(map[string]*RequestBodyOrRef)
	for name, rb := range bodies {
		result[name] = t.transformRequestBodyOrRef(rb)
	}
	return result
}

func (t *Transformer) transformComponentHeaders(headers map[string]*model.HeaderOrRef) map[string]*HeaderOrRef {
	if len(headers) == 0 {
		return nil
	}
	result := make(map[string]*HeaderOrRef)
	for name, h := range headers {
		result[name] = t.transformHeaderOrRef(h)
	}
	return result
}

func (t *Transformer) transformComponentSecuritySchemes(schemes map[string]*model.SecuritySchemeOrRef) map[string]*SecuritySchemeOrRef {
	if len(schemes) == 0 {
		return nil
	}
	result := make(map[string]*SecuritySchemeOrRef)
	for name, ss := range schemes {
		result[name] = t.transformSecuritySchemeOrRef(ss)
	}
	return result
}

func (t *Transformer) transformSecuritySchemeOrRef(ss *model.SecuritySchemeOrRef) *SecuritySchemeOrRef {
	if ss == nil {
		return nil
	}
	if ss.Ref != "" {
		return &SecuritySchemeOrRef{Ref: ss.Ref}
	}
	if ss.Value == nil {
		return nil
	}
	return &SecuritySchemeOrRef{
		Value: t.TransformSecurityScheme(ss.Value),
	}
}

func (t *Transformer) transformComponentLinks(links map[string]*model.LinkOrRef) map[string]*LinkOrRef {
	if len(links) == 0 {
		return nil
	}
	result := make(map[string]*LinkOrRef)
	for name, l := range links {
		result[name] = t.transformLinkOrRef(l)
	}
	return result
}

func (t *Transformer) transformComponentCallbacks(callbacks map[string]*model.CallbackOrRef) map[string]*CallbackOrRef {
	if len(callbacks) == 0 {
		return nil
	}
	result := make(map[string]*CallbackOrRef)
	for name, cb := range callbacks {
		result[name] = t.transformCallbackOrRef(cb)
	}
	return result
}
