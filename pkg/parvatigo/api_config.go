package parvatigo

import (
	"fmt"
	"github.com/misatosangel/gitconfig"
	"math/rand"
	"strings"
)

type ApiConfig struct {
	URI       string              `gcKey:"parvati.uri" gcDefault:"https://parvati.phi.al"`
	Username  string              `gcKey:"parvati.username"`
	Password  string              `gcKey:"parvati.password"`
	Announcer string              `gcKey:"parvati.announcer"`
	Games     map[string]GameInfo `gcKey:"game"`
}

type GameInfo struct {
	Name            string `gcKey:"name"`
	ConfigName      string
	HostMessages    []string `gcKey:"hostMessage"`
	WaitMessages    []string `gcKey:"waitMessage"`
	HostOrder       string   `gcKey:"hostMessageOrder" default:"round-robin"`
	WaitOrder       string   `gcKey:"waitMessageOrder" default:"round-robin"`
	Port            uint     `gcKey:"watchPort" gcRequired:"false" gcDefault:"0"`
	Enabled         bool     `gcKey:"enabled" gcDefault:"true"`
	OnJoined        []string `gcKey:"onJoined" gcRequired:"false"`
	lastHostMessage uint
	lastWaitMessage uint
}

func ReadDefaultConfig() (*ApiConfig, error) {
	path, err := DefaultConfigFile()
	if err != nil {
		return nil, err
	}
	return ReadConfig(path)
}

func ReadConfig(file string) (*ApiConfig, error) {
	conf, err := gitconfig.NewConfigFromFile(file)
	if err != nil {
		return nil, err
	}
	var apiConfig ApiConfig
	err = conf.Load(&apiConfig)
	return &apiConfig, err
}

// return a list of games with enabled flags overridden by
// the lists passed in. Only enabled games will be returned.
func (self *ApiConfig) GetEnabledGames(enable, disable []string) []GameInfo {
	out := make([]GameInfo, 0, len(self.Games))
	for n, g := range self.Games {
		if disable != nil {
			shouldDisable := false
			for _, x := range disable {
				if strings.EqualFold(x, g.Name) || strings.EqualFold(x, g.ConfigName) {
					shouldDisable = true
					break
				}
			}
			if shouldDisable {
				continue
			}
		}
		// take copy so we don't alter the enabled flag in our copy
		cpy := g
		cpy.ConfigName = n
		if cpy.Name == "" {
			cpy.Name = n
		}
		if enable != nil {
			for _, x := range disable {
				if strings.EqualFold(x, g.Name) || strings.EqualFold(x, g.ConfigName) {
					cpy.Enabled = true
					break
				}
			}
		}
		if cpy.Enabled { // because we set it, or because it started set
			out = append(out, cpy)
		}
	}
	return out
}

func (self *GameInfo) HostMessage() string {
	max := uint(len(self.HostMessages))
	if max == 0 {
		return ""
	}
	if self.HostOrder == "round-robin" {
		last := self.lastHostMessage
		if last >= max {
			last = 0
		}
		m := self.HostMessages[last]
		self.lastHostMessage = last + 1
		return m
	}
	self.lastHostMessage = uint(rand.Int31n(int32(max)))
	return self.HostMessages[self.lastHostMessage]
}

func (self *GameInfo) WaitMessage() string {
	max := uint(len(self.WaitMessages))
	if max == 0 {
		return ""
	}
	if self.WaitOrder == "round-robin" {
		last := self.lastWaitMessage
		if last >= max {
			last = 0
		}
		m := self.WaitMessages[last]
		self.lastWaitMessage = last + 1
		return m
	}
	self.lastWaitMessage = uint(rand.Int31n(int32(max)))
	return self.WaitMessages[self.lastWaitMessage]
}

func (self *GameInfo) PrettyName() string {
	if self.ConfigName != "" {
		return self.ConfigName
	}
	return self.Name
}

func (self *GameInfo) String() string {
	mHLen := len(self.HostMessages)
	mWLen := len(self.WaitMessages)
	gName := "'" + self.Name + "'"
	if self.Name != self.ConfigName {
		gName += "(" + self.ConfigName + ")"
	}
	s := fmt.Sprintf("Game: %s %d host message(s) picked %s.", gName, mHLen, self.HostOrder)
	if self.Port != 0 {
		s += fmt.Sprintf(" Watching port: %d", self.Port)
	}
	s += "\n"
	if mHLen == 0 && mWLen == 0 {
		s += "No default messages."
		return s
	}
	if mHLen == 0 {
		s += "No default host messages.\n"
	} else {
		if mHLen == 1 {
			s += "Host Message: " + self.HostMessages[0] + "\n"
			return s
		}
		s += "Host Messages:\n"
		for i, m := range self.HostMessages {
			s += fmt.Sprintf("% 2d: %s\n", i+1, m)
		}
	}
	if mWLen == 0 {
		s += "No default wait messages.\n"
	} else {
		if mWLen == 1 {
			s += "Wait Message: " + self.WaitMessages[0] + "\n"
			return s
		}
		s += "Wait Messages:\n"
		for i, m := range self.WaitMessages {
			s += fmt.Sprintf("% 2d: %s\n", i+1, m)
		}
	}
	return s
}
