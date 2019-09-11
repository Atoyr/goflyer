package models

import (
	"fmt"
	"github.com/atoyr/go-talib"
	"time"
)

type CandleCollections map[string]CandleCollection

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

const Open = "Open"
const Close = "Close"
const High = "High"
const Low = "Low"
const Volume = "Volume"

func NewCandleCollection(productCode string, duration time.Duration) CandleCollection {
	cc := CandleCollection{ProductCode: productCode, Duration: duration}
	return cc
}

func (c *CandleCollection) Name() string {
	return fmt.Sprintf("%s_%s", c.ProductCode, c.Duration)
}

func (c *CandleCollection) MergeCandle(candle Candle) error {
	if candle.Duration != c.Duration {
		// TODO return error
		return nil
	}
	if len(c.Candles) == 0 {
		c.Candles = []Candle{candle}
		c.updateChart()
		return nil
	}

	max := len(c.Candles) - 1
	beforeTime := c.Candles[max].Time
	if candle.Time.Equal(beforeTime) {
		c.Candles[max] = candle
		c.refreshChart()
	} else if candle.Time.Before(c.Candles[max].Time) {
		for i := range c.Candles {
			if candle.Time.Equal(c.Candles[max-i].Time) {
				c.Candles[len(c.Candles)-1-i] = candle
				c.refreshChart()
				break
			} else if candle.Time.Before(beforeTime) && candle.Time.After(c.Candles[max-i].Time) {
				before := c.Candles[:max-i]
				after := c.Candles[max-i+1:]
				c.Candles = append(before, candle)
				c.Candles = append(c.Candles, after...)
				c.refreshChart()
				break
			}
		}
	} else {
		c.Candles = append(c.Candles, candle)
		c.updateChart()
	}
	return nil
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

func (c *CandleCollection) Values(valueType string) []float64 {
	ret, _ := c.LastOfValues(valueType, 0)
	return ret
}

func (c *CandleCollection) LastOfValues(valueType string, from int) ([]float64, error) {
	if len(c.Candles) <= from {
		// TODO return error
		return nil, nil
	}
	// 123456 012345
	ret := make([]float64, len(c.Candles)-from)
	switch valueType {
	case Open:
		for i, v := range c.Candles[from:] {
			ret[i] = v.Open
		}
	case Close:
		for i, v := range c.Candles[from:] {
			ret[i] = v.Close
		}
	case High:
		for i, v := range c.Candles[from:] {
			ret[i] = v.High
		}
	case Low:
		for i, v := range c.Candles[from:] {
			ret[i] = v.Low
		}
	case Volume:
		for i, v := range c.Candles[from:] {
			ret[i] = v.Volume
		}
	default:
	}
	return ret, nil
}

func (c *CandleCollection) updateChart() {
	c.updateSmas()
	c.updateEmas()
}

func (c *CandleCollection) refreshChart() {
	c.refreshSmas()
	c.refreshEmas()
}

// SMA
func (c *CandleCollection) AddSmas(period int) {
	var sma MovingAverage
	sma.Period = period
	if len(c.Candles) > period {
		sma.Values = talib.Sma(c.Values(Close), period)
	} else {
		sma.Values = make([]float64, len(c.Candles))
	}
	c.Smas = append(c.Smas, sma)
}

func (c *CandleCollection) updateSmas() {
	for i, sma := range c.Smas {
		appendlength := len(c.Candles) - len(sma.Values)
		if appendlength > 0 {
			if len(c.Candles) > sma.Period {
				length := appendlength + sma.Period
				if length > len(c.Candles) {
					length = len(c.Candles)
				}
				candles, _ := c.LastOfValues(Close, len(c.Candles)-length)
				values := talib.Sma(candles, sma.Period)
				c.Smas[i].Values = append(c.Smas[i].Values, values[len(values)-appendlength:]...)
			} else {
				sma.Values = make([]float64, len(c.Candles))
			}
		}
	}
}

func (c *CandleCollection) refreshSmas() {
	for i, sma := range c.Smas {
		if len(c.Candles) > sma.Period {
			c.Smas[i].Values = talib.Sma(c.Values(Close), sma.Period)
		} else {
			c.Smas[i].Values = make([]float64, len(c.Candles))
		}
	}
}

// EMA
func (c *CandleCollection) AddEmas(period int) {
	var ema MovingAverage
	ema.Period = period
	if len(c.Candles) > period {
		ema.Values = talib.Ema(c.Values(Close), period)
	} else {
		ema.Values = make([]float64, len(c.Candles))
	}
	c.Emas = append(c.Emas, ema)
}

func (c *CandleCollection) updateEmas() {
	for i, ema := range c.Emas {
		appendlength := len(c.Candles) - len(ema.Values)
		if appendlength > 0 {
			if len(c.Candles) > ema.Period {
				length := appendlength + ema.Period
				if length > len(c.Candles) {
					length = len(c.Candles)
				}
				candles, _ := c.LastOfValues(Close, len(c.Candles)-length)
				values := talib.Ema(candles, ema.Period)
				c.Emas[i].Values = append(c.Emas[i].Values, values[len(values)-appendlength:]...)
			} else {
				ema.Values = make([]float64, len(c.Candles))
			}
		}
	}
}

func (c *CandleCollection) refreshEmas() {
	for i, ema := range c.Emas {
		if len(c.Candles) > ema.Period {
			c.Emas[i].Values = talib.Ema(c.Values(Close), ema.Period)
		} else {
			c.Emas[i].Values = make([]float64, len(c.Candles))
		}
	}
}

// BollingerBand
func (c *CandleCollection) AddBollingerBand(n int, k1, k2 float64) {
	bb := new(BollingerBand)
	bb.N = n
	bb.K1 = k1
	bb.K2 = k2
	if n <= len(c.Candles) {
		closes := c.Values(Close)
		up1, center, down1 := talib.BBands(closes, n, k1, k1, 0)
		up2, center, down2 := talib.BBands(closes, n, k2, k2, 0)
		bb.Up2 = up2
		bb.Up1 = up1
		bb.Center = center
		bb.Down1 = down1
		bb.Down2 = down2
	} else {
		bb.Up2 = make([]float64, len(c.Candles))
		bb.Up1 = make([]float64, len(c.Candles))
		bb.Center = make([]float64, len(c.Candles))
		bb.Down1 = make([]float64, len(c.Candles))
		bb.Down2 = make([]float64, len(c.Candles))
	}
	c.BollingerBand = bb
}
