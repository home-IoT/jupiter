package dht22

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/home-IoT/jupiter/dht22-mock/models"
	"github.com/home-IoT/jupiter/dht22-mock/restapi/operations"
)

func newMockData() *models.SensorData {
	device := "dht-04"
	deltaTime := int64(0)
	stale := int64(0)

	heatIndex := float64(23.4)
	humidity := float64(65.1)
	temperature := float64(22.4)

	dht := models.SensorDataDht22{Humidity: &humidity, Temperature: &temperature, HeatIndex: &heatIndex}

	return &models.SensorData{Device: &device, DeltaTime: &deltaTime, Stale: &stale, Dht22: &dht}
}

func GetMockData(params operations.GetParams) middleware.Responder {
	return operations.NewGetOK().WithPayload(newMockData())
}

func SetConfig(params operations.GetConfigParams) middleware.Responder {
	response := models.ConfigResponse{Ssid: &params.Ssid}
	return operations.NewGetConfigOK().WithPayload(&response)
}
