package dht22

import (
	"math"
	"math/rand"
)

// CalcHeatIndex calculates the heat index based on temp and humidity.
func CalcHeatIndex(tempC float64, humidity float64) *float64 {
	// Based on a code from the https://github.com/chrissnell/gopherwx project.
	// Thanks to https://github.com/chrissnell

	temp := tempCtoF(tempC)

	// Heat indices don't make much sense at temps below 77° F
	if temp < 77 {
		return nil
	}

	// First, we try Steadman's method, which is valid for all heat indices
	// below 80° F
	hi := 0.5 * (temp + 61.0 + ((temp - 68.0) * 1.2) + (humidity + 0.094))
	if hi < 80 {
		// Only return heat index if it's greater than the temperature
		if hi > temp {
			hi = tempFtoC(hi)
			return &hi
		}
		return nil
	}

	// Our heat index is > 80, so we need to use the Rothfusz method instead
	c1 := -42.379
	c2 := 2.04901523
	c3 := 10.14333127
	c4 := 0.22475541
	c5 := 0.00683783
	c6 := 0.05481717
	c7 := 0.00122874
	c8 := 0.00085282
	c9 := 0.00000199

	hi = c1 + (c2 * temp) + (c3 * humidity) - (c4 * temp * humidity) - (c5 * math.Pow(temp, 2)) - (c6 * math.Pow(humidity, 2)) + (c7 * math.Pow(temp, 2) * humidity) + (c8 * temp * math.Pow(humidity, 2)) - (c9 * math.Pow(temp, 2) * math.Pow(humidity, 2))

	// If RH < 13% and temperature is between 80 and 112, we need to subtract an adjustment
	if humidity < 13 && temp >= 80 && temp <= 112 {
		adj := ((13 - humidity) / 4) * math.Sqrt((17-math.Abs(temp-95.0))/17)
		hi = hi - adj
	} else if humidity > 80 && temp >= 80 && temp <= 87 {
		// Likewise, if RH > 80% and temperature is between 80 and 87, we need to add an adjustment
		adj := ((humidity - 85.0) / 10) * ((87.0 - temp) / 5)
		hi = hi + adj
	}

	// Only return heat index if it's greater than the temperature
	if hi > temp {
		hi = tempFtoC(hi)
		return &hi
	}
	return nil
}

func vary(variable *float64, maxVariance float64) {
	delta := rand.Float64()*2*maxVariance - maxVariance
	*variable = *variable + delta
}

func tempCtoF(tempC float64) float64 {
	return tempC*1.8 + 32
}

func tempFtoC(tempF float64) float64 {
	return (tempF - 32) / 1.8
}
