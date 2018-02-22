package jupiter

import (
	clientOps "github.com/home-IoT/jupiter/client/client/operations"
	clientModels "github.com/home-IoT/jupiter/client/models"
	srvModels "github.com/home-IoT/jupiter/server/models"
	serverOps "github.com/home-IoT/jupiter/server/restapi/operations"

	"fmt"
	"log"
	"sync"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/home-IoT/jupiter/client/client"
)

type latestSensorDataType struct {
	sync.Mutex
	SensorsData map[string]*srvModels.SensorData
}

// latestSensorsData keeps the result of the last read of each sensor
var latestSensorsData = latestSensorDataType{SensorsData: map[string]*srvModels.SensorData{}}

// sensorsCardList is the list of sensor cards that are known to this service
var sensorsCardList []*srvModels.SensorCard

type httpClientsType struct {
	sync.Mutex
	Clients map[string]*client.Jupiter
}

// httpClients keeps a client per sensor
var httpClients = httpClientsType{Clients: map[string]*client.Jupiter{}}

// StartUpdatingSensorData prepares the sensor data and begins the concurrent threads of updating sensor data
func StartUpdatingSensorData() {
	updateSensorsList()
	go updateSensorsData()
}

// GetSensorsList prepares the sensor list response
func GetSensorsList(params serverOps.GetSensorsListParams) middleware.Responder {
	links := ServerCreateLinksWithSelf(params.HTTPRequest)
	for _, v := range sensorsCardList {
		if v.Links == nil {
			selfLink := fmt.Sprintf("%s/%s", *links.Self, *v.ID)
			v.Links = &srvModels.GenericLinks{Self: &selfLink}
		}
	}
	sensorList := srvModels.SensorList{Data: sensorsCardList, Links: links}
	return serverOps.NewGetSensorsListOK().WithPayload(&sensorList)
}

// getSensorData prepares the last reading of a given sensor
func getSensorData(sensorID string) (*srvModels.SensorData, int) {
	sensorData, ok := latestSensorsData.SensorsData[sensorID]
	if !ok {
		if _, ok := configuration.Sensors[sensorID]; ok {
			return nil, 504
		}
		return nil, 404
	}

	deltaTime := time.Since(time.Time(*sensorData.Timestamp)).Minutes()

	if deltaTime > float64(configuration.Server.SensorTimeout) {
		fmt.Printf("Removing sensor data for '%s'.", sensorID)

		latestSensorsData.Lock()
		defer latestSensorsData.Unlock()
		delete(latestSensorsData.SensorsData, sensorID)

		return nil, 504
	}

	return sensorData, 0
}

// GetSensorData prepares the last reading of a given sensor
func GetSensorData(params serverOps.GetSensorDataParams) middleware.Responder {
	sensorData, errorCode := getSensorData(params.SensorID)

	switch errorCode {
	case 0:
		sensorResponse := srvModels.SensorResponse{Data: sensorData, Links: ServerCreateLinksWithSelf(params.HTTPRequest)}
		return serverOps.NewGetSensorDataOK().WithPayload(&sensorResponse)

	case 504:
		return serverOps.NewGetSensorDataGatewayTimeout()

	case 404:
		fallthrough

	default:
		return serverOps.NewGetSensorDataNotFound()
	}
}

// GetSensorDataRaw prepares the last reading of a given sensor
func GetSensorDataRaw(params serverOps.GetSensorDataRawParams) middleware.Responder {
	sensorData, errorCode := getSensorData(params.SensorID)

	switch errorCode {
	case 0:
		sensorResponse := srvModels.SensorResponseRaw{Temperature: sensorData.Temperature, Humidity: sensorData.Humidity}
		return serverOps.NewGetSensorDataRawOK().WithPayload(&sensorResponse)

	case 504:
		return serverOps.NewGetSensorDataRawGatewayTimeout()

	case 404:
		fallthrough

	default:
		return serverOps.NewGetSensorDataRawNotFound()
	}
}

// StartUpdatingSensorData concurrently updates the sensor data
func updateSensorsData() {
	c := make(chan *srvModels.SensorData)

	for id := range configuration.Sensors {
		go fetchSensorData(id, c, 0)
	}

	for {
		sensorData := <-c

		if sensorData.Temperature != nil {
			latestSensorsData.Lock()
			latestSensorsData.SensorsData[*sensorData.ID] = sensorData
			latestSensorsData.Unlock()

			log.Printf("%s: {%v, %v}", *sensorData.ID, *sensorData.Temperature, *sensorData.Humidity)
		}
		go fetchSensorData(*sensorData.ID, c, configuration.Server.SensorReadTime)
	}
}

// fetchSensorData fetches the sensor data from a sensor
func fetchSensorData(id string, c chan *srvModels.SensorData, delayInSeconds int64) {
	if delayInSeconds > 0 {
		time.Sleep(time.Duration(delayInSeconds * 1000000000))
	}
	httpClient := getHTTPClient(id)
	params := clientOps.NewGetParams()
	resp, err := httpClient.Operations.Get(params)
	if err != nil {
		log.Println(err)
		log.Printf("Cannot access sensor '%s'.", id)

		// if there is an issue with the access,
		// either send back the old data or create an empty one
		if _, ok := latestSensorsData.SensorsData[id]; !ok {
			data := new(srvModels.SensorData)
			data.ID = &id
			c <- data
		} else {
			c <- latestSensorsData.SensorsData[id]
		}

		return
	}
	// if data is received, update the cache (latestSensorsData)
	c <- copyClientSensorDataToServerSensorData(configuration.Sensors[id], resp.Payload, latestSensorsData.SensorsData[id])

	log.Printf("Fetched sensor data for '%s'.", id)
}

func getHTTPClient(sensorID string) *client.Jupiter {
	c, ok := httpClients.Clients[sensorID]
	if !ok {
		host := configuration.Sensors[sensorID].Host + ":" + *configuration.Sensors[sensorID].Port
		transport := httptransport.New(host, "/", nil)
		c = client.New(transport, strfmt.Default)
		httpClients.Lock()
		defer httpClients.Unlock()

		httpClients.Clients[sensorID] = c
	}
	return c
}

func copyClientSensorDataToServerSensorData(config *sensorConfig, src *clientModels.SensorData, dest *srvModels.SensorData) *srvModels.SensorData {
	if dest == nil {
		dest = new(srvModels.SensorData)
	}

	stale := (*src.Stale != 0)
	timeNow := strfmt.DateTime(time.Now())
	dest.ID = &config.ID
	dest.Name = &config.Name

	dest.Stale = &stale

	if src.Dht22.HeatIndex != nil {
		dest.HeatIndex = *src.Dht22.HeatIndex
	} else {
		dest.HeatIndex = 0
	}
	dest.HeatIndex = *src.Dht22.HeatIndex
	dest.Temperature = src.Dht22.Temperature
	dest.Humidity = src.Dht22.Humidity
	dest.DeltaTime = src.DeltaTime
	dest.Timestamp = &timeNow

	return dest
}

func newSensorCard(id string, name string) *srvModels.SensorCard {
	card := srvModels.SensorCard{ID: &id, Name: &name, Links: nil}

	return &card
}

func updateSensorsList() {
	sensorCards := []*srvModels.SensorCard{}

	for id, sensor := range configuration.Sensors {
		sensorCards = append(sensorCards, newSensorCard(id, sensor.Name))
	}

	sensorsCardList = sensorCards
}
