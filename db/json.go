package db

import (
	"io/ioutil"
	"os"

	"github.com/atoyr/goflyer/models"
	"github.com/atoyr/goflyer/util"
)

type JsonDB struct {
	tickers []models.Ticker
}

func GetJsonDB() (JsonDB, error) {
	var jsonDB JsonDB

	return jsonDB, nil
}

func (j *JsonDB) UpdateTicker(models.Ticker) error {
	return nil
}

func (j *JsonDB) GetTicker(tickID float64) (models.Ticker, error) {
	tickers, err := j.GetTickerAll()
	if err != nil {
		return models.Ticker{}, err
	}
	for i := range tickers {
		if tickers[i].TickID == tickID {
			return tickers[i], nil
		}
	}
	return models.Ticker{}, nil
}

func (j *JsonDB) GetTickerAll() ([]models.Ticker, error) {
	jsonFile, err := os.Open("./testdata/tickers.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	raw, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	tickers, err := models.JsonUnmarshalTickers(raw)
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
