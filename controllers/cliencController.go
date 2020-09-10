package controllers

import (
	"context"

	"github.com/atoyr/goflyer/client"
	"github.com/atoyr/goflyer/client/bitflyer"
	"github.com/atoyr/goflyer/models"
)

type ClientController struct {
	client client.Client
}

func NewClientController(c client.Client) *ClientController {
	cc := new(ClientController)
	cc.client = c

	return cc
}

func FetchTickerAsync(ctx context.Context, callbacks []func(beforeeticker, ticker bitflyer.Ticker)) {
	exe := getExecutor()
	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var tickerChannl = make(chan bitflyer.Ticker)

	before := bitflyer.Ticker{}
	go exe.client.GetRealtimeTicker(childctx, tickerChannl, models.BTC_JPY)

	for ticker := range tickerChannl {
		for i := range callbacks {
			callbacks[i](before, ticker)
		}
		before = ticker
	}
}
