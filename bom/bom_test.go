package bom

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	f, err := os.Open("sample.xml")
	if err != nil {
		t.Error(err.Error())
	}
	forecasts, err := parse(f)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if forecasts[0].ChanceOfRain != "60%" {
		t.Error("Expected 60")
	}
	if len(forecasts) != 8 {
		t.Error("Expected 8 forecasts.")
		return
	}
	if forecasts[6].PreciseForecast == nil {
		t.Error("Expected precise forecast")
		return
	}
	if forecasts[7].MinTemperature != 13 {
		t.Errorf("Expected 13, got %d", forecasts[7].MinTemperature)
	}
}
