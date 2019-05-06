package cmd_lowlevel

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/misatosangel/parvati-api-client/parvatigo"
	"strings"
)

type ConfigHelp struct {
	api          *parvatigo.Api
	ListSections bool `short:"l" long:"ls-sections" description:"Just list (matching) sections."`
	FilePath     bool `short:"p" long:"path" description:"With no section names, just print the default file path amd exit. Otherwise give to print the path at the end of the section info."`
}

func (self *ConfigHelp) AddCommands(base *flags.Command) (*flags.Command, error) {
	c, err := base.AddCommand("ConfigHelp", "Show Config File help", "Use this command to show information on the configuration file. Pass specific section names (or prefixes) to just give docs for those. Pass nothing to show docs for all sections.", self)
	if err != nil {
		return nil, err
	}
	c.Aliases = append(c.Aliases, "help-config")
	return c, err
}

func (self *ConfigHelp) NeedsAPI() bool {
	return false
}

func (self *ConfigHelp) NeedsAPIConfig() bool {
	return false
}

func (self *ConfigHelp) SetAPI(api *parvatigo.Api) {
}

func (self *ConfigHelp) SetAPIConfig(api *parvatigo.ApiConfig) {
}

func (self *ConfigHelp) KnownSections() []string {
	return []string{"parvati", "interfaces", "game"}
}

func (self *ConfigHelp) Execute(args []string) error {
	def, _ := parvatigo.DefaultConfigFile()
	known := self.KnownSections()
	doSections, err := argsToMap(args, known)
	if err != nil {
		return err
	}
	if doSections["!"] {
		if self.FilePath {
			fmt.Printf("%s\n", def)
			return nil
		}
		if !self.ListSections {
			fmt.Print("The default configuration file format is a 'git-config' style file.\n")
			fmt.Print("It can contain the following settings:\n")
			self.FilePath = true // so we print at the end too
		}
	}
	if self.ListSections {
		for _, k := range known {
			if doSections[k] {
				fmt.Println(k)
			}
		}
		if self.FilePath {
			fmt.Printf("The default configuration file path is:\n%s\n", def)
		}
		return nil
	}
	if doSections["parvati"] {
		fmt.Print("Section parvati:\n")
		fmt.Print("  - parvati.username {string}\n")
		fmt.Print("    This is your current username registered to parvati\n\n")
		fmt.Print("  - parvati.password {string}\n")
		fmt.Print("    This is your password, previously registered via e.g. IRC/discord.\n\n")
		fmt.Print("  - parvati.uri {string}\n")
		fmt.Print("    Override the default URI for parvati's backend.\n\n")
		fmt.Print("\n")
	}
	if doSections["interfaces"] {
		fmt.Print("Section interfaces:\n")
		fmt.Print("  - interfaces.ipv4 {string}\n")
		fmt.Print("    Force IPv4 to bind to this interface name or number.\n\n")
		fmt.Print("  - interfaces.ipv6 {string}\n")
		fmt.Print("    Force IPv6 to bind to this interface name or number.\n\n")
		fmt.Print("\n")
	}
	if doSections["game"] {
		fmt.Print("Section game:\n")
		fmt.Print("  - game.NAME.enabled {boolean}\n")
		fmt.Print("    Enable checking of the given game. Defaults to true if other keys\n")
		fmt.Print("    exist.\n\n")
		fmt.Print("  - game.NAME.hostMessage {string}\n")
		fmt.Print("    Host message to use for a game (can be repeated).\n\n")
		fmt.Print("  - game.NAME.hostMessageOrder {string}\n")
		fmt.Print("    Order to use messages (can be either 'round-robin' or 'random').\n\n")
		fmt.Print("  - game.NAME.onJoined {string list}\n")
		fmt.Print("    If defined will attempt to call this program (with optional\n")
		fmt.Print("    arguments) when your host is first joined. The first entry is the\n")
		fmt.Print("    program to run, complete with path as required. Any extra strings\n")
		fmt.Print("    are arguments to the program, one per entry in order.\n")
		fmt.Print("    The following substitutions will be made before running to args:\n")
		fmt.Print("      ${NICK} - will be replaced by the opponent's Parvati NickName'.\n\n")
		fmt.Print("  - game.NAME.watchPort {integer}\n")
		fmt.Print("    Override your online default port with this one to check for\n")
		fmt.Print("    hosting.\n\n")
	}
	if self.FilePath {
		fmt.Printf("The default configuration file path is:\n%s\n", def)
	}
	return nil
}

func argsToMap(args, known []string) (map[string]bool, error) {
	out := make(map[string]bool)
	if len(args) == 0 { // do everything if not sections given
		for _, s := range known {
			out[s] = true
		}
		out["!"] = true
		return out, nil
	}
	knownSections := make(map[string]bool)
	for _, s := range known {
		knownSections[s] = true
	}

	for _, arg := range args {
		arg = strings.ToLower(arg)
		if arg == "all" {
			for _, s := range known {
				out[s] = true
			}
			return out, nil
		}
		if knownSections[arg] {
			out[arg] = true
			continue
		}
		completes := shortestCompletions(arg, known)
		if len(completes) == 0 {
			ksecs := strings.Join(known, ", ")
			return nil, fmt.Errorf("Unknown section name or prefix: '%s'\nKnown sections are: %s\n", arg, ksecs)
		}
		for _, c := range completes {
			out[c] = true
			continue
		}
	}
	return out, nil
}

func shortestCompletions(arg string, known []string) []string {
	out := make([]string, 0, len(known))
	for _, s := range known {
		if strings.HasPrefix(s, arg) {
			out = append(out, s)
		}
	}
	return out
}
