# Jupiter - DYI Home Automation Temperature and Humidity Gateway
    
| `master`   | `develop`     |
| ---------- | --------------|
| [![Build Status](https://travis-ci.org/home-IoT/jupiter.svg?branch=master)](https://travis-ci.org/home-IoT/jupiter) | [![Build Status](https://travis-ci.org/home-IoT/jupiter.svg?branch=develop)](https://travis-ci.org/home-IoT/jupiter) |

This is a very simple gateway service to a collection of HT sensors. The server offers a simple REST API to the sensors. The responses follow the the suggestions of the [{json:api}](http://jsonapi.org) specification. 

The service is regularly reading the sensors. The time interval is defined by the `server.sensorReadTime` configuration item. When reading a sensor through this service, the service returns the last reading of the sensor or a `504` if the last reading is too old or if the sensor is unreachable. 

## Server

### REST API
The API of the server is defined by the [`api/server.yml`](api/server.yml) Swagger specification. There are basically two endpoints:

1. `GET /sensors`: provides the list of sensors that the gateway is aware of. 
2. `GET /sernsors/{sensor-id}`: provides the latest reading of the sensor identified by the given id. It returns a `504` error if the reading is older than the defined period in `server.sensorTimeout` configuration, or if the sensor has been unreachablethe. 
3. `GET /sernsors/{sensor-id}/raw`: provides the latest reading of the sensor identified by the given id in a very simple JSON with only two fields (`temperature` and `humidity`). This is useful in connection with other services, such as [Homebridge HTTP Temp/Humidity plugin](https://github.com/lucacri/homebridge-http-temperature-humidity). 

### Configuration
The [`configs/config-template.yml`](configs/config-template.yml) offers a template for the service configuration. 

### Build 

Make sure you that
* you have `dep` installed. Visit https://github.com/golang/dep 
* your `GOPATH` and `GOROOT` environments are set properly.

#### Makefile
There is a [`Makefile`](Makefile) provided that offers a number of targets for preparing, building and running the service. To build and run the service against the [`configs/test.yml`](configs/test.yml) configuration, simply call the `run` target:
```
make clean dep run
```

#### Systemd
I currently have a very basic systemd unit file defined under [`init/jupiter.service`](init/jupiter.service). This can be later improved. 

Before using the service definition, make sure that you go through the file and update the `WorkingDirectory` and `ExecStart` to match your installation. 

## Sensors
The sensors are expected to expose their readings as a JSON response to a simple HTTP GET endpoint defined by the [`api/client.yml`](api/client.yml) Swagger specficitiaon. 

### Mock Sensor
The [home-IoT/esp-dht22-mock](https://github.com/home-IoT/esp-dht22-mock) projects offers a mock sensor service. You can run a number of instances on your local network and test/play with the Jupiter service against the mock sensors. 

To run the server with the configuration provided in [config/test.yml](config/test.yml), run one mock sensor on port `8081` of your machine:
```
dht22-mock --port 8081
```

## License
The code is published under an [MIT license](LICENSE.md). 

## Contributions
Please report issues or feature requests using Github issues. Code contributions can be done using pull requests. 
