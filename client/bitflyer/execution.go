package bitflyer

import (
	"encoding/json"
	"log"
	"time"
)

type Execution struct {
	ID                         float64 `json:"id"`
	Side                       string  `json:"side"`
	Price                      float64 `json:"price"`
	Size                       float64 `json:"size"`
	ExecDate                   string  `json:"exec_date"`
	BuyChildOrderAcceptanceID  string  `json:"buy_child_order_acceptance_id"`
	SellChildOrderAcceptanceID string  `json:"sell_child_order_acceptance_id"`
}

func JsonUnmarshalExecution(row []byte) (*Execution, error) {
	var execution = new(Execution)
	err := json.Unmarshal(row, execution)
	if err != nil {
		return nil, err
	}
	return execution, nil
}

func JsonUnmarshalExecutions(row []byte) ([]Execution, error) {
	var executions []Execution
	err := json.Unmarshal(row, &executions)
	if err != nil {
		return nil, err
	}
	return executions, nil
}

func (e *Execution) DateTime() time.Time {
	datetime, err := time.Parse(time.RFC3339, e.ExecDate)
	if err != nil {
		log.Printf("action=Execution.GetTimestamp, argslen=0, args=, err=%s", err.Error())
		return time.Now()
	}
	return datetime
}
