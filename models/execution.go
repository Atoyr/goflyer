package models

import (
	"encoding/json"
	"time"
)

// Execution
type Execution struct {
	Side  string    `json:"side"`
	Price float64   `json:"price"`
	Size  float64   `json:"size"`
	Time  time.Time `json:"exec_date"`
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
