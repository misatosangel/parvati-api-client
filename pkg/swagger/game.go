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

type Game struct {

	// Unique identifier for a supported game
	Id int32 `json:"id,omitempty"`

	// Unique short-name identifier that can be used in urls.
	UrlShortName string `json:"url_short_name,omitempty"`

	// Name of the game.
	Name string `json:"name,omitempty"`

	// Name of the game according to Discord.
	DiscordName string `json:"discord_name,omitempty"`

	// Main URL for info on the game.
	Url string `json:"url,omitempty"`

	// Default port used by the game
	Port uint16 `json:"port,omitempty"`

	// API Entries
	APIs []APIEntry `json:"api_entries,omitempty"`

	// Protcols understood by the game
	Protocols []string `json:"protocols,omitempty"`
}

type APIEntry struct {
	Uri    string `json:"uri,omitempty"`
	GameId int32  `json:"game_id,omitempty"`
}

type GamePing struct {
	Request  string `json:"request,omitempty"`
	HostPort string `json:"hostport,omitempty"`
	Result   string `json:"result,omitempty"`
	TimeNS   uint64 `json:"timeNS,omitempty"`
}

type GameCheckResult struct {
	Request  string        `json:"request,omitempty"`
	HostPort string        `json:"hostport,omitempty"`
	Info     GameCheckInfo `json:"result,omitempty"`
	Error    string        `json:"error,omitempty"`
}

type GameCheckInfo struct {
	Address  string   `json:"address,omitempty"`
	Status   string   `json:"status,omitempty"`
	Version  string   `json:"version,omitempty"`
	Spectate int      `json:"spectate,omitempty"` // one of y(es) n(o) or u(nknown)
	Opponent string   `json:"opponent,omitempty"`
	Profiles []string `json:"profiles,omitempty"`
	Error    string   `json:"error,omitempty"`
}