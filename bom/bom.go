// Grab the weather for melbourne
package bom

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/url"

	"github.com/heath-family/dash-home/ftp_get"
)

// const RadarUrl = "http://www.bom.gov.au/radar/IDR02B.gif"

var WeatherUrl, _ = url.Parse("ftp://ftp2.bom.gov.au/anon/gen/fwo/IDV10450.xml")

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

func Latest() ([]Forecast, error) {
	reader, err := ftp_get.Get(WeatherUrl)
	if err != nil {
		return nil, err
	}
	return parse(reader)
}

func parse(r io.Reader) ([]Forecast, error) {
	v := BomXmlExtractor{}
	err := xml.NewDecoder(r).Decode(&v)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("%+v", v)
}

type BomXmlArea struct {
	// forecast-period
	// <text type="warning_summary_footer">Details of warnings are available on the Bureau's website www.bom.gov.au, by telephone 1300-659-217* or through some TV and radio broadcasts.</text>

}

type BomXmlExtractor struct {
	Version  string `xml:"version,attr"`
	Forecast []struct {
		Content string `xml:",chardata"`
	} `xml:"forecast>area"`
}
