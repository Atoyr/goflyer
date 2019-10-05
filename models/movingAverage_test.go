package models_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"os"
	"testing"

	"github.com/atoyr/goflyer/models"
)

func getCandle1440() []models.Candle {
	jsonFile, err := os.Open("../testdata/candle_144000.json")
	if err != nil {
		return nil
	}
	defer jsonFile.Close()
	raw, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil
	}
	var cs []models.Candle
	err = json.Unmarshal(raw, &cs)
	if err != nil {
		log.Print(err)
	}

	return cs
}

func getSma1440() models.Sma {
	jsonFile, _ := os.Open("../testdata/sma_6_1440.golden")
	defer jsonFile.Close()
	raw, _ := ioutil.ReadAll(jsonFile)
	var sma models.Sma
	err := json.Unmarshal(raw, &sma)
	if err != nil {
		log.Print(err)
	}

	return sma

}

func TestSmaCreate(t *testing.T) {
	cs := getCandle1440()
	in := make([]float64, len(cs))
	for i := range cs {
		in[i] = cs[i].Close
	}

	period := 6

	sma := models.NewSma(in, period)
	goldenSma := getSma1440()
	for i := range sma.Values {
		if value := (math.Round(sma.Values[i]*1000) / 1000); goldenSma.Values[i] != value {
			t.Fatalf("No %v want %v, but %v:", i, goldenSma.Values[i], value)
		}
	}
}

func TestSmaUpdate(t *testing.T) {
	cs := getCandle1440()
	in := make([]float64, len(cs))
	for i := range cs {
		in[i] = cs[i].Close
	}

	period := 6

	sma := models.NewSma(in[:0], period)
	for i := 1; i < len(in); i++ {
		sma.Update(in[:i])
	}
	goldenSma := getSma1440()
	for i := range sma.Values {
		if value := (math.Round(sma.Values[i]*1000) / 1000); goldenSma.Values[i] != value {
			t.Fatalf("No %v want %v, but %v:", i, goldenSma.Values[i], value)
		}
	}
}
