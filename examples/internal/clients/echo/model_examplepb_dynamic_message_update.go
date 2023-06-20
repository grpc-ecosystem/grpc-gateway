/*
 * Echo Service
 *
 * Echo Service API consists of a single service which returns a message.
 *
 * API version: version not set
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package echo

type ExamplepbDynamicMessageUpdate struct {
	Body       *ExamplepbDynamicMessage `json:"body,omitempty"`
	UpdateMask string                   `json:"updateMask,omitempty"`
}
