package cmd_parvati

import (
	"github.com/jessevdk/go-flags"
	"github.com/misatosangel/parvati-api-client/parvatigo"
)

type Details struct {
	api            *parvatigo.Api
	User           string `long:"user" short:"u" required:"false" description:"Get details for this user rather than yourself."`
	ShowIdentities bool   `long:"show-identities" short:"i" required:"false" description:"Show identities. Only valid if getting your own data."`
}

func (self *Details) AddCommands(base *flags.Command) (*flags.Command, error) {
	c, err := base.AddCommand("Details", "Show Parvati's user information.", "Use this command to Parvati's information on yourself or someone else.", self)
	if err != nil {
		return nil, err
	}
	c.Aliases = append(c.Aliases, "user-info")
	return c, err
}

func (self *Details) NeedsAPI() bool {
	return true
}

func (self *Details) NeedsAPIConfig() bool {
	return false
}

func (self *Details) SetAPI(api *parvatigo.Api) {
	self.api = api
}

func (self *Details) SetAPIConfig(api *parvatigo.ApiConfig) {
}

func (self *Details) Execute(args []string) error {
	me, err := self.api.GetDetails()
	if err != nil {
		return err
	}
	if self.User == "" {
		DumpUserData(me, self.ShowIdentities, true)
		return nil
	}
	data, err := self.api.GetUserDetails(self.User)
	if err != nil {
		return err
	}
	if data.Id == me.Id {
		DumpUserData(me, self.ShowIdentities, true)
		return nil
	}
	admin := me.PrivLevel == "super"
	showIdent := admin && self.ShowIdentities
	DumpUserData(data, showIdent, admin)
	return nil
}
