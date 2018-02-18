package jupiter

import (
	clientOps "github.com/home-IoT/jupiter/client/client/operations"
	clientModels "github.com/home-IoT/jupiter/client/models"
	srvModels "github.com/home-IoT/jupiter/server/models"
	serverOps "github.com/home-IoT/jupiter/server/restapi/operations"

	"fmt"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/home-IoT/jupiter/client/client"
	"log"
	"sync"
	"time"
)

// latestSensorsData keeps the result of the last read of each sensor
type LatestSensorData struct {
	sync.Mutex
	SensorsData map[string]*srvModels.SensorData
}

var latestSensorsData = LatestSensorData{SensorsData: map[string]*srvModels.SensorData{}}

// sensorsCardList is the list of sensor cards that are known to this service
var sensorsCardList []*srvModels.SensorCard

type HTTPClients struct {
	sync.Mutex
	Clients map[string]*client.Jupiter
}

// httpClients keeps a client per sensor
var httpClients = HTTPClients{Clients: map[string]*client.Jupiter{}}

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
func getSensorData(sensorId string) (*srvModels.SensorData, int) {
	sensorData, ok := latestSensorsData.SensorsData[sensorId]
	if !ok {
		if _, ok := Configuration.Sensors[sensorId]; ok {
			return nil, 504
		} else {
			return nil, 404
		}
	}

	deltaTime := time.Since(time.Time(*sensorData.Timestamp)).Minutes()

	if deltaTime > float64(Configuration.Server.SensorTimeout) {
		fmt.Printf("Removing sensor data for '%s'.", sensorId)

		latestSensorsData.Lock()
		defer latestSensorsData.Unlock()
		delete(latestSensorsData.SensorsData, sensorId)

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

	for id, _ := range Configuration.Sensors {
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
		go fetchSensorData(*sensorData.ID, c, Configuration.Server.SensorReadTime)
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
	c <- copyClientSensorDataToServerSensorData(Configuration.Sensors[id], resp.Payload, latestSensorsData.SensorsData[id])

	log.Printf("Fetched sensor data for '%s'.", id)
}

func getHTTPClient(sensorId string) *client.Jupiter {
	c, ok := httpClients.Clients[sensorId]
	if !ok {
		host := Configuration.Sensors[sensorId].Host + ":" + *Configuration.Sensors[sensorId].Port
		transport := httptransport.New(host, "/", nil)
		c = client.New(transport, strfmt.Default)
		httpClients.Lock()
		defer httpClients.Unlock()

		httpClients.Clients[sensorId] = c
	}
	return c
}

func copyClientSensorDataToServerSensorData(config *SensorConfig, src *clientModels.SensorData, dest *srvModels.SensorData) *srvModels.SensorData {
	if dest == nil {
		dest = new(srvModels.SensorData)
	}

	var stale bool = (*src.Stale != 0)
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

	for id, sensor := range Configuration.Sensors {
		sensorCards = append(sensorCards, newSensorCard(id, sensor.Name))
	}

	sensorsCardList = sensorCards
}
