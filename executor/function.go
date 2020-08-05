package executor

import (
	"time"
	"context"

	"github.com/atoyr/goflyer/models" 
	"github.com/atoyr/goflyer/models/bitflyer"
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

func FetchTickerAsync(ctx context.Context, callbacks []func(beforeeticker, ticker bitflyer.Ticker)) {
	exe := getExecutor()
	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var tickerChannl = make(chan bitflyer.Ticker)
	
	before := bitflyer.Ticker{}
	go exe.client.GetRealtimeTicker(childctx, tickerChannl, models.BTC_JPY)
	for ticker := range tickerChannl {
		for i := range callbacks {
			callbacks[i](before, ticker)
		}
		before = ticker
	}
}

func FetchExecutionAsync(ctx context.Context, callbacks []func(beforeExecution, execution bitflyer.Execution)) {
	exe := getExecutor()
	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var executionChannl = make(chan []bitflyer.Execution)
	
	before := bitflyer.Execution{}
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

func FetchDataFrameAsync(ctx context.Context) {
	exe := getExecutor()
	childctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var tickerChannl = make(chan bitflyer.Ticker)
	go exe.client.GetRealtimeTicker(childctx, tickerChannl, models.BTC_JPY)
	for ticker := range tickerChannl {
		Add(ticker.DateTime(),ticker.Ltp,ticker.Volume)
	}
}

func Add(datetime time.Time, price, volume float64) {
	exe := getExecutor()
	for i := range exe.dataFrames {
		exe.dataFrames[i].Add(datetime,price,volume)
	}
}

func GetTicker(count int, before, after float64) ([]bitflyer.Ticker, error) {
	exe := getExecutor()
	// TODO Ticker is getting filter
	tickers, err := exe.db.GetTickerAll()
	return tickers, err
}

func SaveTicker(beforeticker,ticker bitflyer.Ticker) {
	exe := getExecutor()
	if ticker.Message == "" {
		exe.db.UpdateTicker(ticker)
	}
}

func GetExecution(count int, before, after float64) ([]bitflyer.Execution, error) {
	exe := getExecutor()
	// TODO Execution is getting filter
	executions, err := exe.db.GetExecutionAll()
	return executions, err
}

func SaveExecution(beforeexecution,execution bitflyer.Execution) {
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
