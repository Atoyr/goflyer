package executor

import (
	"time"
	"context"

	"github.com/atoyr/goflyer/models" 
)

// GetDataFrame is getting dataframe?
func DataFrame(duration time.Duration) models.DataFrame {
	var df models.DataFrame
	exe := getExecutor()
	for i := range exe.dataFrames {
		if exe.dataFrames[i].Duration == duration {
			df = exe.dataFrames[i]
		}
	}
	return df
}

func FetchTickerAsync(ctx context.Context, callbacks []func(beforeeticker, ticker models.Ticker)) {
	exe := getExecutor()
	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var tickerChannl = make(chan models.Ticker)
	
	before := models.Ticker{}
	go exe.client.GetRealtimeTicker(childctx, tickerChannl, models.BTC_JPY)
	for ticker := range tickerChannl {
		for i := range callbacks {
			callbacks[i](before, ticker)
		}
		before = ticker
	}
}

func FetchExecutionAsync(ctx context.Context, callbacks []func(beforeExecution, execution models.Execution)) {
	exe := getExecutor()
	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var executionChannl = make(chan []models.Execution)
	
	before := models.Execution{}
	go exe.client.GetRealtimeExecutions(childctx, executionChannl, models.BTC_JPY)
	for executions := range executionChannl{
		for i := range executions{
			for j := range callbacks {
				callbacks[j](before, executions[i])
			}
			before = executions[i]
		}
	}
}

func Add(datetime time.Time, price, volume float64) {
	exe := getExecutor()
	for i := range exe.dataFrames {
		exe.dataFrames[i].Add(datetime,price,volume)
	}
}

func GetTicker(count int, before, after float64) ([]models.Ticker, error) {
	exe := getExecutor()
	// TODO Ticker is getting filter
	tickers, err := exe.db.GetTickerAll()
	return tickers, err
}

func SaveTicker(beforeticker,ticker models.Ticker) {
	exe := getExecutor()
	if ticker.Message == "" {
		exe.db.UpdateTicker(ticker)
	}
}

func GetExecution(count int, before, after float64) ([]models.Execution, error) {
	exe := getExecutor()
	// TODO Execution is getting filter
	executions, err := exe.db.GetExecutionAll()
	return executions, err
}

func SaveExecution(beforeexecution,execution models.Execution) {
	exe := getExecutor()
		exe.db.UpdateExecution(execution)
}

// GetCandles is getting candles 
func GetCandles(duration time.Duration) models.Candles {
	df := DataFrame(duration)
	cs := df.GetCandles()
	return cs
}

func SaveDataFrame() {
	exe := getExecutor()
	for i := range exe.dataFrames {
		exe.db.UpdateDataFrame(exe.dataFrames[i])
	}
}
