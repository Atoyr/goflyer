package db

import (
	"github.com/atoyr/goflyer/models"
)

type DB interface {
	Init() error
	UpdateTicker(t model.Ticker)
	MergeCandle(candle model.Candle)
	GetCandles(productCode string) ([]model.Candle, error)
}
