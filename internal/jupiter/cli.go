package jupiter

import (
	"fmt"
	"github.com/go-openapi/swag"
	"github.com/home-IoT/jupiter/server/restapi/operations"
	"os"
)

var GitRevision string
var BuildVersion string
var BuildTime string

type JupiterCommandLineGroup struct {
	Version    bool   `short:"v" long:"version" description:"Show version"`
	ConfigFile string `short:"c" long:"config" description:"Config file"`
}

var GatewayCommandLineGroup = swag.CommandLineOptionsGroup{
	ShortDescription: "Jupiter",
	LongDescription:  "Jupiter options",
	Options:          new(JupiterCommandLineGroup),
}

var CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
	GatewayCommandLineGroup,
}

func GetConfigurationOptions(api *operations.JupiterAPI) *JupiterCommandLineGroup {
	for _, v := range api.CommandLineOptionsGroups {
		options, ok := v.Options.(*JupiterCommandLineGroup)
		if ok {
			return options
		}
	}
	return nil
}

func ShowVersion() {
	fmt.Printf("app version : %s\n", BuildVersion)
	fmt.Printf("git revision: %s\n", GitRevision)
	fmt.Printf("build time  : %s\n", BuildTime)
}

func PrintError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}
