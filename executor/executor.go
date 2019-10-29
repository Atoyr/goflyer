package executor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/atoyr/goflyer/client"
	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/models"
)

// Executor is singleton
type Executor struct {
	dataFrames models.DataFrames
	db         db.DB
	client     client.APIClient
}

var (
	once sync.Once
	exe  *Executor
)

// GetExecutor is getting executor. executor is singleton
func GetExecutor(db db.DB) *Executor {
	once.Do(func() {
		e := new(Executor)
		e.dataFrames = make(map[string]models.DataFrame, 0)
		e.client = *client.New("", "")

		e.dataFrames["1m"] = models.NewDataFrame(models.BTC_JPY, models.GetDuration("1m"))
		e.dataFrames["3m"] = models.NewDataFrame(models.BTC_JPY, models.GetDuration("3m"))
		e.dataFrames["5m"] = models.NewDataFrame(models.BTC_JPY, models.GetDuration("5m"))
		e.dataFrames["10m"] = models.NewDataFrame(models.BTC_JPY, models.GetDuration("10m"))
		e.dataFrames["15m"] = models.NewDataFrame(models.BTC_JPY, models.GetDuration("15m"))
		e.dataFrames["30m"] = models.NewDataFrame(models.BTC_JPY, models.GetDuration("30m"))
		e.dataFrames["1h"] = models.NewDataFrame(models.BTC_JPY, models.GetDuration("1h"))
		e.dataFrames["2h"] = models.NewDataFrame(models.BTC_JPY, models.GetDuration("2h"))
		e.dataFrames["4h"] = models.NewDataFrame(models.BTC_JPY, models.GetDuration("4h"))
		e.dataFrames["6h"] = models.NewDataFrame(models.BTC_JPY, models.GetDuration("6h"))
		e.dataFrames["12h"] = models.NewDataFrame(models.BTC_JPY, models.GetDuration("12h"))
		e.dataFrames["24h"] = models.NewDataFrame(models.BTC_JPY, models.GetDuration("24h"))
		exe = e
	})
	exe.db = db
	return exe
}

// RunClient is running Executor
func RunClient() {
	client := client.New("", "")
	var tickerChannl = make(chan models.Ticker)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client.GetRealtimeTicker(ctx, tickerChannl, "BTC_JPY")
	for ticker := range tickerChannl {
		fmt.Println(ticker)
	}
}

// GetDataFrame is getting dataframe?
func (e *Executor) GetDataFrame(key string) models.DataFrame {
	if df, ok := e.dataFrames[key]; ok {
		return df
	}
	return e.dataFrames["24h"]
}

func (e *Executor) GetCandleOHLCs(key string) []models.CandleOHLC {
	var dataFrame models.DataFrame
	if df, ok := e.dataFrames[key]; ok {
		dataFrame = df
	} else {
		dataFrame = e.dataFrames["24h"]
	}
	return dataFrame.Candles.GetCandleOHLCs()
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

func (e *Executor)AddValue(datetime time.Time, id, price, volume float64) {
	for k := range e.dataFrames {
		e.dataFrames[k].AddValue(datetime , id, price, volume )
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

func (e *Executor) SaveExecution(beforeexecution,execution models.Execution) {
		e.db.UpdateExecution(execution)
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
