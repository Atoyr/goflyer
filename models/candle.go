package models

import (
	"encoding/json"
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
	LastID        float64   `json:"last_id"`
}

func NewCandle(productCode string, duration time.Duration, time time.Time,id, price ,volume float64) *Candle {
	c := new(Candle)
	c.ProductCode = productCode
	c.Duration = duration.Nanoseconds()
	c.Time = time.Truncate(duration)
	c.Open = price
	c.Close = price
	c.High = price
	c.Low = price
	c.Volume = volume
	c.OpenDateTime = time
	c.CloseDateTime = time
	c.LastID = id
	return c
}

func JsonUnmarshalCandle(row []byte) (*Candle, error) {
	var candle = new(Candle)
	err := json.Unmarshal(row, candle)
	if err != nil {
		return nil, err
	}
	return candle, nil
}

func JsonUnmarshalCandles(row []byte) ([]Candle, error) {
	var candles []Candle
	err := json.Unmarshal(row, &candles)
	if err != nil {
		return nil, err
	}
	return candles, nil
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

func (c *Candle) Add(time time.Time, id,price,volume float64) error {
	toTime := c.Time.Add(c.GetDuration())
	if !c.Time.After(time) && toTime.Before(time) {
		if c.High < price {
			c.High = price
		}
		if c.Low > price {
			c.Low = price
		}
		c.Volume += volume
		if time.Before(c.OpenDateTime) {
			c.OpenDateTime = time
			c.Open = price
		}
		if c.LastID < id {
			c.CloseDateTime = time
			c.Close = price
			c.LastID = id
		}
		return nil
	} else {
		// TODO return error
		return nil
	}
}

func (c *Candle) GetCandleOHLC() CandleOHLC {
	ohlc := new(CandleOHLC)
	ohlc.Time = c.Time.Format(time.RFC3339)
	ohlc.Open = c.Open
	ohlc.High = c.High
	ohlc.Low = c.Low
	ohlc.Close = c.Close
	return *ohlc
}
