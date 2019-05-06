package cmd_lowlevel

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/misatosangel/parvati-api-client/iface"
	"github.com/misatosangel/parvati-api-client/parvatigo"
	"github.com/misatosangel/parvati-api-types-golang"
	"github.com/misatosangel/traceroute"
	"os"
	"strings"
)

type GameConfig struct {
	BackendGame *swagger.Game
	ConfigInfo  *parvatigo.GameInfo
}

func AddCommands(base *flags.Command) (*flags.Command, error) {
	if _, err := (&ShowGames{}).AddCommands(base); err != nil {
		return base, err
	}
	if _, err := (&IfaceList{}).AddCommands(base); err != nil {
		return base, err
	}
	if _, err := (&ConfigHelp{}).AddCommands(base); err != nil {
		return base, err
	}
	return base, nil
}

func ShowDefaultList(wantV4, wantV6, filtered bool) error {
	flags := traceroute.WANT_LIVE_IP
	fStr := ""
	if filtered {
		fStr = " [auto-filtered]"
	}
	ipTypeStr := "(v4 only" + fStr + ")"
	if wantV6 {
		flags |= traceroute.WANT_PUBLIC_V6 | traceroute.WANT_PRIVATE_V6
		if wantV4 {
			flags |= traceroute.WANT_PUBLIC_V4 | traceroute.WANT_PRIVATE_V4
			ipTypeStr = "(v4 and v6" + fStr + ")"
		} else {
			ipTypeStr = "(v6 only" + fStr + ")"
		}
	} else {
		flags |= traceroute.WANT_PUBLIC_V4 | traceroute.WANT_PRIVATE_V4
	}
	list, err := iface.NewList(flags)
	if err != nil {
		return err
	}
	fmt.Printf("Default Interface list " + ipTypeStr + ":\n")
	list.Show(os.Stdout)
	return nil
}

func FilterGames(knownGames []swagger.Game, configGames []parvatigo.GameInfo, ignoreDups bool) ([]*GameConfig, error) {
	out := make([]*GameConfig, 0, 3)
	found := make(map[string]parvatigo.GameInfo, 3)
	for _, sGame := range knownGames {
		name := sGame.UrlShortName
		for _, config := range configGames {
			if !config.Enabled {
				continue
			}
			if !strings.EqualFold(config.Name, name) {
				continue
			}
			lc_vers := strings.ToLower(config.Name)
			matched, exists := found[lc_vers]
			if exists {
				if ignoreDups {
					continue
				}
				return nil, fmt.Errorf("Multiple active conflicting configurations found:\n" + matched.String() + "--vs-\n" + config.String())
			}
			found[lc_vers] = config
			cpyGame := sGame
			cpyConfig := config
			out = append(out, &GameConfig{BackendGame: &cpyGame, ConfigInfo: &cpyConfig})
		}
	}
	return out, nil
}

func GetIpFlagsFromGames(games []*GameConfig) int {
	flags := 0
	all := traceroute.WANT_PUBLIC_V4 | traceroute.WANT_PUBLIC_V6
	for _, game := range games {
		for _, proto := range game.BackendGame.Protocols {
			l := len(proto)
			if l == 0 {
				continue
			}
			if proto[l-1] == '4' {
				flags |= traceroute.WANT_PUBLIC_V4
			} else if proto[l-1] == '6' {
				flags |= traceroute.WANT_PUBLIC_V6
			} else {
				// unknown or busted
				continue
			}
			if flags == all {
				return all
			}
		}
	}
	return flags
}
