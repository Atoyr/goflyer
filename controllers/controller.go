package controllers

func StartFetchTickerAsync(ctx context.Context) {
  callback := func(beforeticker,ticker bitflyer.Ticker) {
		Executor.Add(ticker.DateTime(), ticker.Ltp,ticker.Volume)
	})
	FetchTickerAsync(ctx, callback )
}
