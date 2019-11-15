package models

import "encoding/json"

type Log struct {
	LogLevel int    `json:"log_level"`
	Message  string `json:"message"`
}

func JsonUnmarshalLog(row []byte) (*Log, error) {
	var log = new(Log)
	err := json.Unmarshal(row, log)
	if err != nil {
		return nil, err
	}
	return log, nil
}
