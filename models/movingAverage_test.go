package models_test

import (
	"testing"
	"io/ioutil"
	"os"
	"log"
	"encoding/json"
	"github.com/atoyr/goflyer/models"
)

func getCandle1440() models.Candles{
	jsonFile, err := os.Open("../testdata/candle_144000.json")
	if err != nil {
		return nil
	}
	defer jsonFile.Close()
	raw, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil
	}
	log.Println(string(raw))
	var cs  []models.Candle
	err = json.Unmarshal(raw,&cs)
	if err != nil {
		log.Print(err)
	}

	return  cs
}

func TestSmaCreate(t *testing.T) {
	cs := getCandle1440()
	t.Log(cs)
	in := make([]float64, len(cs))
	for i := range cs {
		in[i] = cs[i].Close
	}

	period := 6

sma := models.NewSma(in, period)
	for i := range sma.Values {
		//if out[i] != sma.Values[i] {
		//	t.Fatalf("No %v want %v, but %v:", i, out[i], sma.Values[i])
		//}
		t.Log(sma.Values[i])

	}
}
