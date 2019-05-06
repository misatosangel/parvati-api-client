package cmd_parvati

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/misatosangel/parvati-api-client/commands/lowlevel"
	"github.com/misatosangel/parvati-api-client/iface"
	"github.com/misatosangel/parvati-api-client/parvatigo"
	"github.com/misatosangel/parvati-api-types-golang"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

type HostWatch struct {
	api           *parvatigo.Api
	apiConfig     *parvatigo.ApiConfig
	configFile    string
	EnabledGames  []string `short:"E" long:"enable" description:"Enable a game by (game) name or config section name." value-name:"<game>"`
	DisabledGames []string `short:"D" long:"disable" description:"Disable a game by (game) name or config section name. This wins over --enable." value-name:"<game>"`
	V4Iface       string   `long:"iface4" required:"false" value-name:"<name>|<id>" description:"Use this interface name/number for public v4 IP."`
	V6Iface       string   `long:"iface6" required:"false" value-name:"<name>|<id>" description:"Use this interface name/number for public v6 IP."`
	HostMessage   string   `short:"m" long:"host-message" decription:"Use this message to host (overrides config files)" value-name:"<text>"`
	NoIPUpdate    bool     `long:"no-ip-update" required:"false" description:"Do not also update IPs."`
}

func (self *HostWatch) AddCommands(base *flags.Command) (*flags.Command, error) {
	c, err := base.AddCommand("HostWatch", "Watch for you hosting.", "Use this to auto-check for you hosting and advertise.", self)
	if err != nil {
		return nil, err
	}
	c.Aliases = append(c.Aliases, "watch")
	return c, err
}

func (self *HostWatch) NeedsAPI() bool {
	return true
}

func (self *HostWatch) NeedsAPIConfig() bool {
	return true
}

func (self *HostWatch) SetAPI(api *parvatigo.Api) {
	self.api = api
}

func (self *HostWatch) SetAPIConfig(api *parvatigo.ApiConfig) {
	self.apiConfig = api
}

func (self *HostWatch) SetConfigFile(filePath string) {
	self.configFile = filePath
}

func (self *HostWatch) Execute(args []string) error {
	// check available configured games
	knownGames, err := self.api.GetGames()
	if err != nil {
		return err
	}
	if len(knownGames) == 0 {
		return fmt.Errorf("Parvati's backend is not configured; no known games were found.\n")
	}
	enabledGames := self.apiConfig.GetEnabledGames(self.EnabledGames, self.DisabledGames)
	if len(enabledGames) == 0 {
		return fmt.Errorf("Your configuration file and/or options does not enable any games.\n")
	}
	games, _ := cmd_lowlevel.FilterGames(knownGames, enabledGames, true)
	if len(games) == 0 {
		return fmt.Errorf("You have filtered out all known games.\n")
	}
	ipFlags := cmd_lowlevel.GetIpFlagsFromGames(games)
	if ipFlags == 0 {
		return fmt.Errorf("Your filtered games have no IP information.\n")
	}

	ifaceConfig, err := ConfigureIfacePrefs(self.configFile, self.V4Iface, self.V6Iface)
	if err != nil {
		return err
	}
	return self.noCuiMode(games, ifaceConfig, ipFlags)
}

func (self *HostWatch) noCuiMode(games []*cmd_lowlevel.GameConfig, ifaceConfig *iface.Config, ipFlags int) error {
	// now ready to do it
	fmt.Printf("Running update continually at 10s intervals. Hit CTRL+C to stop.\n")
	signalC := make(chan os.Signal, 1)
	signal.Notify(signalC, os.Interrupt)
	updateTicker := time.NewTicker(10 * time.Second)
	defer updateTicker.Stop()
	statMap := make(map[string]string)
	first := true

	for {
		select {
		case <-updateTicker.C:
			delta, err := UpdateIPs(self.api, ipFlags, ifaceConfig.V4ID, ifaceConfig.V6ID, !self.NoIPUpdate)
			if err != nil {
				log.Println(err)
				fmt.Println("Aborting host check on this iteration\n")
				continue
			}
			ProcessIPDelta(delta, self.NoIPUpdate, first)
			first = false
			for _, game := range games {
				name := game.ConfigInfo.PrettyName()
				status, err := CheckAutoHost(self.api, game, statMap[name], &delta.Player, self.HostMessage)
				if err != nil {
					log.Println(err)
					continue
				}
				stat := status.Status
				if statMap[name] != stat && stat == "Playing" {
					log.Println("You have been joined by opponent " + status.Opponent + "\n")
					joinLen := len(game.ConfigInfo.OnJoined)
					if joinLen > 0 {
						go func() {
							args := make([]string, joinLen-1, joinLen-1)
							if joinLen > 0 {
								for i, arg := range game.ConfigInfo.OnJoined[1:] {
									args[i] = strings.Replace(arg, "${NICK}", status.Opponent, -1)
								}
							}
							cmd := exec.Command(game.ConfigInfo.OnJoined[0], args...)
							cmd.Stdin = nil
							cmd.Stdout = os.Stdout
							cmd.Stderr = os.Stderr
							if runErr := cmd.Run(); runErr != nil {
								log.Println(runErr)
							}
						}()
					}
				}
				statMap[name] = stat
			}
		case sig := <-signalC:
			fmt.Println("Stopping on signal:", sig)
			return nil
		}
	}

	return nil
}

func CheckAutoHost(api *parvatigo.Api, gameConfig *cmd_lowlevel.GameConfig, lastStat string, user *swagger.User, hostMessage string) (*swagger.GameCheckInfo, error) {
	game := gameConfig.BackendGame
	hoster, waiter, err := api.UserInHostlist(game, user)
	if err != nil {
		return nil, fmt.Errorf("Unable to check existing hostlist: %s\n", err.Error())
	}
	if hoster != nil {
		info := api.HostAsCheckInfo(hoster)
		return info, nil // already listed
	}
	result, err := api.CheckHosting(game, user, "basic", uint(gameConfig.ConfigInfo.Port))
	if err != nil {
		return nil, err
	}
	if result.HostPort == "" {
		return nil, fmt.Errorf("%s host checking failed: %s\n", game.Name, result.Error)
	}
	switch result.Info.Status {
	case "Waiting", "Playing", "Relay":
		// post the host!
		ipStr, portStr, err := net.SplitHostPort(result.HostPort)
		if err != nil {
			return &result.Info, fmt.Errorf("Failed to parse ip:port result '%s': %s\n", result.HostPort, err.Error())
		}
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return &result.Info, fmt.Errorf("Failed to parse ip address in '%s': %s\n", result.HostPort, ipStr)
		}
		port, err := strconv.ParseUint(portStr, 10, 32)
		if err != nil {
			return &result.Info, fmt.Errorf("Failed to parse port of '%s' as numeric '%s': %s\n", result.HostPort, portStr, err.Error())
		}
		mes := hostMessage
		if mes == "" {
			mes = gameConfig.ConfigInfo.HostMessage()
		}
		_, err = api.PostUserHost(game, user, ip, uint(port), mes)
		if err != nil {
			return &result.Info, fmt.Errorf("%s host announce on %s failed: %s\n", game.Name, result.HostPort, err.Error())
		}
		fmt.Printf("%s host announce succeeded.\n", game.Name)
		return &result.Info, nil
	default:
		if result.Info.Status != lastStat {
			fmt.Printf("%s host check on %s gave result %s\n", game.Name, result.HostPort, result.Info.Status)
		}
		if waiter != nil {
			return nil, nil
		}
		return &result.Info, nil
	}
}
