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

import (
	"time"
)

type StatusCheckResult struct {
	// Check id
	Id uint64 `json:"id,omitempty"`

	// Current Status of this host.
	Status string `json:"status,omitempty"`

	// Is spectateable?
	CanSpec string `json:"can_spec"`

	// Confused version string; needs to be fixed
	Version string `json:"roll"`

	// Last time this was verified.
	LastCheck time.Time `json:"last_check,omitempty"`

	// First time of this status result.
	CheckDate time.Time `json:"check_date,omitempty"`

	// Profile name of player one, if known.
	P1Profile string `json:"p1Profile,omitempty"`

	// Profile name of player two, if known.
	P2Profile string `json:"p2Profile,omitempty"`
}
