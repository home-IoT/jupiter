package jupiter

import (
	"errors"
	"github.com/home-IoT/jupiter/server/restapi/operations"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// defaultSensorReadTime holds the default reading time of sensors in seconds
const defaultSensorReadTimeInSeconds int64 = 30

// defaultSensorTimeoutInMinutes holds the default sensor read timeout duration in minutes
const defaultSensorTimeoutInMinutes int64 = 60

type SensorConfig struct {
	ID   string  `yaml:"id"`
	Name string  `yaml:"name"`
	Host string  `yaml:"host"`
	Port *string `yaml:"port,omitempty"`
}

type ServerConfig struct {
	SensorReadTime int64 `yaml:"sensorReadTime"`
	SensorTimeout  int64 `yaml:"sensorTimeout"`
}

type jupiterConfigYAML struct {
	Server  ServerConfig    `yaml:"server,omitempty"`
	Sensors []*SensorConfig `yaml:"sensors"`
}

type JupiterConfig struct {
	Sensors map[string]*SensorConfig
	Server  *ServerConfig
}

var Configuration *JupiterConfig

// Configure configures the server with a given configuration file
func Configure(api *operations.JupiterAPI) {
	options := GetConfigurationOptions(api)
	loadSensorsConfig(options.ConfigFile)
}

func loadSensorsConfig(configFile string) (*JupiterConfig, error) {

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config := new(jupiterConfigYAML)

	err = yaml.Unmarshal(file, config)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Error loading the configuration file.")
	}

	Configuration = processConfigYAML(config)

	return Configuration, nil
}

func processConfigYAML(yamlConfig *jupiterConfigYAML) *JupiterConfig {
	config := new(JupiterConfig)
	config.Sensors = make(map[string]*SensorConfig)

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
