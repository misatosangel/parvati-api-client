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

type ChallongeDetails struct {
	ID       uint64 `json:"id,omitempty"`
	PID      uint64 `json:"player_id,omitempty"`
	ApiKey   string `json:"api_key,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

type ChallongeDelta struct {
	Added    *ChallongeDetails   `json:"added,omitempty"`
	Removed  *ChallongeDetails   `json:"removed,omitempty"`
	Modified []*ChallongeDetails `json:"modified,omitempty"`
}

type Delta struct {
	IPPort     []string        `json:"ip,omitempty"`
	IPv4       []string        `json:"ipv4,omitempty"`
	IPv6       []string        `json:"ipv6,omitempty"`
	Port       []uint16        `json:"port,omitempty"`
	Challonge  *ChallongeDelta `json:"challonge,omitempty"`
	StaticIP   []uint8         `json:"ip_lock,omitempty"`
	Password   []string        `json:"password,omitempty"`
	JoinNotify []uint8         `json:"join_notify,omitempty"`
	Gender     []string        `json:"preferred_gender,omitempty"`
	Avatar     []string        `json:"avatar,omitempty"`
	Nick       []string        `json:"nick,omitempty"`
	PVsMsg     []string        `json:"playing_vs_message,omitempty"`
	HostMsg    []string        `json:"hosting_message,omitempty"`
}
