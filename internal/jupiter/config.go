package jupiter

import (
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/home-IoT/jupiter/server/restapi/operations"
	"gopkg.in/yaml.v2"
)

// defaultSensorReadTime holds the default reading time of sensors in seconds
const defaultSensorReadTimeInSeconds int64 = 30

// defaultSensorTimeoutInMinutes holds the default sensor read timeout duration in minutes
const defaultSensorTimeoutInMinutes int64 = 60

type sensorConfig struct {
	ID   string  `yaml:"id"`
	Name string  `yaml:"name"`
	Host string  `yaml:"host"`
	Port *string `yaml:"port,omitempty"`
}

type serverConfig struct {
	SensorReadTime int64 `yaml:"sensorReadTime"`
	SensorTimeout  int64 `yaml:"sensorTimeout"`
}

type jupiterConfigYAML struct {
	Server  serverConfig    `yaml:"server,omitempty"`
	Sensors []*sensorConfig `yaml:"sensors"`
}

type jupiterConfig struct {
	Sensors map[string]*sensorConfig
	Server  *serverConfig
}

var configuration *jupiterConfig

// Configure configures the server with a given configuration file
func Configure(api *operations.JupiterAPI) {
	options := getConfigurationOptions(api)

	if options.Version {
		showVersion()
		os.Exit(0)
	}

	if options.ConfigFile == "" {
		printError("Configuration file is missing. Use flag `-c, --config' to provide a config file.")
		os.Exit(1)
	}
	loadSensorsConfig(options.ConfigFile)
}

func loadSensorsConfig(configFile string) (*jupiterConfig, error) {

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config := new(jupiterConfigYAML)

	err = yaml.Unmarshal(file, config)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error loading the configuration file")
	}

	configuration = processConfigYAML(config)

	return configuration, nil
}

func processConfigYAML(yamlConfig *jupiterConfigYAML) *jupiterConfig {
	config := new(jupiterConfig)
	config.Sensors = make(map[string]*sensorConfig)

	for _, v := range yamlConfig.Sensors {
		config.Sensors[v.ID] = v
	}

	config.Server = &yamlConfig.Server

	if config.Server.SensorReadTime <= 0 {
		config.Server.SensorReadTime = defaultSensorReadTimeInSeconds
	}

	if config.Server.SensorTimeout <= 0 {
		config.Server.SensorTimeout = defaultSensorTimeoutInMinutes
	}

	return config
}
