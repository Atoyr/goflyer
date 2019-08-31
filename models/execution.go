package models

import (
	"log"
	"time"
)

type Execution struct {
	ID                         int     `json:"id"`
	Side                       string  `json:"side"`
	Price                      int64   `json:"price"`
	Size                       float64 `json:"size"`
	ExecDate                   string  `json:"exec_date"`
	BuyChildOrderAcceptanceID  string  `json:"buy_child_order_acceptance_id"`
	SellChildOrderAcceptanceID string  `json:"sell_child_order_acceptance_id"`
}

func (e *Execution) GetExecDate() time.Time {
	datetime, err := time.Parse(time.RFC3339, e.ExecDate)
	if err != nil {
		log.Printf("action=Execution.GetExecDate, argslen=0, args=, err=%s", err.Error())
	}
	return datetime
}
