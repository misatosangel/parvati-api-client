package parvatigo

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty"
	"github.com/misatosangel/parvati-api-types-golang"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const defUri = "https://parvati.phi.al"

type Api struct {
	HApi      *swagger.HostsApi
	GApi      *swagger.GamesApi
	UApi      *swagger.UsersApi
	Config    *swagger.Configuration
	userID    string
	announcer string
	Verbose   bool
	log       *log.Logger
}

func NewApi(conf *ApiConfig, buildVersion string) (Api, error) {
	c := swagger.NewConfiguration()
	if buildVersion == "" {
		buildVersion = "dev"
	}
	c.UserAgent = "Parvati-Client/" + buildVersion + "/go"
	a := Api{Config: c}
	if conf == nil {
		var err error
		conf, err = ReadDefaultConfig()
		if err != nil {
			if !os.IsNotExist(err) {
				return a, err
			}
			conf = &ApiConfig{}
		}
	}
	if conf.URI == "" {
		c.BasePath = defUri
	} else {
		c.BasePath = conf.URI
	}
	if conf.Username != "" {
		c.UserName = conf.Username
	}
	if conf.Password != "" {
		c.Password = conf.Password
	}
	a.announcer = conf.Announcer
	if a.announcer == "" {
		a.announcer = "ApiClient"
	}
	a.HApi = &swagger.HostsApi{Configuration: *c}
	a.GApi = &swagger.GamesApi{Configuration: *c}
	a.UApi = &swagger.UsersApi{Configuration: *c}
	a.log = log.New(os.Stderr, "API> ", log.LstdFlags)
	return a, nil
}

func (self *Api) GetGames() ([]swagger.Game, error) {
	data, _, err := self.GApi.GamesGet("")
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (self *Api) UpdateIPs(v4, v6 net.IP) (swagger.UserDelta, error) {
	if v4 == nil && v6 == nil {
		return swagger.UserDelta{}, nil
	}
	if self.userID == "" {
		return swagger.UserDelta{}, fmt.Errorf("No user id for user: " + self.Config.UserName + " must get details before updating.\n")
	}
	ipMap := make(map[string]string, 2)
	if v4 != nil {
		ipMap["ip"] = v4.String()
	}
	if v6 != nil {
		ipMap["ipv6"] = v6.String()
	}
	delta, _, err := self.UApi.UpdateUser(self.userID, ipMap)
	return delta, err
}

func (self *Api) GetDetails() (*swagger.User, error) {
	lookupId := self.userID
	first := false
	if lookupId == "" {
		first = true
		lookupId = self.Config.UserName
		if self.Verbose {
			self.log.Printf("Looking up self details with id: '%s'\n", lookupId)
		}
	}

	data, _, err := self.UApi.UserGet(lookupId)
	if err != nil {
		return nil, err
	}
	if err == nil && first {
		self.userID = fmt.Sprintf("%d", data.Id)
		if self.Verbose {
			self.log.Printf("Got back user id: '%s'\n", self.userID)
		}
	}
	return &data, nil
}

func (self *Api) GetUserDetails(user string) (*swagger.User, error) {
	data, _, err := self.UApi.UserGet(user)
	return &data, err
}

// default config file is ~/.parvati.config
func DefaultConfigFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	home := usr.HomeDir
	return filepath.Join(home, ".parvati.config"), nil
}

func (self *Api) Info() string {
	s := "Connection to: " + self.Config.BasePath
	if self.Config.UserName != "" && self.Config.Password != "" {
		s += " as " + self.Config.UserName
	} else {
		s += " (no credentials)"
	}
	return s
}

func (self *Api) PostUserHost(game *swagger.Game, user *swagger.User, ip net.IP, port uint, hostMessage string) (*swagger.HosterStatus, error) {
	if self.Verbose {
		self.log.Printf("Posting host for user: '%d' for game: '%s' on ip: '%s' port: '%d' in host list\n", user.Id, game.UrlShortName, ip.String(), port)
	}
	stat, _, err := self.HApi.DeclareHost(user.Id, game.UrlShortName, self.announcer, hostMessage, ip, int(port))
	return stat, err
}

// returns whether the user is in the hostlist.
// waiter will always be set if present. Hoster is only set if the user is hosting (otherwise they are waiting)
func (self *Api) UserInHostlist(game *swagger.Game, user *swagger.User) (*swagger.Host, *swagger.Waiter, error) {
	if self.Verbose {
		self.log.Printf("Checking for user: '%d' in host list\n", user.Id)
	}
	list, err := self.CheckListedHosts(game, user)
	if err != nil {
		if self.Verbose {
			self.log.Printf("Getting list failed: '%s'\n", err.Error())
		}
		return nil, nil, err
	}
	for _, host := range list.Hosts {
		if host.Host.BaseInfo.Id == user.Id {
			if self.Verbose {
				self.log.Printf("Found host match\n")
			}
			return &(host.Host), &(host.Host.BaseInfo), nil
		}
		if self.Verbose {
			self.log.Printf("Host id %d did not match %d\n", host.Host.BaseInfo.Id, user.Id)
		}
	}
	for _, wait := range list.Waits {
		if wait.Waiter.Id == user.Id {
			if self.Verbose {
				self.log.Printf("Found waiting user\n")
			}
			return nil, &(wait.Waiter), nil
		}
		if self.Verbose {
			self.log.Printf("Waiting user id %d did not match %d\n", wait.Waiter.Id, user.Id)
		}
	}
	if self.Verbose {
		self.log.Printf("Id %d was not present in the list\n", user.Id)
	}
	return nil, nil, nil
}

// returns items in the hostlist for the given game (required) and user (optional)
func (self *Api) CheckListedHosts(game *swagger.Game, user *swagger.User) (*swagger.HostList, error) {
	userName := ""
	if user != nil {
		userName = fmt.Sprintf("%d", user.Id)
	}
	if self.Verbose {
		self.log.Printf("Checking listed hosts for '%s' with host id: '%s'\n", game.UrlShortName, userName)
	}
	list, _, err := self.HApi.GamesGameIdHostsGet(game.UrlShortName, "", nil, userName, nil)
	if self.Verbose {
		if err != nil {
			self.log.Println("Checking listed hosts failed with error: " + err.Error())
		}
	}

	return list, err
}

// check is one of "basic", "state", "full"
// basic - gives just a quick check on if the host is live
// state - attempts to check who is playing, whether spectate is possible
// full - as above, but includes all current game info, if playing. Not currently supported.
func (self *Api) CheckHosting(game *swagger.Game, user *swagger.User, check string, forcePort uint) (swagger.GameCheckResult, error) {
	var lastErrResult swagger.GameCheckResult
	if len(game.APIs) == 0 {
		return lastErrResult, fmt.Errorf("No test APIs associated with game %s.\n", game.Name)
	}
	if check == "" {
		check = "basic"
	}
	var ipv4, ipv6 net.IP
	for _, proto := range game.Protocols {
		l := len(proto)
		if l == 0 {
			continue
		}
		if proto[l-1] == '4' {
			if ipv4 != nil || user.Ipv4 == "" {
				continue
			}
			ipv4 = net.ParseIP(user.Ipv4)
		} else if proto[l-1] == '6' {
			if ipv6 != nil || user.Ipv6 == "" {
				continue
			}
			ipv6 = net.ParseIP(user.Ipv6)
		}

	}
	ips := make([]net.IP, 0, 2)
	if ipv6 != nil {
		ips = append(ips, ipv6)
	}
	if ipv4 != nil {
		ips = append(ips, ipv4)
	}
	if len(ips) == 0 {
		return lastErrResult, fmt.Errorf("No IPs associated with user %s for game %s.\n", user.Nick, game.Name)
	}
	userPort := forcePort
	if userPort == 0 {
		userPort = uint(user.Port)
		if userPort == 0 {
			userPort = uint(game.Port)
		}
	}
	userPortStr := fmt.Sprintf("%d", userPort)
	var lastErr error
	for _, api := range game.APIs {
		base_uri := api.Uri
		if !strings.HasSuffix(base_uri, "/") {
			base_uri += "/"
		}
		base_uri += "check/"
		for _, ip := range ips {
			hp := net.JoinHostPort(ip.String(), userPortStr)
			uri := base_uri + hp
			if self.Verbose {
				self.log.Printf("Checking host status of: '%s'\n", hp)
			}

			request := resty.R()
			request.SetHeader("X-APIUser", self.Config.UserName)
			request.SetQueryParam("level", check)
			response, err := request.Get(uri)
			if err != nil {
				lastErr = err
				if self.Verbose {
					self.log.Println("Check failed: " + err.Error())
				}
				continue
			}
			var result swagger.GameCheckResult

			err = json.Unmarshal(response.Body(), &result)
			if err != nil {
				if self.Verbose {
					self.log.Println("Host check did not produce valid JSON: " + err.Error())
				}
				return lastErrResult, err
			}
			if response.StatusCode() != 200 || result.Error != "" ||
				result.Info.Status == "Unreachable" || result.Info.Status == "Unknown" {
				if self.Verbose {
					self.log.Printf("Check returned status: %d and error: '%s'", response.StatusCode(), result.Error)
				}
				lastErrResult = result
				continue
			}
			return result, nil
		}

	}
	if lastErrResult.Request == "" {
		return lastErrResult, lastErr
	}
	return lastErrResult, nil
}

func (self *Api) HostAsCheckInfo(host *swagger.Host) *swagger.GameCheckInfo {
	chkLen := len(host.Checks)
	stat := "New"
	if chkLen > 0 {
		stat = host.Checks[chkLen-1].Status
	}
	ip := host.Ipv4
	if ip == "" {
		ip = host.Ipv6
	}
	var spec int
	switch host.Spectateable {
	case "Yes":
		spec = 'y'
	case "No":
		spec = 'n'
	default:
		spec = 'u'
	}
	var uErr = ""
	op := ""
	if host.OpPrivate {
		op = "Anonymous"
	} else if host.Opponent.Id != 0 {
		uidStr := fmt.Sprintf("%d", host.Opponent.Id)
		u, err := self.GetUserDetails(uidStr)
		if err != nil {
			uErr = err.Error()
			op = "id: " + uidStr
		} else {
			op = u.Nick
		}
	}
	g := swagger.GameCheckInfo{
		Address:  net.JoinHostPort(ip, fmt.Sprintf("%d", host.Port)),
		Status:   stat,
		Version:  host.Version,
		Spectate: spec,
		Opponent: op,
		Error:    uErr,
	}
	return &g
}
