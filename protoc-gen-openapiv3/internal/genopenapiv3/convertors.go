package genopenapiv3

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv3/options"
)

func convertSecurityRequiremnt(requirements []*options.SecurityRequirement) *openapi3.SecurityRequirements {
	oAPIReqs := openapi3.NewSecurityRequirements()

	for _, req := range requirements {
		oAPISecReq := openapi3.NewSecurityRequirement()
		for authenticator, scopes := range req.GetAdditionalProperties() {
			oAPISecReq.Authenticate(authenticator, scopes.GetScopes()...)
		}

		oAPIReqs.With(oAPISecReq)
	}

	return oAPIReqs
}
