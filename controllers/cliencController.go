package controllers

import (
	"log"

	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/models"
)

type ClientController struct {
	DB db.Bolt
}

func NewClientController(bolt db.Bolt) *ClientController {
	cc := new(ClientController)
	cc.DB = bolt

	return cc
}

func (cc *ClientController) ExecuteTickerRoutin(tickerChannel <-chan models.Ticker) {
	for ticker := range tickerChannel {
		log.Printf("action=ExecuteTickerRoutin, %v", ticker)
		cc.DB.UpdateTicker(ticker)
	}
}
