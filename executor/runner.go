package executor

import (
  "context"
	"github.com/atoyr/goflyer/models"
)

func RunAsync(ctx context.Context) {
	executor := GetExecutor()
	callback := make([]func(beforeticker,ticker models.Ticker),0)
	callback = append(callback , func(beforeticker,ticker models.Ticker) {
		executor.AddValue(ticker.DateTime(), ticker.GetMidPrice(),ticker.Volume)
		executor.SaveCandles()
	})
	executor.FetchTickerAsync(ctx, callback )
}
