// Grab the weather for melbourne
package bom

import (
	"io"
)

type (
	Percentage int
	Celcius    int
	Forecast   struct {
		Description    string
		MaxTemperature Celcius
		MinTemperature Celcius
		Precipitation  Percentage
	}
)

var Latest = Forecast{}

var Sample = []Forecast{
	Forecast{
		Description:    "Hot and sunny",
		MaxTemperature: 36,
		MinTemperature: 12,
		Precipitation:  10,
	},
	Forecast{
		Description:    "Cold and rainy",
		MaxTemperature: 19,
		MinTemperature: 4,
		Precipitation:  90,
	},
}

func parse(io.Reader) ([]Forecast, error) {
	return nil, nil
}
