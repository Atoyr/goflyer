package db

import (
	"time"

	"github.com/atoyr/goflyer/models"
)

type DB interface {
	UpdateTicker(ticker models.Ticker) error
	GetTicker(tickerID float64) (models.Ticker, error)
	GetTickerAll() ([]models.Ticker, error)
	UpdateExecution(execution models.Execution) error
	GetExecutionAll() ([]models.Execution, error)
	UpdateDataFrame(models.DataFrame)
	GetDataFrame(duration time.Duration) models.DataFrame
// 	UpdateCandle(duration time.Duration, c models.Candle) error 
// 	GetCandles(duration time.Duration) (models.Candles,error)
	//	MergeCandle(candle models.Candle)
	//	GetCandles(productCode string) ([]models.Candle, error)
}
