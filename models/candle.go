package models

import (
	"encoding/json"
	"time"
)

type Candle struct {
	Time          time.Time `json:"time"`
	Open          float64   `json:"open"`
	Close         float64   `json:"close"`
	High          float64   `json:"high"`
	Low           float64   `json:"low"`
}

func NewCandle(duration time.Duration, time time.Time, price float64) Candle {
	c := new(Candle)
	c.Time = time.Truncate(duration)
	c.Open = price
	c.Close = price
	c.High = price
	c.Low = price
	return *c
}

func JsonUnmarshalCandle(row []byte) (*Candle, error) {
	var candle = new(Candle)
	err := json.Unmarshal(row, candle)
	if err != nil {
		return nil, err
	}
	return candle, nil
}

func (c *Candle) GetTimeString() string {
	return c.Time.Format(time.RFC3339)
}

func (c *Candle) Add(price float64)  {
	if c.High < price {
		c.High = price
	}
	if c.Low > price {
		c.Low = price
	}
	c.Close = price
}
