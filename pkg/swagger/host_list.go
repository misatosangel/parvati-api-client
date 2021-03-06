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

type HostList struct {
	Hosts []HosterStatus `json:"hosts,omitempty"`

	Waits []WaiterStatus `json:"waits,omitempty"`
}
