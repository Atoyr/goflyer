package models

import (
	"fmt"
	"time"
)

type Candle struct {
	ProductCode   string    `json:"product_code"`
	Duration      int64     `json:"duration"`
	Time          time.Time `json:"time"`
	Open          float64   `json:"open"`
	Close         float64   `json:"close"`
	High          float64   `json:"high"`
	Low           float64   `json:"low"`
	Volume        float64   `json:"volume"`
	OpenDateTime  time.Time `json:"open_date_time"`
	CloseDateTime time.Time `json:"close_date_time"`
}

type Candles []Candle

func NewCandle(productCode string, duration time.Duration, ticker Ticker) *Candle {
	c := new(Candle)
	c.ProductCode = productCode
	c.Duration = duration.Nanoseconds()
	c.Time = ticker.DateTime()
	c.Open = ticker.GetMidPrice()
	c.Close = ticker.GetMidPrice()
	c.High = ticker.GetMidPrice()
	c.Low = ticker.GetMidPrice()
	c.Volume = ticker.Volume
	c.OpenDateTime = ticker.DateTime()
	c.CloseDateTime = ticker.DateTime()
	return c
}

func (c *Candle) CollectionKey() string {
	return fmt.Sprintf("%s_%s", c.ProductCode, c.GetDuration())
}

func (c *Candle) Key() string {
	return fmt.Sprintf("%s_%s", c.ProductCode, c.GetTimeString())
}

func (c *Candle) GetTimeString() string {
	return c.Time.Format(time.RFC3339)
}

func (c *Candle) GetDuration() time.Duration {
	return time.Duration(c.Duration)
}

func (c *Candle) AddTicker(ticker Ticker) error {
	toTime := c.Time.Add(c.GetDuration())
	tickerTime := ticker.DateTime()
	if !c.Time.After(tickerTime) && tickerTime.Before(toTime) {
		price := ticker.GetMidPrice()
		if c.High < price {
			c.High = price
		}
		if c.Low > price {
			c.Low = price
		}
		c.Volume += ticker.Volume
		if tickerTime.Before(c.OpenDateTime) {
			c.OpenDateTime = tickerTime
		}
		if tickerTime.After(c.CloseDateTime) {
			c.CloseDateTime = tickerTime
		}
		return nil
	} else {
		// TODO return error
		return nil
	}
}
