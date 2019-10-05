package main

import (
	///	"context"

	//"github.com/atoyr/goflyer/client"
	"context"
	"log"
	"path/filepath"

	"github.com/atoyr/goflyer/client"
	"github.com/atoyr/goflyer/controllers"
	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/models"
	"github.com/atoyr/goflyer/util"
)

func main() {
	clientClient := client.New("", "")
	dirPath, err := util.CreateConfigDirectoryIfNotExists("goflyer")
	if err != nil {
		log.Println(err)
	}
	var tickerChannl = make(chan models.Ticker)
	//	var boardCannl = make(chan models.Board)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//	board, err := clientClient.GetBoard("BTC_JPY")
	if err != nil {
		log.Println(err)
	}
	dbfile := filepath.Join(dirPath, "goflyer.db")
	d, err := db.GetBolt(dbfile)
	if err != nil {
		log.Println(err)
	}
	d.Init()

	cc := controllers.NewClientController(d)
	go clientClient.GetRealtimeTicker(ctx, tickerChannl, "BTC_JPY")
	cc.ExecuteTickerRoutin(tickerChannl)

	//tickers, err := d.GetAllTicker()
	//if err != nil {
	//	log.Println(err)
	//}
	//	for _, t := range tickers {
	//		log.Println(t)
	//	}

	//	p, err := clientClient.GetPermissions()
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	for k, v := range p {
	//		log.Printf("%t : %s", v, k)
	//	}

	// cc := models.NewDataFrame("test", 3*time.Minute)
	// cc.AddSmas(3)
	// cc.AddEmas(3)
	// cc.AddMacd(2, 4, 4)
	// cc.AddRsis(2)
	// candles := make([]models.Candle, 10)
	// candles[0] = *models.NewCandle("test", 3*time.Minute, time.Now().Add(-30*time.Minute), 100, 120, 150, 90, 5)
	// candles[1] = *models.NewCandle("test", 3*time.Minute, time.Now().Add(-27*time.Minute), 120, 110, 150, 90, 5)
	// candles[2] = *models.NewCandle("test", 3*time.Minute, time.Now().Add(-24*time.Minute), 110, 120, 150, 90, 5)
	// candles[3] = *models.NewCandle("test", 3*time.Minute, time.Now().Add(-21*time.Minute), 120, 100, 150, 90, 5)
	// candles[4] = *models.NewCandle("test", 3*time.Minute, time.Now().Add(-18*time.Minute), 100, 110, 150, 90, 5)
	// candles[5] = *models.NewCandle("test", 3*time.Minute, time.Now().Add(-15*time.Minute), 110, 130, 150, 90, 5)
	// candles[6] = *models.NewCandle("test", 3*time.Minute, time.Now().Add(-12*time.Minute), 130, 150, 150, 90, 5)
	// candles[7] = *models.NewCandle("test", 3*time.Minute, time.Now().Add(-9*time.Minute), 150, 120, 150, 90, 5)
	// candles[8] = *models.NewCandle("test", 3*time.Minute, time.Now().Add(-6*time.Minute), 120, 100, 150, 90, 5)
	// candles[9] = *models.NewCandle("test", 3*time.Minute, time.Now().Add(-3*time.Minute), 100, 110, 150, 90, 5)
	// for _, c := range candles {
	// 	cc.MergeCandle(c)
	// }
	// ccs := models.DataFrames{}
	// ccs["hoge"] = cc
	// e := api.AppendHandler(api.GetEcho(ccs))
	// e.Start(":8080")

	//	go clientClient.GetRealtimeTicker(ctx, tickerChannl, "BTC_JPY")
	//	go clientClient.GetRealtimeBoard(ctx, boardCannl, "BTC_JPY", false)
	//	go func() {
	//		for b := range boardCannl {
	//			board.Merge(b)
	//			log.Printf("action=strealBoard, midPrice: %f", board.MidPrice)
	//			log.Printf("action=strealBoard, bibs")
	//			for _, bid := range board.Bids[:20] {
	//				log.Printf("action=strealBoard, %v", bid)
	//			}
	//			log.Printf("action=strealBoard, asks\n")
	//			for _, ask := range board.Asks[:20] {
	//				log.Printf("action=strealBoard,%v ", ask)
	//			}
	//		}
	//
	//	}()
	//for ticker := range tickerChannl {
	//	log.Printf("action=StreamIngestionData, %v", ticker)
	//	d.UpdateTicker(ticker)
	//}
}
