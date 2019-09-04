package models

import (
	"time"
	"fmt"
)

type CandleCollection struct {
	ProductCode string
	Duration time.Duration
	Candles map[string]Candle
}

func (c *CandleCollection) Name() string {
	return fmt.Sprintf("%s_%s",c.ProductCode, c.Duration)
}

func (c *CandleCollection) MergeCandle(candle Candle) {
	c.Candles[candle.Key()] = candle
}
