package main

import (
	"context"
	"github.com/atoyr/goflyer/api"
	"github.com/atoyr/goflyer/models"
	"log"
)

func main() {
	apiClient := api.New("", "")
	var tickerChannl = make(chan models.Ticker)
	var boardCannl = make(chan models.Board)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go apiClient.GetRealtimeTicker(ctx, tickerChannl, "BTC_JPY")
	go apiClient.GetRealtimeBoard(ctx, boardCannl, "BTC_JPY", false)
	go func() {
		for board := range boardCannl {
			log.Printf("action=strealBoard, midPrice: %f", board.MidPrice)
			log.Printf("action=strealBoard, bibs")
			for _, bid := range board.Bids[:10] {
				log.Printf("action=strealBoard, %v", bid)
			}
			log.Printf("action=strealBoard, asks\n")
			for _, ask := range board.Asks[:10] {
				log.Printf("action=strealBoard,%v ", ask)
			}
		}

	}()
	for ticker := range tickerChannl {
		log.Printf("action=StreamIngestionData, %v", ticker)
	}
}
