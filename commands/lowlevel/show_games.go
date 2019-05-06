package cmd_lowlevel

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/misatosangel/parvati-api-client/parvatigo"
	"github.com/misatosangel/parvati-api-types-golang"
	"sort"
	"strings"
)

type ShowGames struct {
	api       *parvatigo.Api
	APIs      bool `short:"a" long:"apis" description:"Print known checking API endpoints"`
	Info      bool `short:"i" long:"info" description:"Print additiona name/url information"`
	Protocols bool `short:"p" long:"protocols" description:"Print Supported protocol information"`
}

func (self *ShowGames) AddCommands(base *flags.Command) (*flags.Command, error) {
	c, err := base.AddCommand("ListKnownGames", "Show Parvati's game information", "Use this command to show all Parvati's known games", self)
	if err != nil {
		return nil, err
	}
	c.Aliases = append(c.Aliases, "ls-games")
	return c, err
}

func (self *ShowGames) NeedsAPI() bool {
	return true
}

func (self *ShowGames) NeedsAPIConfig() bool {
	return false
}

func (self *ShowGames) SetAPI(api *parvatigo.Api) {
	self.api = api
}

func (self *ShowGames) SetAPIConfig(api *parvatigo.ApiConfig) {
}

func (self *ShowGames) Execute(args []string) error {
	knownGames, err := self.api.GetGames()
	if err != nil {
		return err
	}
	for _, g := range knownGames {
		MakupGame(g, self.APIs, self.Protocols, self.Info)
	}
	return nil
}

func MakupGame(g swagger.Game, apis, protocols, info bool) {
	fmt.Printf("%d. %s\n", g.Id, g.Name)
	if info {
		fmt.Printf("Aka: %s\n"+
			"Discord name: %s\n"+
			"Info Url: %s\n", g.UrlShortName, g.DiscordName, g.Url)
	}
	if g.Port != 0 {
		fmt.Printf("Default Port: %d\n", g.Port)
	}
	if protocols {
		proto := make([]string, 0, 2)
		protocols := make(map[string]bool)
		ipv4Sup := false
		ipv6Sup := false
		for _, pv := range g.Protocols {
			pLen := len(pv)
			p := pv[0 : pLen-1]
			switch pv[pLen-1] {
			case '4':
				ipv4Sup = true
			case '6':
				ipv6Sup = true
			}
			if !protocols[p] {
				proto = append(proto, p)
				protocols[p] = true
			}
		}
		if len(proto) == 0 {
			fmt.Printf("Protocols: unknown\n")
		} else {
			sort.Strings(proto)
			fmt.Printf("Protocols: %s\n", strings.Join(proto, ", "))
			fmt.Printf("IPv4 Support: %t\n", ipv4Sup)
			fmt.Printf("IPv6 Support: %t\n", ipv6Sup)
		}
	}
	if apis {
		apiCnt := len(g.APIs)
		if apiCnt == 0 {
			fmt.Printf("No Parvati checking API information\n")
			return
		}
		fmt.Printf("Parvati Checking URIs:\n")
		for _, u := range g.APIs {
			fmt.Printf(" - %s\n", u)
		}
	}
}
