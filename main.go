package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/misatosangel/parvati-api-client/commands/generic"
	"github.com/misatosangel/parvati-api-client/commands/lowlevel"
	"github.com/misatosangel/parvati-api-client/commands/parvati"
	"github.com/misatosangel/parvati-api-client/iface"
	"github.com/misatosangel/parvati-api-client/parvatigo"
	"log"
	"os"
)

// Variables used for command line parameters
var settings struct {
	Username   string `short:"u" long:"username" required:"false" description:"Parvati username" value-name:"<nick>"`
	URI        string `long:"uri"  required:"false" description:"Parvati API Uri" value-name:"<url>"`
	ConfigFile string `short:"c" long:"config" required:"false" value-name:"<path>" description:"Location of a gitconfig style file holding your credentials and password and other preferences."`

	Version     func() `long:"version" required:"false" description:"Print tool version and exit."`
	Debug       bool   `short:"d" long:"debug" description:"Debug API load errors."`
	ifaceConfig *iface.Config
}

var buildVersion = "dev"
var buildDate = "dev"
var buildCommit = "dev"

func init() {
}

func main() {
	os.Exit(run())
}

func run() int {
	CliParse()
	return 0

}

func LoadParvatiApi() (*parvatigo.Api, *parvatigo.ApiConfig, error) {
	var config *parvatigo.ApiConfig
	if settings.ConfigFile == "" {
		var err error
		config, err = parvatigo.ReadDefaultConfig()
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, nil, err
			}
			// no config and no parvati credentials, so just do default list show and leave
			def, _ := parvatigo.DefaultConfigFile()
			return nil, nil, fmt.Errorf("No default config file (expected: %s), and no config file given with --config / -c\n"+
				"You must supply a standard gitconfig style file containing at least a value for parvati.password\n", def)
		}
	} else {
		var err error
		config, err = parvatigo.ReadConfig(settings.ConfigFile)
		if err != nil {
			return nil, nil, err
		}
	}
	if config != nil && config.Password == "" {
		return nil, config, fmt.Errorf("Configuration file did not specify a password with key parvati.password\n")
	}
	if settings.URI != "" {
		config.URI = settings.URI
	}
	if settings.Username != "" {
		config.Username = settings.Username
	}
	api, err := parvatigo.NewApi(config, buildVersion)
	if err != nil {
		return nil, config, err
	}
	return &api, config, err
}

func CliParse() {
	parser := flags.NewParser(&settings, flags.Default)
	gaveVersion := false
	settings.Version = func() {
		parser.SubcommandsOptional = true
		fmt.Printf("Parvati client version %s\nBuilt: %s\nCommit: %s\n", buildVersion, buildDate, buildCommit)
		gaveVersion = true
	}
	parser.CommandHandler = func(cmd flags.Commander, args []string) error {
		if gaveVersion {
			return nil
		}
		if cmd == nil {
			return fmt.Errorf("No command given to execute\n")
		}
		if apiCmd, ok := cmd.(cmd_generic.APICommand); ok {
			api, apiConfig, err := LoadParvatiApi()
			if err != nil {
				if apiCmd.NeedsAPI() || (apiCmd.NeedsAPIConfig() && apiConfig == nil) {
					log.Fatalln(err)
				} else if settings.Debug {
					log.Printf("Warning: unable to load parvati api information: %s", err.Error())
				}
			} else {
				if settings.Debug {
					api.Verbose = true
				}
				apiCmd.SetAPIConfig(apiConfig)
				apiCmd.SetAPI(api)
			}
		} else {
			log.Fatalf("Not an API command: %+v\n", cmd)
		}
		if confCmd, ok := cmd.(cmd_generic.IfaceCommand); ok {
			f := settings.ConfigFile
			if f == "" {
				f, _ = parvatigo.DefaultConfigFile()
			}
			confCmd.SetConfigFile(f)
		}
		return cmd.Execute(args)
	}
	_, err := cmd_lowlevel.AddCommands(parser.Command)
	if err != nil {
		log.Fatal(err)
	}
	_, err = cmd_parvati.AddCommands(parser.Command)
	if err != nil {
		log.Fatal(err)
	}
	_, err = parser.Parse()
	if err != nil {
		switch err.(type) {
		case *flags.Error:
			if err.(*flags.Error).Type == flags.ErrHelp {
				os.Exit(0)
			}
		}
		log.Fatalln(err)
	}
}
