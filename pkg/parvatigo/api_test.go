package parvatigo

import (
    "testing"
    "os"
	"net"
)

func getConfig( t *testing.T ) (*Api, *ApiConfig) {
    config_file := os.Getenv( "API_TEST_CONFIG" );
    if config_file == "" {
    	var err error
		config_file, err = DefaultConfigFile()
		if err != nil {
			t.Errorf("Unable to find default config file:\n%s\nAborting test\n", err.Error() )
		}
    }
	config, err := ReadConfig(config_file)
	if err != nil {
		t.Errorf("Unable to read config file: '%s'\n%s\nAborting test\n", config_file, err.Error() )
	}
    if config == nil {
		t.Errorf("Config file: '%s' produce nil configuration\nAborting test\n", config_file )
    }
	api, err := NewApi(config, "dev-test")
	if err != nil {
		t.Errorf("Creating API from Config file: '%s' failed:\n%s\nAborting test\n", config_file, err.Error() )
	}
	return &api, config
}

func FailTest( a *Api, t *testing.T, err error ) {
	baseUri := a.Config.BasePath
	t.Errorf("Failed on test against %s:\n%s", baseUri, err.Error())
}

func TestMakeUser( t *testing.T ) {
    api, _ := getConfig(t)
	ip := net.ParseIP("10.0.0.1")
    _, err := api.MakeUnregisteredUser( ip, 0 )
    if err != nil {
    	FailTest(api,t,err)
    }
}
