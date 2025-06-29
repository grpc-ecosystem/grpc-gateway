package genopenapiv3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	"google.golang.org/protobuf/proto"
)

func (g *generator) convertFileOptions(target *descriptor.File) (*openapi3.T, bool) {
	if openAPIAno := proto.GetExtension(target.GetOptions(), options.E_Openapiv3Document).(*options.OpenAPI); openAPIAno != nil {
		doc := &openapi3.T{
			OpenAPI:  OpenAPIVersion,
			Info:     g.convertInfo(openAPIAno.GetInfo()),
			Security: *convertSecurityRequiremnt(openAPIAno.GetSecurity()),
			Servers:  g.convertServers(openAPIAno.GetServers()),
		}
		// TODO: implement other openapi file annotation fields
		return doc, true
	}

	return nil, false
}

func (g *generator) convertServers(servers []*options.Server) openapi3.Servers {
	oAPIservers := make(openapi3.Servers, len(servers))

	for i, srv := range servers {
		vars := map[string]*openapi3.ServerVariable{}

		for k, v := range srv.GetVariables() {
			vars[k] = &openapi3.ServerVariable{
				Enum:        v.GetEnum(),
				Default:     v.GetDefault(),
				Description: v.GetDescription(),
			}
		}

		oAPIservers[i] = &openapi3.Server{
			URL:         srv.GetUrl(),
			Description: srv.GetDescription(),
			Variables:   vars,
		}
	}

	return oAPIservers
}

func (g *generator) convertInfo(openAPIInfo *options.Info) *openapi3.Info {
	return &openapi3.Info{
		Title:          openAPIInfo.GetTitle(),
		Description:    openAPIInfo.GetDescription(),
		Version:        openAPIInfo.GetVersion(),
		TermsOfService: openAPIInfo.GetTermsOfService(),
		Contact:        g.convertContact(openAPIInfo.GetContact()),
		License:        g.convertLicense(openAPIInfo.GetLicense()),
	}
}

func (g *generator) convertContact(contactOption *options.Contact) *openapi3.Contact {
	if contactOption == nil {
		return nil
	}

	return &openapi3.Contact{
		Name:  contactOption.GetName(),
		URL:   contactOption.GetUrl(),
		Email: contactOption.GetEmail(),
	}
}

func (g *generator) convertLicense(licenseOption *options.License) *openapi3.License {
	if licenseOption == nil {
		return nil
	}

	return &openapi3.License{
		Name: licenseOption.GetName(),
		URL:  licenseOption.GetUrl(),
	}
}
