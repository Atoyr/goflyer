package models

import (
	"fmt"
	"time"
)

type Candles struct {
	productCode string
	duration    time.Duration
	candles     []Candle
}

func NewCandles(productCode string, duration time.Duration) Candles {
	var c Candles
	c.productCode = productCode
	c.duration = duration
	return c
}

func (cs *Candles) Add(datetime time.Time, price float64) {
	if length := len(cs.candles); length == 0 {
		c := NewCandle(cs.duration, datetime, price)
		cs.candles = append(cs.candles, c)
	} else {
		length--
		for i := range cs.candles {
			index := length - i
			if cs.candles[index].Time.Before(datetime) {
				c := NewCandle(cs.duration, datetime, price)

				if index == length {
					cs.candles = append(cs.candles,c)
				}else {
					temp := cs.candles[index + 1:]
					cs.candles = append(cs.candles[:index],c)
					cs.candles = append(cs.candles,temp...)
				}
				return 
			}else if cs.candles[index].Time.Equal(datetime) {
				cs.candles[index].Add(price)
        return 
			}
		}
		// append HEAD
		c := NewCandle(cs.duration, datetime, price)
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
