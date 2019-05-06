package cmd_generic

import (
	"github.com/misatosangel/parvati-api-client/parvatigo"
)

type APICommand interface {
	NeedsAPI() bool
	NeedsAPIConfig() bool
	SetAPI(api *parvatigo.Api)
	SetAPIConfig(api *parvatigo.ApiConfig)
}

type IfaceCommand interface {
	SetConfigFile(string)
}
