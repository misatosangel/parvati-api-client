package cmd_parvati

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/misatosangel/parvati-api-client/internal/iface"
	"github.com/misatosangel/parvati-api-client/pkg/parvatigo"
	"github.com/misatosangel/parvati-api-client/pkg/swagger"
	"github.com/misatosangel/traceroute"
	"net"
	"time"
)

func AddCommands(base *flags.Command) (*flags.Command, error) {
	if _, err := (&Details{}).AddCommands(base); err != nil {
		return base, err
	}
	if _, err := (&UpdateIP{}).AddCommands(base); err != nil {
		return base, err
	}
	if _, err := (&HostWatch{}).AddCommands(base); err != nil {
		return base, err
	}
	return base, nil
}

func DumpUserData(user *swagger.User, show_ids, has_admin bool) {
	def := "[not set]"
	bDef := "false"
	if !has_admin {
		def = "[redacted]"
		bDef = def
	}
	fmt.Printf("      Account id: %d / %s [%s]\n"+
		"         Created: %s\n"+
		"    Password set: %s\n"+
		"     Stored IPv4: %s\n"+
		"     Stored IPv6: %s\n"+
		"       Static IP: %s\n"+
		"    Default port: %d\n"+
		"    Private Join: %s\n"+
		"   Stated Gender: %s\n"+
		"      Avatar Url: %s\n"+
		"Registered Email: %s\n", user.Id, user.Nick, user.PrivLevel, user.Created.String(),
		TrueOrDefault(user.HasPassword, bDef), StringOrDefault(user.Ipv4, def),
		StringOrDefault(user.Ipv6, def), TrueOrDefault(user.StaticIP, bDef),
		user.Port,
		TrueOrDefault(user.Private, bDef), StringOrDefault(user.Gender, def),
		StringOrDefault(user.Picture, def), StringOrDefault(user.Email, def))
	if !show_ids {
		return
	}
	if cnt := len(user.Credentials); cnt > 0 {
		fmt.Printf("Identities:\n")
		for i, id := range user.Credentials {
			fmt.Printf(" - % 3d. [%s] %s - %s [Created: %s]\n", i+1, id.AuthRealm.Name, id.Nick, id.Credential, id.Created)
		}
	}
	if cnt := len(user.ChallongeInfo); cnt > 0 {
		d := "[not set]"
		fmt.Printf("Challonge Identities:\n")
		for _, id := range user.ChallongeInfo {
			fmt.Printf(" - % 6d. Username: %s  Email: %s API-Key: %s\n", id.ID,
				StringOrDefault(id.Username, d), StringOrDefault(id.Email, d), StringOrDefault(id.ApiKey, d))
		}
	}

}

func TrueOrDefault(b bool, d string) string {
	if b {
		return "true"
	}
	return d
}

func StringOrDefault(s, d string) string {
	if s == "" {
		return d
	}
	return s
}

func TranslateDate(d string) string {
	t, err := time.Parse("2006-01-02T15:04:05.000000Z", d)
	if err != nil {
		return d
	}
	return t.String()
}

func UpdateIPs(api *parvatigo.Api, ipFlags, v4INum, v6INum int, doIt bool) (*swagger.UserDelta, error) {
	list, err := iface.NewList(ipFlags | traceroute.WANT_LIVE_IP)
	if err != nil {
		return nil, err
	}
	user, err := api.GetDetails()
	if err != nil {
		return nil, err
	}

	var v4, v6 net.IP
	var v4Err, v6Err error

	if ipFlags&traceroute.WANT_PUBLIC_V4 != 0 {
		addr, err := list.GetPublicIP(v4INum, traceroute.WANT_PUBLIC_V4|traceroute.WANT_PRIVATE_V4|traceroute.WANT_LIVE_IP)
		if err != nil {
			v4Err = err
		} else if user.Ipv4 == "" {
			v4 = addr.RemoteIP
		} else {
			uIP := net.ParseIP(user.Ipv4)
			if !addr.RemoteIP.Equal(uIP) {
				v4 = addr.RemoteIP
			}
		}
	}
	if ipFlags&traceroute.WANT_PUBLIC_V6 != 0 {
		addr, err := list.GetPublicIP(v4INum, traceroute.WANT_PUBLIC_V6|traceroute.WANT_PRIVATE_V6|traceroute.WANT_LIVE_IP)
		if err != nil {
			v6Err = err
		} else if user.Ipv4 == "" {
			v6 = addr.RemoteIP
		} else {
			uIP := net.ParseIP(user.Ipv6)
			if !addr.RemoteIP.Equal(uIP) {
				v6 = addr.RemoteIP
			}
		}
	}
	// nothing to do, assuming we didn't error earlier
	if v4 == nil && v6 == nil {
		if v4Err != nil {
			return nil, v4Err
		}
		if v6Err != nil {
			return nil, v6Err
		}
		// not asked for v4 or v6 give empty delta
		delta := &swagger.UserDelta{Player: *user}
		return delta, nil
	}
	if !doIt {
		// construct a manual delta
		delta := &swagger.UserDelta{Player: *user, Delta: &swagger.Delta{}}
		if v4 != nil {
			delta.Delta.IPv4 = []string{user.Ipv4, v4.String()}
		}
		if v6 != nil {
			delta.Delta.IPv4 = []string{user.Ipv6, v6.String()}
		}
		return delta, nil
	}
	d, err := api.UpdateIPs(v4, v6)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func ConfigureIfacePrefs(confFile, inV4, inV6 string) (*iface.Config, error) {
	conf, err := iface.ReadConfig(confFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to read interface information from '%s'\n:%s", confFile, err.Error())
	}
	if inV4 != "" {
		conf.V4Iface = inV4
	}
	if inV6 != "" {
		conf.V6Iface = inV6
	}
	err = conf.Configure(nil)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func ProcessIPDelta(delta *swagger.UserDelta, checkOnly, commentSame bool) error {
	if delta.Delta == nil {
		foundIp := false
		if delta.Player.Ipv4 != "" {
			foundIp = true
			if commentSame {
				fmt.Printf("IPv4 already set to '%s', no need to update\n", delta.Player.Ipv4)
			}
		}
		if delta.Player.Ipv6 != "" {
			foundIp = true
			if commentSame {
				fmt.Printf("IPv6 already set to '%s', no need to update\n", delta.Player.Ipv6)
			}
		}
		if foundIp {
			return nil
		}
		return fmt.Errorf("Unable to find any public IPs for you\n")
	}
	upCnt := 0
	what := "Updated"
	if checkOnly {
		what = "Would update"
	}
	if delta.Delta.IPv4 != nil && len(delta.Delta.IPv4) > 0 {
		fmt.Printf("%s IPv4 from %s to %s\n", what, delta.Delta.IPv4[0], delta.Delta.IPv4[1])
		upCnt++
	}
	if delta.Delta.IPv6 != nil && len(delta.Delta.IPv6) > 0 {
		fmt.Printf("%s IPv6 from %s to %s\n", what, delta.Delta.IPv6[0], delta.Delta.IPv6[1])
		upCnt++
	}
	if upCnt == 0 {
		return fmt.Errorf("No update details returned.\n")
	}
	return nil
}
