package db

import (
	"github.com/atoyr/goflyer/models"
)

type DB interface {
	UpdateTicker(ticker models.Ticker) error
	GetTickerAll() ([]models.Ticker, error)
	//	MergeCandle(candle models.Candle)
	//	GetCandles(productCode string) ([]models.Candle, error)
}
