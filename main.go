package main

import (
	"context"
	"github.com/atoyr/goflyer/api"
	"github.com/atoyr/goflyer/models"
	"github.com/atoyr/goflyer/db"
	"log"
)

func main() {
	apiClient := api.New("", "")
	var tickerChannl = make(chan models.Ticker)
	var boardCannl = make(chan models.Board)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	board, err := apiClient.GetBoard("BTC_JPY")
	if err != nil {
		log.Println(err)
	}
	d ,err := db.GetBolt("hoge")
	if err != nil {
		log.Println(err)
	}
	d.Init()
	p, err := apiClient.GetPermissions()
	if err != nil {
		log.Println(err)
	}
	for k, v := range p {
		log.Printf("%t : %s", v, k)
	}
	go apiClient.GetRealtimeTicker(ctx, tickerChannl, "BTC_JPY")
	go apiClient.GetRealtimeBoard(ctx, boardCannl, "BTC_JPY", false)
	go func() {
		for b := range boardCannl {
			board.Merge(b)
			log.Printf("action=strealBoard, midPrice: %f", board.MidPrice)
			log.Printf("action=strealBoard, bibs")
			for _, bid := range board.Bids[:20] {
				log.Printf("action=strealBoard, %v", bid)
			}
			log.Printf("action=strealBoard, asks\n")
			for _, ask := range board.Asks[:20] {
				log.Printf("action=strealBoard,%v ", ask)
			}
		}

	}()
	for ticker := range tickerChannl {
		log.Printf("action=StreamIngestionData, %v", ticker)
		d.UpdateTicker(ticker)
	}
}
