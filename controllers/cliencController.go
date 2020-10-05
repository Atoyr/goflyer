package controllers

import (
	"context"
	"sync"
	"time"
  "fmt"

	"github.com/atoyr/goflyer/client"
	"github.com/atoyr/goflyer/client/bitflyer"
	"github.com/atoyr/goflyer/models"
	"github.com/mattn/go-pubsub"
)

type ClientController struct {
	client client.Client
	ps     *pubsub.PubSub
	m      sync.Mutex

	scheduleActions []ScheduleAction
}

type ScheduleAction struct {
	Time   time.Time
	Action Action
}

type Action int

const actionKey = "action"
const (
  EmptyAction = iota
  StartAction
  StopAction
  ExitAction
)

type FetchTickerCallback func(ticker bitflyer.Ticker)

func NewClientController(c client.Client) *ClientController {
	cc := new(ClientController)
	cc.client = c
	cc.ps = pubsub.New()

	// Scheduler
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		ctx = context.WithValue(ctx, actionKey, EmptyAction)
		defer cancel()

		t := time.NewTicker(1 * time.Second)
		defer t.Stop()

		f := func(cc *ClientController) (context.Context, func()) {
			cc.m.Lock()
			ctx, cancel := context.WithDeadline(context.Background(), cc.scheduleActions[0].Time)
			ctx = context.WithValue(ctx, actionKey, cc.scheduleActions[0].Action)
			cc.scheduleActions = cc.scheduleActions[1:]
			cc.m.Unlock()
			return ctx, cancel
		}

		for {
			select {
			case <-ctx.Done():
        action,_ := ctx.Value(actionKey).(Action)
        cc.ps.Pub(action)
				if len(cc.scheduleActions) > 0 {
					ctx, cancel = f(cc)
				} else {
					ctx, cancel = context.WithCancel(context.Background())
					ctx = context.WithValue(ctx, actionKey, EmptyAction)
				}
			case <-t.C:
				if ctx.Value(actionKey) == EmptyAction && len(cc.scheduleActions) > 0 {
					ctx, cancel = f(cc)
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

func (cc *ClientController) RegisterScheduleAction(t time.Time, a Action) {
	cc.m.Lock()
	cc.scheduleActions = append(cc.scheduleActions, ScheduleAction{Time: t, Action: a})
	cc.m.Unlock()
}

func (cc *ClientController) ExecuteFetchTicker(ctx context.Context) {
	childctx, cancel := context.WithCancel(ctx)
	var tickerChannl = make(chan bitflyer.Ticker)
  actionChan := make(chan Action)

  cc.ps.Sub(func(a Action) {
    actionChan <- a
  })

  start := func () {
    go cc.client.GetRealtimeTicker(childctx, tickerChannl, models.BTC_JPY)
  }
  stop := func() {
    cancel()
  }
  

	for {
		select {
		case <-ctx.Done():
        cancel()
			return
		case ticker := <-tickerChannl:
			cc.ps.Pub(ticker)
    case a := <-actionChan:
      fmt.Println(a)
      switch a {
      case StartAction:
        fmt.Println("start")
        start()
      case StopAction:
        fmt.Println("stop")
        stop()
      case ExitAction:
        fmt.Println("exit")
        stop()
        return
      case EmptyAction:
      default:
      }
		}
	}
}
