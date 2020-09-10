package controllers

import (
	"log"

	"github.com/atoyr/goflyer/client"
	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/models"
)

type ClientController struct {
	client client.Client
}

func NewClientController(d db.Bolt) *ClientController {
	cc := new(ClientController)

	return cc
}

func (cc *ClientController) ExecuteTickerRoutin(tickerChannel <-chan models.Ticker) {
	for ticker := range tickerChannel {
		log.Printf("action=ExecuteTickerRoutin, %v", ticker)
		cc.DB.UpdateTicker(ticker)
	}
}
