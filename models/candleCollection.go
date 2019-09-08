package models

import (
	"fmt"
	"github.com/atoyr/go-talib"
	"time"
)

type CandleCollection struct {
	ProductCode string
	Duration    time.Duration
	Candles     []Candle

	Smas          []MovingAverage
	Emas          []MovingAverage
	BollingerBand *BollingerBand
	// 	IchimokuCloud *IchimokuCloud `json:"ichimoku,omitempty"`
	// 	Rsi *Rsi `json:"rsi,omitempty"`
	// 	Macd *Macd `json:"macd,omitempty"`
	// 	Hvs []Hv `json:"hvs,omitempty"`
	// 	Events *SignalEvents `json:"events,omitempty"`
}

func (c *CandleCollection) Name() string {
	return fmt.Sprintf("%s_%s", c.ProductCode, c.Duration)
}

func (c *CandleCollection) AppendCnadle(candle Candle) {
	c.Candles = append(c.Candles, candle)
}

func (c *CandleCollection) Alls() (opens, closes, highs, lows, volumes []float64) {
	opens = make([]float64, len(c.Candles))
	closes = make([]float64, len(c.Candles))
	highs = make([]float64, len(c.Candles))
	lows = make([]float64, len(c.Candles))
	volumes = make([]float64, len(c.Candles))

	for i, v := range c.Candles {
		opens[i] = v.Open
		closes[i] = v.Close
		highs[i] = v.High
		lows[i] = v.Low
		volumes[i] = v.Volume
	}
	return
}

func (c *CandleCollection) Opens() []float64 {
	ret := make([]float64, len(c.Candles))
	for i, v := range c.Candles {
		ret[i] = v.Open
	}
	return ret
}

func (c *CandleCollection) Closes() []float64 {
	ret := make([]float64, len(c.Candles))
	for i, v := range c.Candles {
		ret[i] = v.Close
	}
	return ret
}

func (c *CandleCollection) Highs() []float64 {
	ret := make([]float64, len(c.Candles))
	for i, v := range c.Candles {
		ret[i] = v.High
	}
	return ret
}

func (c *CandleCollection) Lows() []float64 {
	ret := make([]float64, len(c.Candles))
	for i, v := range c.Candles {
		ret[i] = v.Low
	}
	return ret
}

func (c *CandleCollection) Volumes() []float64 {
	ret := make([]float64, len(c.Candles))
	for i, v := range c.Candles {
		ret[i] = v.Volume
	}
	return ret
}

func (c *CandleCollection) AddEmas(period int) {
	if len(c.Candles) > period {
		var ema MovingAverage
		ema.Period = period
		ema.Values = talib.Ema(c.Closes(), period)
	}
}

func (c *CandleCollection) AddBollingerBand(n int, k1, k2 float64) {
	if n <= len(c.Candles) {

		closes := c.Closes()
		up1, center, down1 := talib.BBands(closes, n, k1, k1, 0)
		up2, center, down2 := talib.BBands(closes, n, k2, k2, 0)
		bb := new(BollingerBand)
		bb.N = n
		bb.K1 = k1
		bb.K2 = k2
		bb.Up2 = up2
		bb.Up1 = up1
		bb.Center = center
		bb.Down1 = down1
		bb.Down2 = down2
		c.BollingerBand = bb
	}
}
