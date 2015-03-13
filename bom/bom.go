// Grab the weather for melbourne

// WIP: Tried using goquery / html parser.
// Doesn't seem to be working; probably go back to encoding/xml and structs.
package bom

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"

	"github.com/heath-family/dash-home/ftp_get"
)

// const RadarUrl = "http://www.bom.gov.au/radar/IDR02B.gif"

var WeatherUrl, _ = url.Parse("ftp://ftp2.bom.gov.au/anon/gen/fwo/IDV10450.xml")

type (
	Celcius  int
	Forecast struct {
		Date               string
		Description        string
		ChanceOfRain       string
		PrecipitationRange string
		*PreciseForecast
	}
	PreciseForecast struct {
		MaxTemperature Celcius
		MinTemperature Celcius
	}
)

var Sample = []Forecast{
	Forecast{
		Date:            "3 Jan",
		Description:     "Hot and sunny",
		ChanceOfRain:    "10%",
		PreciseForecast: nil,
	},
	Forecast{
		Date:               "4 Jan",
		Description:        "Cold and rainy",
		PrecipitationRange: "1-2mm",
		PreciseForecast: &PreciseForecast{
			MaxTemperature: 36,
			MinTemperature: 12,
		},
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
	var area *BomXmlForecast

	v := BomXmlExtractor{}
	err := xml.NewDecoder(r).Decode(&v)
	if err != nil {
		return nil, err
	}
	for _, f := range v.Forecast {
		if f.Aac == "VIC_PT042" { // Melbourne metro
			area = &f
			break
		}
	}
	if area == nil {
		return nil, fmt.Errorf("Could not find area VIC_PT042 in BOM xml")
	}

	forecasts := make([]Forecast, len(area.ForecastPeriod))
	for i, day := range area.ForecastPeriod {
		precis, err := day.Text.Find("precis")
		if err != nil {
			return nil, fmt.Errorf("Couldn't find a precis for forecast period %d; %s", i, err)
		}
		precipitation, _ := day.Text.Find("probability_of_precipitation")
		precipRange, _ := day.Text.Find("precipitation_range")
		forecast, _ := day.Elements.ToPreciseForecast()

		forecasts[i] = Forecast{
			// Month:              day.StartTimeLocal.Month(),
			// Day:                day.StartTimeLocal.Day(),
			Date:               day.StartTimeLocal.Format("02 Jan"),
			Description:        precis,
			ChanceOfRain:       precipitation,
			PrecipitationRange: precipRange,
			PreciseForecast:    forecast,
		}
	}
	return forecasts, nil
}

type BomXmlText struct {
	Content string `xml:",chardata"`
	Type    string `xml:"type,attr"`
}

type BomXmlTexts []BomXmlText

func (t BomXmlTexts) ToPreciseForecast() (*PreciseForecast, error) {
	maxTemp, err := t.GetNum("air_temperature_maximum")
	if err != nil {
		return nil, err
	}
	minTemp, err := t.GetNum("air_temperature_minimum")
	if err != nil {
		return nil, err
	}

	return &PreciseForecast{
		MaxTemperature: Celcius(maxTemp),
		MinTemperature: Celcius(minTemp),
	}, nil
}

type ErrKeyMissing string

func (e ErrKeyMissing) Error() string {
	return "Did not find key " + string(e)
}

func (t BomXmlTexts) GetNum(key string) (int, error) {
	text, err := t.Find(key)
	if err == nil {
		result, err := strconv.ParseInt(text[:2], 10, 8)
		if err != nil {
			return 0, err
		} else {
			return int(result), nil
		}
	} else {
		return 0, ErrKeyMissing(key)
	}
}

func (t BomXmlTexts) Find(key string) (string, error) {
	for _, txt := range t {
		if txt.Type == key {
			return txt.Content, nil
		}
	}
	return "", ErrKeyMissing(key)
}

type BomXmlForecast struct {
	// Content string `xml:",chardata"`
	Aac            string `xml:"aac,attr"`
	ForecastPeriod []struct {
		Index          int         `xml:"index,attr"`
		StartTimeLocal time.Time   `xml:"start-time-local,attr"`
		EndTimeLocal   time.Time   `xml:"end-time-local,attr"`
		Elements       BomXmlTexts `xml:"element"`
		Text           BomXmlTexts `xml:"text"`
		// start-time-local="2015-02-22T17:00:00+11:00" end-time-local="2015-02-23T00:00:00+11:00" start-time-utc="2015-02-22T06:00:00Z" end-time-utc="2015-02-22T13:00:00Z">
	} `xml:"forecast-period"`
}

type BomXmlExtractor struct {
	Version  string           `xml:"version,attr"`
	Forecast []BomXmlForecast `xml:"forecast>area"`
}
