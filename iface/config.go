package iface

import (
	"fmt"
	"github.com/misatosangel/gitconfig"
	"strconv"
)

type Config struct {
	V4Iface string `gcKey:"interfaces.ipv4"`
	V6Iface string `gcKey:"interfaces.ipv6"`
	v4Name  string
	v6Name  string
	V4ID    int
	V6ID    int
}

func ReadConfig(file string) (*Config, error) {
	conf, err := gitconfig.NewConfigFromFile(file)
	if err != nil {
		return nil, err
	}
	var ifaceConfig Config
	err = conf.Load(&ifaceConfig)
	return &ifaceConfig, err
}

func (self *Config) Configure(list *InterfaceList) error {
	if err := self.ConfigureV4(list); err != nil {
		return err
	}
	if err := self.ConfigureV6(list); err != nil {
		return err
	}
	return nil
}

func (self *Config) ConfigureV4(list *InterfaceList) error {
	if self.V4Iface == "" {
		return nil
	}
	v, err := strconv.ParseInt(self.V4Iface, 10, 32)
	if err == nil {
		self.V4ID = int(v)
		return nil
	}
	if list == nil {
		var listErr error
		list, listErr = NewList(0)
		if listErr != nil {
			return listErr
		}
	}
	self.v4Name = self.V4Iface
	val := list.GetInterfaceNumber(self.V4Iface)
	if val == 0 {
		return fmt.Errorf("Unable to parse interface name: '%s', no matching interface.\n", self.V4Iface)
	}
	self.V4ID = val
	return nil
}

func (self *Config) ConfigureV6(list *InterfaceList) error {
	if self.V6Iface == "" {
		return nil
	}
	v, err := strconv.ParseInt(self.V6Iface, 10, 32)
	if err == nil {
		self.V6ID = int(v)
		return nil
	}
	if list == nil {
		var listErr error
		list, listErr = NewList(0)
		if listErr != nil {
			return listErr
		}
	}
	self.v6Name = self.V6Iface
	val := list.GetInterfaceNumber(self.V6Iface)
	if val == 0 {
		return fmt.Errorf("Unable to parse interface name: '%s', no matching interface.\n", self.V6Iface)
	}
	self.V6ID = val
	return nil
}
