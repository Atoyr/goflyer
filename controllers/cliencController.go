package controllers

import (
	"context"

	"github.com/atoyr/goflyer/client"
	"github.com/atoyr/goflyer/client/bitflyer"
	"github.com/atoyr/goflyer/models"
)

type ClientController struct {
	client              client.Client
	fetchTickerCallback []FetchTickerCallback
}

type FetchTickerCallback func(before, ticker bitflyer.Ticker)

func NewClientController(c client.Client) *ClientController {
	cc := new(ClientController)
	cc.client = c
	cc.fetchTickerCallback = make([]FetchTickerCallback, 0)

	return cc
}

// RegisterTickerCallback is registed FetchTicker callback
func (cc *ClientController) RegisterTickerCallback(callback FetchTickerCallback) FetchTickerCallback {
	cc.FetchTickerCallback = append(cc.FetchTickerCallback, callback)
}

func (cc *ClientController) UnregisterTickerCallback(callback FetchTickerCallback) {
	if l := len(cc.fetchTickerCallback) - 1; l >= 0 {
		for i, v := range cc.fetchTickerCallback {
			if v == callback {
				if i == 0 && l == 0 {
					cc.fetchTickerCallback = make([]FetchTickerCallback)
				} else if i == 0 {
					cc.fetchTickerCallback = cc.fetchTickerCallback[1:]
				} else if i == l {
					cc.fetchTickerCallback = cc.fetchTickerCallback[:l-1]
				} else {
					cc.fetchTickerCallback = append(cc.fetchTickerCallback[:i-1], cc.fetchTickerCallback[i+1:]...)
				}
			}
		}
	}
}

func (cc *ClientController) FetchTicker(ctx context.Context) {
	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var tickerChannl = make(chan bitflyer.Ticker)

	before := bitflyer.Ticker{}
	go cc.client.GetRealtimeTicker(childctx, tickerChannl, models.BTC_JPY)
	for {
		select {
		case <-ctx.Done():
			return

		case ticker <- tickerChannl:
			for i := range callbacks {
				callbacks[i](before, ticker)
			}
			before = ticker
		}
	}
}
