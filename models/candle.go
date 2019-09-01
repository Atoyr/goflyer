package models

import (
	"time"
)

type Candle struct {
	ProductCode string
	Time        time.Time
	Open        float64
	Close       float64
	High        float64
	Low         float64
	Volume      float64
}

func NewCandle(productCode string, time time.Time, open, close, high, low, volume float64) *Candle {
	c := new(Candle)
	c.ProductCode = productCode
	c.Time = time
	c.Open = open
	c.Close = close
	c.High = high
	c.Low = low
	c.Volume = volume
	return c
}
