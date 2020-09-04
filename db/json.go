package db

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/atoyr/goflyer/client/bitflyer"
	"github.com/atoyr/goflyer/models"
	"github.com/atoyr/goflyer/util"
)

type JsonDB struct {
	tickers []bitflyer.Ticker
}

func GetJsonDB() (JsonDB, error) {
	var jsonDB JsonDB

	return jsonDB, nil
}

func (j *JsonDB) UpdateTicker(bitflyer.Ticker) error {
	return nil
}

func (j *JsonDB) GetTicker(tickID float64) (bitflyer.Ticker, error) {
	tickers, err := j.GetTickerAll()
	if err != nil {
		return bitflyer.Ticker{}, err
	}
	for i := range tickers {
		if tickers[i].TickID == tickID {
			return tickers[i], nil
		}
	}
	return bitflyer.Ticker{}, nil
}

func (j *JsonDB) GetTickerAll() ([]bitflyer.Ticker, error) {
	jsonFile, err := os.Open("./testdata/tickers.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	raw, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	tickers, err := bitflyer.JsonUnmarshalTickers(raw)
	if err != nil {
		return nil, err
	}
	j.tickers = tickers
	return j.tickers, nil
}

func ExportJsonForTickers(db DB, path string) error {
	tickers, err := db.GetTickerAll()
	if err != nil {
		return err
	}
	return util.SaveJsonMarshalIndent(tickers, path)
}

func (j *JsonDB) UpdateExecution(execution bitflyer.Execution) error {
	return nil
}
func (j *JsonDB) GetExecutionAll() ([]bitflyer.Execution, error) {
	return nil, nil
}

func (j *JsonDB) UpdateCandle(duration time.Duration, c models.Candle) error {
	return nil
}

func (j *JsonDB) GetCandles(duration time.Duration) (models.Candles, error) {
	return models.Candles{}, nil
}
func (j *JsonDB) UpdateDataFrame(models.DataFrame) error {
	return nil
}
func (j *JsonDB) GetDataFrame(duration time.Duration) models.DataFrame {
	return models.DataFrame{}
}
