package dht22

import (
	"fmt"
	"github.com/go-openapi/swag"
	"github.com/home-IoT/jupiter/dht22-mock/restapi/operations"
	"os"
)

var GitRevision string
var BuildVersion string
var BuildTime string

type MockCommandLineOptions struct {
	Version bool `short:"v" long:"version" description:"Show version"`
}

var MockCommandLineGroup = swag.CommandLineOptionsGroup{
	ShortDescription: "Mock",
	LongDescription:  "Mock options",
	Options:          new(MockCommandLineOptions),
}

var CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
	MockCommandLineGroup,
}

func GetConfigurationOptions(api *operations.Dht22MockAPI) *MockCommandLineOptions {
	for _, v := range api.CommandLineOptionsGroups {
		options, ok := v.Options.(*MockCommandLineOptions)
		if ok {
			return options
		}
	}
	return nil
}

func CheckVersionFlag(api *operations.Dht22MockAPI) {
	options := GetConfigurationOptions(api)

	if options.Version {
		showVersion()
		os.Exit(0)
	}
}

func showVersion() {
	fmt.Printf("app version : %s\n", BuildVersion)
	fmt.Printf("git revision: %s\n", GitRevision)
	fmt.Printf("build time  : %s\n", BuildTime)
}
