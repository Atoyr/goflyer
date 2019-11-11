package executor

import (
  "context"
	"fmt"

	"github.com/atoyr/goflyer/models"
)

func RunAsync(ctx context.Context) {
	callback := make([]func(beforeticker,ticker models.Ticker),0)
	callback = append(callback , func(beforeticker,ticker models.Ticker) {
		Add(ticker.DateTime(), ticker.Ltp,ticker.Volume)
		SaveDataFrame()
		exe := getExecutor()
		df := exe.dataFrames[1]
		last := len(df.Datetimes) -1
		lastDatetime := ticker.DateTime()
		x := beforeticker.DateTime().Truncate(models.GetDuration(models.Duration_3m))
		y := ticker.DateTime().Truncate(models.GetDuration(models.Duration_3m))
		if !x.Equal(y){
			lastDatetime = beforeticker.DateTime()
			fmt.Println()
		}
		fmt.Printf("\r%s open : %f high : %f low : %f close : %f %s",df.Datetimes[last],df.Opens[last],df.Highs[last],df.Lows[last],df.Closes[last],lastDatetime)
	})
	FetchTickerAsync(ctx, callback )
} 
