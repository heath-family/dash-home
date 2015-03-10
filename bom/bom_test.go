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
	if forecasts[0].Precipitation != 60 {
		t.Error("Expected 60")
	}
	if len(forecasts) != 8 {
		t.Error("Expected 8 forecasts.")
	}
	if forecasts[7].MinTemperature != 13 {
		t.Error("Expected 13")
	}
}

// "VIC_PT042" -> melbourne, bit we care about
