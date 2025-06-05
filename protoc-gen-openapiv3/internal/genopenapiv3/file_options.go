package genopenapiv3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
	"google.golang.org/protobuf/proto"
)

func (g *generator) extractFileOptions(target *descriptor.File) (*openapi3.T, bool) {
	if openAPIAno := proto.GetExtension(target.GetOptions(), options.E_Openapiv3Document).(*options.OpenAPI); openAPIAno != nil {
		doc := &openapi3.T{
			OpenAPI: "3.0.2",
		}
		doc.Info = g.extractInfo(openAPIAno.GetInfo())
		// TODO: implement other openapi file annotation fields
		return doc, true
	}

	return nil, false
}

func (g *generator) extractInfo(openAPIInfo *options.Info) *openapi3.Info {
	return &openapi3.Info{
		Title: openAPIInfo.GetTitle(),
		Description: openAPIInfo.GetDescription(),
		Version: openAPIInfo.GetVersion(),
		TermsOfService: openAPIInfo.GetTermsOfService(),
		Contact: g.extractContact(openAPIInfo.GetContact()),
		License: g.extractLicense(openAPIInfo.GetLicense()),
	}
}

func (g *generator) extractContact(contactOption *options.Contact) *openapi3.Contact {
	if contactOption == nil {
		return nil
	}

	contact := &openapi3.Contact{}
	contact.Name = contactOption.GetName()
	contact.URL = contactOption.GetUrl()
	contact.Email = contactOption.GetEmail()

	return contact
}

func (g *generator) extractLicense(licenseOption *options.License) *openapi3.License {
	if licenseOption == nil {
		return nil
	}

	license := &openapi3.License{}
	license.Name = licenseOption.GetName()
	license.URL = licenseOption.GetUrl()

	return license
}
