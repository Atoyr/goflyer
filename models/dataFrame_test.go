package models_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/atoyr/goflyer/models"
)

//
// func TestMergeCandle(t *testing.T) {
// 	cc := NewDataFrame("test", 3*time.Minute)
// 	var candles Candles
// 	candles = make([]Candle, 10)
// 	candles[0] = *NewCandle("test", 3*time.Minute, time.Now().Add(-30*time.Minute), 100, 120, 150, 90, 5)
// 	candles[1] = *NewCandle("test", 3*time.Minute, time.Now().Add(-27*time.Minute), 120, 110, 150, 90, 5)
// 	candles[2] = *NewCandle("test", 3*time.Minute, time.Now().Add(-24*time.Minute), 110, 120, 150, 90, 5)
// 	candles[3] = *NewCandle("test", 3*time.Minute, time.Now().Add(-21*time.Minute), 120, 100, 150, 90, 5)
// 	candles[4] = *NewCandle("test", 3*time.Minute, time.Now().Add(-18*time.Minute), 100, 110, 150, 90, 5)
// 	candles[5] = *NewCandle("test", 3*time.Minute, time.Now().Add(-15*time.Minute), 110, 130, 150, 90, 5)
// 	candles[6] = *NewCandle("test", 3*time.Minute, time.Now().Add(-12*time.Minute), 130, 150, 150, 90, 5)
// 	candles[7] = *NewCandle("test", 3*time.Minute, time.Now().Add(-9*time.Minute), 150, 120, 150, 90, 5)
// 	candles[8] = *NewCandle("test", 3*time.Minute, time.Now().Add(-6*time.Minute), 120, 100, 150, 90, 5)
// 	candles[9] = *NewCandle("test", 3*time.Minute, time.Now().Add(-3*time.Minute), 100, 110, 150, 90, 5)
// 	for _, c := range candles {
// 		cc.MergeCandle(c)
// 	}
//
// }
func TestAddTicker(t *testing.T) {
	jsonFile, err := os.Open("../testdata/tickers.json")
	if err != nil {
		t.Fatal(err)
		t.Fail()
	}
	defer jsonFile.Close()
	raw, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		t.Fatal(err)
		t.Fail()
	}
	cc := models.NewDataFrame("test", 3*time.Minute)
	tickers, err := models.JsonUnmarshalTickers(raw)
	if err != nil {
		t.Fatal(err)
		t.Fail()
	}
	for i := range tickers {
		cc.AddTicker(tickers[i])
	}
}
