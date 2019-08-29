package main

import (
	"context"
	"github.com/atoyr/goflyer/api"
	"github.com/atoyr/goflyer/api/model"
	"log"
)

func main() {
	apiClient := api.New("", "")
	var tickerChannl = make(chan model.Ticker)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go apiClient.GetRealtimeTicker(ctx, tickerChannl, "BTC_JPY")
	for ticker := range tickerChannl {
		log.Printf("action=StreamIngestionData, %v", ticker)
	}
}
