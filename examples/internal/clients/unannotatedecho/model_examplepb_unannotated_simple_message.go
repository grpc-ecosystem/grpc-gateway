/*
 * Unannotated Echo
 *
 * Unannotated Echo Service Similar to echo_service.proto but without annotations. See unannotated_echo_service.yaml for the equivalent of the annotations in gRPC API configuration format.  Echo Service API consists of a single service which returns a message.
 *
 * API version: 1.0
 * Contact: none@example.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package unannotatedecho

// A simple message with many types
type ExamplepbUnannotatedSimpleMessage struct {
	// Id represents the message identifier.
	Id string `json:"id"`
	// Int value field
	Num string `json:"num"`
	Duration string `json:"duration,omitempty"`
	LineNum string `json:"lineNum,omitempty"`
	Lang string `json:"lang,omitempty"`
	Status *ExamplepbUnannotatedEmbedded `json:"status,omitempty"`
	En string `json:"en,omitempty"`
	No *ExamplepbUnannotatedEmbedded `json:"no,omitempty"`
	ResourceId string `json:"resourceId,omitempty"`
	NId *ExamplepbUnannotatedNestedMessage `json:"nId,omitempty"`
}
