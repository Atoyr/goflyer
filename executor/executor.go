package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/atoyr/goflyer/client"
	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/models"
)

type executor struct {
	dataFrames models.DataFrames
	db         db.DB
	client     client.APIClient
}

var (
	once sync.Once
	exe  *executor
)

func GetExecutor(db db.DB) *executor {
	once.Do(func() {
		e := new(executor)
		e.dataFrames = make(map[string]models.DataFrame, 0)
		e.client = *client.New("", "")

		e.dataFrames["3m"] = models.NewDataFrame(models.BTC_JPY, models.GetDuration("3m"))
		e.dataFrames["24h"] = models.NewDataFrame(models.BTC_JPY, models.GetDuration("24h"))
		exe = e
	})
	exe.db = db
	return exe
}

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

func (e *executor) GetDataFrame(key string) models.DataFrame {
	if df, ok := e.dataFrames[key]; ok {
		tickers, _ := e.db.GetTickerAll()
		for i := range tickers {
			df.AddTicker(tickers[i])
		}
		return df
	}
	return e.dataFrames["24h"]
}

func (e *executor) GetCandleOHLCs(key string) []models.CandleOHLC {
	var cs models.Candles
	var dataFrame models.DataFrame
	if df, ok := e.dataFrames[key]; ok {
		tickers, _ := e.db.GetTickerAll()
		for i := range tickers {
			df.AddTicker(tickers[i])
		}
		dataFrame = df
	} else {
		dataFrame = e.dataFrames["24h"]
	}
	cs = dataFrame.Candles
	return cs.GetCandleOHLCs()
}

func (e *executor) FetchTickerAsync(ctx context.Context, callbacks []func(models.Ticker)) {
	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var tickerChannl = make(chan models.Ticker)

	go e.client.GetRealtimeTicker(childctx, tickerChannl, "BTC_JPY")
	for ticker := range tickerChannl {
		e.db.UpdateTicker(ticker)
		for i := range callbacks {
			callbacks[i](ticker)
		}
	}
}

func (e *executor) GetTicker(count int, before, after float64) ([]models.Ticker, error) {
	tickers, err := e.db.GetTickerAll()
	return tickers, err
}

func (e *executor) MigrationDB(db db.DB) error {
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
