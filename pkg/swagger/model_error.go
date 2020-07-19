/*
 * Host API
 *
 * API for posting hosts and waiting for hosts
 *
 * OpenAPI spec version: 1.0.0
 *
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 *
 * Licensed under the MIT License
 */

package swagger

type ModelError struct {
	Code int32 `json:"code,omitempty"`

	Message string `json:"message,omitempty"`

	Fields string `json:"fields,omitempty"`
}
