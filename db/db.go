package db

import (
	"github.com/atoyr/goflyer/models"
)

type DB interface {
	Init() error
	UpdateTicker(t models.Ticker)
	MergeCandle(candle models.Candle)
	GetCandles(productCode string) ([]models.Candle, error)
}
