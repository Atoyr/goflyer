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

func (c *CandleCollection) MergeCandle(candle Candle) {
	key := fmt.Sprintf("%s_%s",c.ProductCode, c.Duration)
	c.Candles[key] = candle
}
