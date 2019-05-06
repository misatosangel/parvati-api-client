package cmd_lowlevel

import (
	"github.com/jessevdk/go-flags"
	"github.com/misatosangel/parvati-api-client/parvatigo"

	"fmt"
	"github.com/misatosangel/traceroute"
	"log"
)

type IfaceList struct {
	api          *parvatigo.Api
	apiConfig    *parvatigo.ApiConfig
	IgnoreConfig bool `short:"i" long:"ignore-config" required:"false" description:"Ignore game configuration pointers for filtering IP families."`
	ShowV6       bool `short:"6" required:"false" description:"Include v6 IPs. Implies --ignore-config."`
	ShowV4       bool `short:"4" required:"false" description:"Include v4 IPs. Implies --ignore-config."`
}

func (self *IfaceList) AddCommands(base *flags.Command) (*flags.Command, error) {
	c, err := base.AddCommand("InterfaceList", "Show your active interface list.", "Use this command to show all active interfaces and NAT setup, possibly filtered by those used by your enabled games.", self)
	if err != nil {
		return nil, err
	}
	c.Aliases = append(c.Aliases, "ls-interfaces")
	c.Aliases = append(c.Aliases, "ls-ifaces")
	return c, err
}

func (self *IfaceList) NeedsAPI() bool {
	return false
}

func (self *IfaceList) NeedsAPIConfig() bool {
	return false
}

func (self *IfaceList) SetAPI(api *parvatigo.Api) {
	self.api = api
}

func (self *IfaceList) SetAPIConfig(apiConfig *parvatigo.ApiConfig) {
	self.apiConfig = apiConfig
}

func (self *IfaceList) Execute(args []string) error {
	if self.ShowV4 || self.ShowV6 {
		self.IgnoreConfig = true
	} else if self.IgnoreConfig { // ignore by itself implies show all
		self.ShowV4 = true
		self.ShowV6 = true
	}
	if self.api == nil || self.apiConfig == nil || self.IgnoreConfig {
		return ShowDefaultList(self.ShowV4, self.ShowV6, false)
	}
	knownGames, err := self.api.GetGames()
	if err != nil {
		return err
	}
	if len(knownGames) == 0 {
		return fmt.Errorf("Parvati's backend is not configured; no known games were found.\n")
	}
	enabledGames := self.apiConfig.GetEnabledGames(nil, nil)
	if len(enabledGames) == 0 {
		return fmt.Errorf("Your configuration file does not enable any games specify interfaces to check.\n")
	}
	games, _ := FilterGames(knownGames, enabledGames, true)
	if len(games) == 0 {
		log.Fatalf("You have filtered out all known games.\n")
	}
	ipFlags := GetIpFlagsFromGames(games)
	wantIPV4 := false
	if ipFlags&traceroute.WANT_PUBLIC_V4 != 0 {
		wantIPV4 = true
	}
	wantIPV6 := false
	if ipFlags&traceroute.WANT_PUBLIC_V6 != 0 {
		wantIPV6 = true
	}
	return ShowDefaultList(wantIPV4, wantIPV6, true)
}
