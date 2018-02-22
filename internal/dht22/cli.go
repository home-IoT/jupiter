package dht22

import (
	"fmt"
	"os"

	"github.com/go-openapi/swag"
	"github.com/home-IoT/jupiter/dht22-mock/restapi/operations"
)

// GitRevision holds the git revision based on which this service is compiled
var GitRevision string

// BuildVersion holds the version of this service
var BuildVersion string

// BuildTime holds the time of build
var BuildTime string

type mockCommandLineOptions struct {
	Version bool `short:"v" long:"version" description:"Show version"`
}

var mockCommandLineGroup = swag.CommandLineOptionsGroup{
	ShortDescription: "Mock",
	LongDescription:  "Mock options",
	Options:          new(mockCommandLineOptions),
}

// CommandLineOptionsGroups holds the CL option groups
var CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
	mockCommandLineGroup,
}

func getConfigurationOptions(api *operations.Dht22MockAPI) *mockCommandLineOptions {
	for _, v := range api.CommandLineOptionsGroups {
		options, ok := v.Options.(*mockCommandLineOptions)
		if ok {
			return options
		}
	}
	return nil
}

// CheckVersionFlag checks if the version flag is set, prints the version and stops the process
func CheckVersionFlag(api *operations.Dht22MockAPI) {
	options := getConfigurationOptions(api)

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
