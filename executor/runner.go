package executor

import (
  "context"

	"github.com/atoyr/goflyer/models"
)

func RunAsync(ctx context.Context) {
	callback := make([]func(beforeticker,ticker models.Ticker),0)
	callback = append(callback , func(beforeticker,ticker models.Ticker) {
		Add(ticker.DateTime(), ticker.GetMidPrice(),ticker.Volume)
		SaveDataFrame()
	})
	FetchTickerAsync(ctx, callback )
} 
