/*
 * examples/internal/proto/examplepb/response_body_service.proto
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: version not set
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package responsebody

type StreamResultOfExamplepbResponseBodyOut struct {
	Result *ExamplepbResponseBodyOutResponse `json:"result,omitempty"`
	Error_ *RpcStatus                        `json:"error,omitempty"`
}
