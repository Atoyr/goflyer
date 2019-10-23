package models

import (
	"time"
)

type Candles struct {
	productCode string
	duration    int64
	candles     []Candle
}

func NewCandles(productCode string, duration int64) Candles {
	var c Candles
	c.productCode = productCode
	c.duration = duration
	return c
}

func (cs Candles) Add(datetime time.Time, id, price, volume float64) {
	if c, index, ok := cs.whereCandle(datetime); ok {
		c.Add(datetime, id, price, volume)
	} else {
		c := NewCandle(cs.productCode, time.Duration(cs.duration), datetime, id, price, volume)
		if index == 0 {
			// append HEAD
			if cap(cs.candles) == 0 {
				cs.candles = append(cs.candles, *c)
			} else {
				cs.candles, cs.candles[0] = append(cs.candles[:1], cs.candles[0:]...), *c
			}
		} else if index == len(cs.candles) {
			// append TAIL
			cs.candles = append(cs.candles, *c)
		} else {
			// insert
			temp := cs.candles[index:]
			cs.candles = append(cs.candles[:index-1], *c)
			cs.candles = append(cs.candles, temp...)
		}
	}
}

func (cs *Candles) Candles() []Candle {
	return cs.candles
}

func (cs Candles) GetCandleOHLCs() []CandleOHLC {
	ohlcs := make([]CandleOHLC, len(cs.candles))
	for i := range cs.candles {
		ohlcs[i] = cs.candles[i].GetCandleOHLC()
	}

	return ohlcs
}

func (cs Candles) AppendCandle(candles ...Candle) {
	cs.candles = append(cs.candles, candles...)
}

func (cs *Candles) Len() int {
	return len(cs.candles)
}

func (cs *Candles) whereCandle(datetime time.Time) (candle *Candle, index int, ok bool) {
	truncateTime := datetime.Truncate(time.Duration(cs.duration))
	// no find and truncateTime position is tail
	if cs.candles[len(cs.candles)-1].Time.Before(truncateTime) {
		return nil, len(cs.candles), false
	}
	for i := 0; i < len(cs.candles); i++ {
		index := len(cs.candles) - i
		if cs.candles[index].Time.Equal(truncateTime) {
			return &cs.candles[index], index, true
		}
		if index > 0 {
			if cs.candles[index].Time.After(truncateTime) {
				return nil, index, false
			}
		} else {
			return nil, index, false
		}
	}
	return nil, len(cs.candles), false
}
