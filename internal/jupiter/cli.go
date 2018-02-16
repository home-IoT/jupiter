package jupiter

import (
	"github.com/go-openapi/swag"
	"github.com/home-IoT/jupiter/server/restapi/operations"
)

type JupiterCommandLineGroup struct {
	ConfigFile string `short:"c" long:"config" description:"Config file" required:"true"`
}

var GatewayCommandLineGroup = swag.CommandLineOptionsGroup{
	ShortDescription: "Gateway",
	LongDescription:  "Gateway options",
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
