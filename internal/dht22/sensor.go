package dht22

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/home-IoT/jupiter/dht22-mock/models"
	"github.com/home-IoT/jupiter/dht22-mock/restapi/operations"
	"math/rand"
	"os"
)

var sensorData *models.SensorData

func getMockData() *models.SensorData {
	if sensorData == nil {
		newMockData()
	} else {
		vary(sensorData.Dht22.Temperature, 1)
		vary(sensorData.Dht22.Humidity, 1)
		sensorData.Dht22.HeatIndex = CalcHeatIndex(*sensorData.Dht22.Temperature, *sensorData.Dht22.Humidity)
	}
	return sensorData
}

func newMockData() {
	pid := int64(os.Getpid())
	device := fmt.Sprintf("dht-%d", pid)
	rand.Seed(pid)

	deltaTime := int64(0)
	stale := int64(0)

	temperature := (rand.Float64() * 60) - 10
	humidity := (rand.Float64() * 90) + 10
	heatIndex := CalcHeatIndex(temperature, humidity)

	dht := models.SensorDataDht22{Humidity: &humidity, Temperature: &temperature, HeatIndex: heatIndex}
	sensorData = &models.SensorData{Device: &device, DeltaTime: &deltaTime, Stale: &stale, Dht22: &dht}
}

// GetMockData returns a mock sensor data
func GetMockData(params operations.GetParams) middleware.Responder {
	return operations.NewGetOK().WithPayload(getMockData())
}

// SetConfig pretends to set the network configuration of the sensor
func SetConfig(params operations.GetConfigParams) middleware.Responder {
	response := models.ConfigResponse{Ssid: &params.Ssid}
	return operations.NewGetConfigOK().WithPayload(&response)
}
