package transform

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/internal/genopenapiv3/model"
)

// BaseTransformer provides common transformation logic for OpenAPI adapters.
// Version-specific transformers embed this to inherit shared transformations.
type BaseTransformer struct{}

// TransformInfo converts canonical Info to shared Info type.
func (b *BaseTransformer) TransformInfo(info *model.Info) *Info {
	if info == nil {
		return nil
	}
	return &Info{
		Title:          info.Title,
		Summary:        info.Summary,
		Description:    info.Description,
		TermsOfService: info.TermsOfService,
		Contact:        b.TransformContact(info.Contact),
		License:        b.TransformLicense(info.License),
		Version:        info.Version,
	}
}

// TransformContact converts canonical Contact to shared Contact type.
func (b *BaseTransformer) TransformContact(c *model.Contact) *Contact {
	if c == nil {
		return nil
	}
	return &Contact{
		Name:  c.Name,
		URL:   c.URL,
		Email: c.Email,
	}
}

// TransformLicense converts canonical License to shared License type.
func (b *BaseTransformer) TransformLicense(l *model.License) *License {
	if l == nil {
		return nil
	}
	return &License{
		Name:       l.Name,
		Identifier: l.Identifier,
		URL:        l.URL,
	}
}

// TransformServers converts a slice of canonical Servers.
func (b *BaseTransformer) TransformServers(servers []*model.Server) []*Server {
	if len(servers) == 0 {
		return nil
	}
	result := make([]*Server, len(servers))
	for i, s := range servers {
		result[i] = b.TransformServer(s)
	}
	return result
}

// TransformServer converts canonical Server to shared Server type.
func (b *BaseTransformer) TransformServer(s *model.Server) *Server {
	if s == nil {
		return nil
	}
	var vars map[string]*ServerVariable
	if len(s.Variables) > 0 {
		vars = make(map[string]*ServerVariable)
		for name, v := range s.Variables {
			vars[name] = &ServerVariable{
				Enum:        v.Enum,
				Default:     v.Default,
				Description: v.Description,
			}
		}
	}
	return &Server{
		URL:         s.URL,
		Description: s.Description,
		Variables:   vars,
	}
}

// TransformExternalDocs converts canonical ExternalDocs to shared ExternalDocs type.
func (b *BaseTransformer) TransformExternalDocs(ed *model.ExternalDocs) *ExternalDocs {
	if ed == nil {
		return nil
	}
	return &ExternalDocs{
		Description: ed.Description,
		URL:         ed.URL,
	}
}

// TransformTags converts a slice of canonical Tags.
func (b *BaseTransformer) TransformTags(tags []*model.Tag) []*Tag {
	if len(tags) == 0 {
		return nil
	}
	result := make([]*Tag, len(tags))
	for i, t := range tags {
		result[i] = &Tag{
			Name:         t.Name,
			Description:  t.Description,
			ExternalDocs: b.TransformExternalDocs(t.ExternalDocs),
		}
	}
	return result
}

// TransformSecurityScheme converts canonical SecurityScheme to shared SecurityScheme type.
func (b *BaseTransformer) TransformSecurityScheme(ss *model.SecurityScheme) *SecurityScheme {
	if ss == nil {
		return nil
	}
	return &SecurityScheme{
		Type:             ss.Type,
		Description:      ss.Description,
		Name:             ss.Name,
		In:               ss.In,
		Scheme:           ss.Scheme,
		BearerFormat:     ss.BearerFormat,
		Flows:            b.TransformOAuthFlows(ss.Flows),
		OpenIDConnectURL: ss.OpenIDConnectURL,
	}
}

// TransformOAuthFlows converts canonical OAuthFlows to shared OAuthFlows type.
func (b *BaseTransformer) TransformOAuthFlows(flows *model.OAuthFlows) *OAuthFlows {
	if flows == nil {
		return nil
	}
	return &OAuthFlows{
		Implicit:          b.TransformOAuthFlow(flows.Implicit),
		Password:          b.TransformOAuthFlow(flows.Password),
		ClientCredentials: b.TransformOAuthFlow(flows.ClientCredentials),
		AuthorizationCode: b.TransformOAuthFlow(flows.AuthorizationCode),
	}
}

// TransformOAuthFlow converts canonical OAuthFlow to shared OAuthFlow type.
func (b *BaseTransformer) TransformOAuthFlow(flow *model.OAuthFlow) *OAuthFlow {
	if flow == nil {
		return nil
	}
	return &OAuthFlow{
		AuthorizationURL: flow.AuthorizationURL,
		TokenURL:         flow.TokenURL,
		RefreshURL:       flow.RefreshURL,
		Scopes:           flow.Scopes,
	}
}

// TransformSecurityRequirements converts canonical SecurityRequirements.
func (b *BaseTransformer) TransformSecurityRequirements(reqs []model.SecurityRequirement) []SecurityRequirement {
	if len(reqs) == 0 {
		return nil
	}
	result := make([]SecurityRequirement, len(reqs))
	for i, req := range reqs {
		result[i] = SecurityRequirement(req)
	}
	return result
}

// TransformDiscriminator converts canonical Discriminator to shared Discriminator type.
func (b *BaseTransformer) TransformDiscriminator(d *model.Discriminator) *Discriminator {
	if d == nil {
		return nil
	}
	return &Discriminator{
		PropertyName: d.PropertyName,
		Mapping:      d.Mapping,
	}
}

// TransformLink converts canonical Link to shared Link type.
func (b *BaseTransformer) TransformLink(l *model.Link) *Link {
	if l == nil {
		return nil
	}
	return &Link{
		OperationRef: l.OperationRef,
		OperationID:  l.OperationID,
		Parameters:   l.Parameters,
		RequestBody:  l.RequestBody,
		Description:  l.Description,
		Server:       b.TransformServer(l.Server),
	}
}

// TransformExample converts canonical Example to shared Example type.
// Note: Ref handling is done in version-specific code since 3.0.x and 3.1.0
// may handle $ref differently in examples.
func (b *BaseTransformer) TransformExample(ex *model.Example) *Example {
	if ex == nil {
		return nil
	}
	return &Example{
		Summary:       ex.Summary,
		Description:   ex.Description,
		Value:         ex.Value,
		ExternalValue: ex.ExternalValue,
	}
}

// TransformExamples converts a map of canonical Examples.
func (b *BaseTransformer) TransformExamples(examples map[string]*model.Example) map[string]*Example {
	if len(examples) == 0 {
		return nil
	}
	result := make(map[string]*Example)
	for name, ex := range examples {
		if ex == nil {
			continue
		}
		result[name] = b.TransformExample(ex)
	}
	return result
}
