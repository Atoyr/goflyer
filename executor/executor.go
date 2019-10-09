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
}

var (
	once sync.Once
	exe  *executor
)

func GetExecutor(db db.DB) *executor {
	once.Do(func() {
		e := new(executor)
		e.dataFrames = make(map[string]models.DataFrame, 0)
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
