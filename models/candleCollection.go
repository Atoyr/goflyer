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
	TimeValue   []string

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

func (c *CandleCollection) Name() string {
	return fmt.Sprintf("%s_%s", c.ProductCode, c.Duration)
}

func (c *CandleCollection) AppendCnadle(candle Candle) {
	c.Candles = append(c.Candles, candle)
	c.TimeValue = append(c.TimeValue, candle.GetTimeString())
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
func (c *CandleCollection) AddSmas(period int) {
	if len(c.Candles) > period {
		var sma MovingAverage
		sma.Period = period
		sma.Values = talib.Sma(c.Values(Close), period)
	}
}

func (c *CandleCollection) updateSmas() error {
	for _, sma := range c.Smas {
		err := c.updateSma(sma)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CandleCollection) updateSma(sma MovingAverage) error {
	length := len(c.Candles) - len(sma.Values) + sma.Period
	candles, err := c.LastOfValues(Close, length)
	if err != nil {
		return err
	}
	sma.Values = append(sma.Values, talib.Sma(candles, sma.Period)...)
	return nil
}

func (c *CandleCollection) AddEmas(period int) {
	if len(c.Candles) > period {
		var ema MovingAverage
		ema.Period = period
		ema.Values = talib.Ema(c.Values(Close), period)
		c.Emas = append(c.Emas, ema)
	}
}

func (c *CandleCollection) updateEmas() error {
	for _, ema := range c.Emas {
		err := c.updateEma(ema)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CandleCollection) updateEma(ema MovingAverage) error {
	length := len(c.Candles) - len(ema.Values) + ema.Period
	candles, err := c.LastOfValues(Close, length)
	if err != nil {
		return err
	}
	ema.Values = append(ema.Values, talib.Ema(candles, ema.Period)...)
	return nil
}

func (c *CandleCollection) AddBollingerBand(n int, k1, k2 float64) {
	if n <= len(c.Candles) {

		closes := c.Values(Close)
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
