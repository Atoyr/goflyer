package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/atoyr/goflyer/models/bitflyer"
)

// GetRealtimeTicker get Ticker for websocket
func (api *APIClient) GetRealtimeTicker(ctx context.Context, ch chan<- bitflyer.Ticker, productCode string) {
	jsonRPC2 := NewJsonRPC2Subscribe()

	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	jsonRPC2.Params = SubscribeParams{Channel: fmt.Sprintf("lightning_ticker_%s", productCode)}

	var paramCh = make(chan interface{})
	go api.doWebsocketRequest(childctx, *jsonRPC2, paramCh)

OUTER:
	for {
		select {
		case <-ctx.Done():
			return

		default:
			param := <-paramCh
			ticker := new(bitflyer.Ticker)
			marchalTick, err := json.Marshal(param)
			if err != nil {
				ticker.Message = err.Error()
				ch <- *ticker
				continue OUTER
			}
			if err := json.Unmarshal(marchalTick, &ticker); err != nil {
				ticker.Message = err.Error()
				ch <- *ticker
				continue OUTER
			}
			ch <- *ticker
		}
	}
}

func (api *APIClient) GetRealtimeBoard(ctx context.Context, ch chan<- bitflyer.Board, productCode string, isDiff bool) {
	jsonRPC2 := NewJsonRPC2Subscribe()

	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if isDiff {
		jsonRPC2.Params = SubscribeParams{Channel: fmt.Sprintf("lightning_board_%s", productCode)}
	} else {
		jsonRPC2.Params = SubscribeParams{Channel: fmt.Sprintf("lightning_board_snapshot_%s", productCode)}

	}

	var paramCh = make(chan interface{})
	go api.doWebsocketRequest(childctx, *jsonRPC2, paramCh)

OUTER:
	for {
		select {
		case <-ctx.Done():
			return

		default:
			param := <-paramCh
			marchalBoard, err := json.Marshal(param)
			if err != nil {
				log.Printf("error : %s", err)
				continue OUTER
			}
			board := new(bitflyer.Board)
			if err := json.Unmarshal(marchalBoard, &board); err != nil {
				log.Printf("error : %s", err)
				continue OUTER
			}
			ch <- *board
		}
	}
}

func (api *APIClient) GetRealtimeExecutions(ctx context.Context, ch chan<- []bitflyer.Execution, productCode string) {
	jsonRPC2 := NewJsonRPC2Subscribe()

	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	jsonRPC2.Params = SubscribeParams{Channel: fmt.Sprintf("lightning_executions_%s", productCode)}

	var paramCh = make(chan interface{})
	go api.doWebsocketRequest(childctx, *jsonRPC2, paramCh)

OUTER:
	for {
		select {
		case <-ctx.Done():
			return

		default:
			param := <-paramCh
			marshalExecutions, err := json.Marshal(param)
			if err != nil {
				continue OUTER
			}
			executions := make([]bitflyer.Execution, 0)
			if err := json.Unmarshal(marshalExecutions, &executions); err != nil {
				continue OUTER
			}
			ch <- executions
		}
	}
}
