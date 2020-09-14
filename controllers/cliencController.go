package controllers

import (
	"context"
  "time"

	"github.com/atoyr/goflyer/client"
	"github.com/atoyr/goflyer/client/bitflyer"
	"github.com/atoyr/goflyer/models"
	"github.com/mattn/go-pubsub"
)

type ClientController struct {
	client client.Client
	ps     *pubsub.PubSub

  scheduleActions []ScheduleAction
}

type ScheduleAction struct {
  Time time.Time
  Action string
}

const actionKey = "action"
const emptyAction = ""

type FetchTickerCallback func(ticker bitflyer.Ticker)

func NewClientController(c client.Client) *ClientController {
	cc := new(ClientController)
	cc.client = c
	cc.ps = pubsub.New()


  // Scheduler
  go func() {
    ctx,cancel := context.WithCancel(context.Background())
    ctx =context.WithValue(ctx,actionKey, emptyAction)
    defer cancel()
    t := time.NewTicker(1 * time.Minute)
    defer t.Stop()

    for {
      select {
      case <-ctx.Done():
        if len(cc.scheduleActions) > 0 {
          ctx, cancel = context.WithDeadline(context.Background(),cc.scheduleActions[0].Time)
          ctx = context.WithValue(ctx, "action", cc.scheduleActions[0].Action)
          cc.scheduleActions = cc.scheduleActions[1:]

        } else {
          ctx,cancel = context.WithCancel(context.Background())
          ctx = context.WithValue(ctx,actionKey, emptyAction)
        }
      case <-t.C :
        if ctx.Value(actionKey) == emptyAction && len(cc.scheduleActions) > 0 {
          ctx, cancel = context.WithDeadline(context.Background(),cc.scheduleActions[0].Time)
          ctx = context.WithValue(ctx, "action", cc.scheduleActions[0].Action)
          cc.scheduleActions = cc.scheduleActions[1:]
        }
      }
    }
  }()

	return cc
}

// RegisterTickerCallback is registed FetchTicker callback
func (cc *ClientController) SubscribeTicker(callback FetchTickerCallback) {
	cc.ps.Sub(callback)
}

func (cc *ClientController) UnsubscribeTicker(callback FetchTickerCallback) {
	cc.ps.Leave(callback)
}

func (cc *ClientController) RegisterTimerAction() {

}

func (cc *ClientController) ExecuteFetchTicker(ctx context.Context) {
	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	var tickerChannl = make(chan bitflyer.Ticker)

	go cc.client.GetRealtimeTicker(childctx, tickerChannl, models.BTC_JPY)
	for {
		select {
		case <-ctx.Done():
			return

		case ticker := <-tickerChannl:
			cc.ps.Pub(ticker)
		}
	}
}
