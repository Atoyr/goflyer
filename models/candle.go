package models

import (
	"time"
)

type Candle struct {
	ProductCode string
	Duration    time.Duration
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

func (c *Candle) AddTicker(ticker Ticker) (*Candle, error) {
	toTime := c.Time.Add(c.Duration)
	tickerTime := ticker.GetTimestamp()
	if !c.Time.After(tickerTime) && tickerTime.Before(toTime) {
		price := ticker.GetMidPrice()
		if c.High < price {
			c.High = price
		}
		if c.Low > price {
			c.Low = price
		}
		c.Volume += ticker.Volume
		c.Close = price
		return c, nil
	} else {
		// TODO return error
		return nil, nil
	}
}
