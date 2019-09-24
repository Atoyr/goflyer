package models

import (
	"time"
	"fmt"
)

type Candle struct {
	ProductCode string        `json:"product_code"`
	Duration    time.Duration `json:"duration"`
	Time        time.Time     `json:"time"`
	Open        float64       `json:"open"`
	Close       float64       `json:"close"`
	High        float64       `json:"high"`
	Low         float64       `json:"low"`
	Volume      float64       `json:"volume"`
}

func NewCandle(productCode string, duration time.Duration, time time.Time, open, close, high, low, volume float64) *Candle {
	c := new(Candle)
	c.ProductCode = productCode
	c.Duration = duration
	c.Time = time
	c.Open = open
	c.Close = close
	c.High = high
	c.Low = low
	c.Volume = volume
	return c
}

func (c *Candle) CollectionKey() string {
	return fmt.Sprintf("%s_%s",c.ProductCode, c.Duration)
}

func (c *Candle) Key() string {
	return fmt.Sprintf("%s_%s",c.ProductCode, c.GetTimeString())
}

func (c *Candle) GetTimeString() string {
	return c.Time.Format(time.RFC3339)
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
