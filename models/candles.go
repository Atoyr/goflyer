package models

import (
	"fmt"
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

func (cs *Candles) Add(c Candle) {
	if length := len(cs.candles); length == 0 {
			cs.candles = append(cs.candles, c)
		} else {
			length--
			for i := range cs.candles {
				index := length - i
				if cs.candles[index].Time.Before(c.Time) {
					if index == length {
						cs.candles = append(cs.candles,c)
					}else {
						temp := cs.candles[index + 1:]
						cs.candles = append(cs.candles[:index],c)
						cs.candles = append(cs.candles,temp...)
					}
					return
				}
			}
			// append HEAD
			cs.candles, cs.candles[0] = append(cs.candles[:1], cs.candles[0:]...), c
		}
}

func (cs *Candles) Get(from,to int) (*Candles ,error ){
	ret := NewCandles(cs.productCode,cs.duration)
	if from > to {
		return &ret, fmt.Errorf("from value is large of to")
	}
	if len(cs.candles) == 0 {
		return &ret, nil
	}
	f, t := from,to

	if f < 0 {
		f = 0
	} else if f >= len(cs.candles) {
		f = len(cs.candles) -1
	}
	if t < 0 {
		t = 0
	} else if t >= len(cs.candles) {
		t = len(cs.candles) -1
	}
	ret.candles = cs.candles[f:t]
	return &ret, nil
}
