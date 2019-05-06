package cmd_parvati

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/misatosangel/parvati-api-client/commands/lowlevel"
	"github.com/misatosangel/parvati-api-client/parvatigo"
	"github.com/misatosangel/traceroute"
	"os"
	"os/signal"
	"time"
)

type UpdateIP struct {
	api        *parvatigo.Api
	apiConfig  *parvatigo.ApiConfig
	configFile string
	SetV6      bool   `short:"6" required:"false" description:"Update v6 IP (ignores enabled games)."`
	SetV4      bool   `short:"4" required:"false" description:"Update v4 IP (ignores enabled games)."`
	Check      bool   `short:"n" long:"no-update" required:"false" description:"Just show what would be done, do not actually update."`
	V4Iface    string `long:"iface4" required:"false" value-name:"<name>|<id>" description:"Use this interface name/number for public v4 IP."`
	V6Iface    string `long:"iface6" required:"false" value-name:"<name>|<id>" description:"Use this interface name/number for public v6 IP."`
	Repeat     bool   `long:"repeat" short:"r" required:"false" description:"Constantly updated over time."`
}

func (self *UpdateIP) AddCommands(base *flags.Command) (*flags.Command, error) {
	c, err := base.AddCommand("UpdateIP", "Update Parvati's stored IP for you.", "Use this to set/update Parvati's IP based on your current public IP address.", self)
	if err != nil {
		return nil, err
	}
	c.Aliases = append(c.Aliases, "set-ip")
	return c, err
}

func (self *UpdateIP) NeedsAPI() bool {
	return true
}

func (self *UpdateIP) NeedsAPIConfig() bool {
	return true
}

func (self *UpdateIP) SetAPI(api *parvatigo.Api) {
	self.api = api
}

func (self *UpdateIP) SetAPIConfig(api *parvatigo.ApiConfig) {
	self.apiConfig = api
}

func (self *UpdateIP) SetConfigFile(filePath string) {
	self.configFile = filePath
}

func (self *UpdateIP) Execute(args []string) error {
	var ipFlags int
	if self.SetV6 {
		ipFlags = traceroute.WANT_PUBLIC_V6
		if self.SetV4 {
			ipFlags |= traceroute.WANT_PUBLIC_V4
		}
	} else if self.SetV4 {
		ipFlags = traceroute.WANT_PUBLIC_V4
	} else {
		// check available configured games
		knownGames, err := self.api.GetGames()
		if err != nil {
			return err
		}
		msgTry := "You can still force set v4/v6 with -4 and/or -6.\n"
		if len(knownGames) == 0 {
			return fmt.Errorf("Parvati's backend is not configured; no known games were found.\n" + msgTry)
		}
		enabledGames := self.apiConfig.GetEnabledGames(nil, nil)
		if len(enabledGames) == 0 {
			return fmt.Errorf("Your configuration file does not enable any games which specify interfaces to check.\n" + msgTry)
		}
		games, _ := cmd_lowlevel.FilterGames(knownGames, enabledGames, true)
		if len(games) == 0 {
			return fmt.Errorf("You have filtered out all known games.\n" + msgTry)
		}
		ipFlags = cmd_lowlevel.GetIpFlagsFromGames(games)
		if ipFlags == 0 {
			return fmt.Errorf("Your filtered games have no IP information.\n" + msgTry)
		}
	}
	ifaceConfig, err := ConfigureIfacePrefs(self.configFile, self.V4Iface, self.V6Iface)
	if err != nil {
		return err
	}
	// no ready to do it
	delta, err := UpdateIPs(self.api, ipFlags, ifaceConfig.V4ID, ifaceConfig.V6ID, !self.Check)
	if err != nil {
		return err
	}
	if !self.Repeat {
		return ProcessIPDelta(delta, self.Check, true)
	}
	fmt.Printf("Running update continually at 15s intervals. Hit CTRL+C to stop.\n")
	err = ProcessIPDelta(delta, self.Check, true)
	if err != nil {
		fmt.Println(err)
	}
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, os.Interrupt)
	updateTicker := time.NewTicker(15 * time.Second)
	defer updateTicker.Stop()
	for {
		select {
		case <-updateTicker.C:
			delta, err := UpdateIPs(self.api, ipFlags, ifaceConfig.V4ID, ifaceConfig.V6ID, !self.Check)
			if err == nil {
				err = ProcessIPDelta(delta, self.Check, true)
			}
			if err != nil {
				fmt.Println(err)
			}
		case sig := <-signalC:
			fmt.Println("Stopping on signal:", sig)
			return nil
		}
	}

	return nil
}
