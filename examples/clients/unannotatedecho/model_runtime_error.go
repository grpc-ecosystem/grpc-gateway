/*
 * examples/proto/examplepb/unannotated_echo_service.proto
 *
 * Unannotated Echo Service Similar to echo_service.proto but without annotations. See unannotated_echo_service.yaml for the equivalent of the annotations in gRPC API configuration format.  Echo Service API consists of a single service which returns a message.
 *
 * API version: version not set
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package unannotatedecho

type RuntimeError struct {
	Error_ string `json:"error,omitempty"`
	Code int32 `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Details []ProtobufAny `json:"details,omitempty"`
}
