package executor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/atoyr/goflyer/client"
	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/models"
	"github.com/atoyr/goflyer/configs"
)

// Executor is singleton
type Executor struct {
	dataFrames []models.DataFrame
	db         db.DB
	client     client.APIClient
}

var (
	once sync.Once
	exe  *Executor
)

// GetExecutor is getting executor. executor is singleton
func GetExecutor() *Executor {
	once.Do(func() {
		config ,err := configs.GetGeneralConfig()
		if err != nil {
			panic(err)
		}
		e := new(Executor)
		e.dataFrames = make([]models.DataFrame, 0)
		e.client = *client.New(config.Apikey(),config.Secretkey())

		e.dataFrames = append(e.dataFrames, models.NewDataFrame(models.BTC_JPY, models.GetDuration("1m")))
		e.dataFrames = append(e.dataFrames, models.NewDataFrame(models.BTC_JPY, models.GetDuration("3m")))
		e.dataFrames = append(e.dataFrames, models.NewDataFrame(models.BTC_JPY, models.GetDuration("5m")))
		e.dataFrames = append(e.dataFrames, models.NewDataFrame(models.BTC_JPY, models.GetDuration("10m")))
		e.dataFrames = append(e.dataFrames, models.NewDataFrame(models.BTC_JPY, models.GetDuration("15m")))
		e.dataFrames = append(e.dataFrames, models.NewDataFrame(models.BTC_JPY, models.GetDuration("30m")))
		e.dataFrames = append(e.dataFrames, models.NewDataFrame(models.BTC_JPY, models.GetDuration("1h")))
		e.dataFrames = append(e.dataFrames, models.NewDataFrame(models.BTC_JPY, models.GetDuration("2h")))
		e.dataFrames = append(e.dataFrames, models.NewDataFrame(models.BTC_JPY, models.GetDuration("4h")))
		e.dataFrames = append(e.dataFrames, models.NewDataFrame(models.BTC_JPY, models.GetDuration("6h")))
		e.dataFrames = append(e.dataFrames, models.NewDataFrame(models.BTC_JPY, models.GetDuration("12h")))
		e.dataFrames = append(e.dataFrames, models.NewDataFrame(models.BTC_JPY, models.GetDuration("24h")))
		e.db = config.GetDB()
		exe = e
	})
	return exe
}

func (e *Executor) ChangeDB(db db.DB) {
	e.db = db
}

// GetDataFrame is getting dataframe?
func (e *Executor) GetDataFrame(duration time.Duration) models.DataFrame {
	var df models.DataFrame
	for i := range e.dataFrames {
		if e.dataFrames[i].Duration == duration {
			df = e.dataFrames[i]
		}
	}
	return df
}

func (e *Executor) FetchTickerAsync(ctx context.Context, callbacks []func(beforeeticker, ticker models.Ticker)) {
	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var tickerChannl = make(chan models.Ticker)
	
	before := models.Ticker{}
	go e.client.GetRealtimeTicker(childctx, tickerChannl, "BTC_JPY")
	for ticker := range tickerChannl {
		for i := range callbacks {
			callbacks[i](before, ticker)
		}
		before = ticker
	}
}

func (e *Executor) FetchExecutionAsync(ctx context.Context, callbacks []func(beforeExecution, execution models.Execution)) {
	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var executionChannl = make(chan []models.Execution)
	
	before := models.Execution{}
	go e.client.GetRealtimeExecutions(childctx, executionChannl, "BTC_JPY")
	for executions := range executionChannl{
		for i := range executions{
			for j := range callbacks {
				callbacks[j](before, executions[i])
			}
			before = executions[i]
		}
	}
}

func (e *Executor)Add(datetime time.Time, price, volume float64) {
	for i := range e.dataFrames {
		e.dataFrames[i].Add(datetime,price,volume)
	}
}

func (e *Executor) GetTicker(count int, before, after float64) ([]models.Ticker, error) {
	// TODO Ticker is getting filter
	tickers, err := e.db.GetTickerAll()
	return tickers, err
}

func (e *Executor) SaveTicker(beforeticker,ticker models.Ticker) {
	if ticker.Message == "" {
		e.db.UpdateTicker(ticker)
	}
}

func (e *Executor) GetExecution(count int, before, after float64) ([]models.Execution, error) {
	// TODO Execution is getting filter
	executions, err := e.db.GetExecutionAll()
	return executions, err
}

func (e *Executor) SaveExecution(beforeexecution,execution models.Execution) {
		e.db.UpdateExecution(execution)
}

// GetCandles is getting candles 
func (e *Executor) GetCandles(duration time.Duration) models.Candles {
	df := e.GetDataFrame(duration)
	cs := df.GetCandles()
	return cs
}

func (e *Executor)SaveCandles() {
	// TODO savecandles
/// 	for k := range e.dataFrames {
/// 		df := e.dataFrames[k]
/// 		cs := df.Candles.Candles()
/// 		for i := range cs {
/// 			e.db.UpdateCandle(cs[i])
/// 		}
/// 	}
}

func (e *Executor) MigrationDB(db db.DB) error {
	tickers, err := e.db.GetTickerAll()
	if err != nil {
		return err
	}
	fmt.Printf("Execute Migration...")
	for i := range tickers {
		db.UpdateTicker(tickers[i])
	}
	fmt.Printf("finish!!")
	return nil
}
