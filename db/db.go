package db

import (
	"time"

	"github.com/atoyr/goflyer/client/bitflyer"
	"github.com/atoyr/goflyer/models"
)

type DB interface {
	UpdateTicker(ticker bitflyer.Ticker) error
	GetTicker(tickerID float64) (bitflyer.Ticker, error)
	GetTickerAll() ([]bitflyer.Ticker, error)
	UpdateExecution(execution bitflyer.Execution) error
	GetExecutionAll() ([]bitflyer.Execution, error)
	UpdateDataFrame(models.DataFrame) error
	GetDataFrame(duration time.Duration) models.DataFrame
	// 	UpdateCandle(duration time.Duration, c models.Candle) error
	// 	GetCandles(duration time.Duration) (models.Candles,error)
	//	MergeCandle(candle models.Candle)
	//	GetCandles(productCode string) ([]models.Candle, error)
}
