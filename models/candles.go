package models

import (
	"fmt"
	"time"
)

type Candles struct {
	ProductCode string        `json:"product_code"`
	Duration    time.Duration `json:"duration"`
	Candles     []Candle      `json:"candles"`
}

func NewCandles(productCode string, duration time.Duration) Candles {
	var c Candles
	c.ProductCode = productCode
	c.Duration = duration
	return c
}

func (cs *Candles) Add(datetime time.Time, price float64) {
	if length := len(cs.Candles); length == 0 {
		c := NewCandle(cs.Duration, datetime, price)
		cs.Candles = append(cs.Candles, c)
	} else {
		length--
		for i := range cs.Candles {
			index := length - i
			if cs.Candles[index].Time.Before(datetime) {
				c := NewCandle(cs.Duration, datetime, price)

				if index == length {
					cs.Candles = append(cs.Candles,c)
				}else {
					temp := cs.Candles[index + 1:]
					cs.Candles = append(cs.Candles[:index],c)
					cs.Candles = append(cs.Candles,temp...)
				}
				return 
			}else if cs.Candles[index].Time.Equal(datetime) {
				cs.Candles[index].Add(price)
        return 
			}
		}
		// append HEAD
		c := NewCandle(cs.Duration, datetime, price)
		cs.Candles, cs.Candles[0] = append(cs.Candles[:1], cs.Candles[0:]...), c
	}
}

func (cs *Candles) Get(from,to int) (*Candles ,error ){
	ret := NewCandles(cs.ProductCode,cs.Duration)
	if from > to {
		return &ret, fmt.Errorf("from value is large of to")
	}
	if len(cs.Candles) == 0 {
		return &ret, nil
	}
	f, t := from,to

	if f < 0 {
		f = 0
	} else if f >= len(cs.Candles) {
		f = len(cs.Candles) -1
	}
	if t < 0 {
		t = 0
	} else if t >= len(cs.Candles) {
		t = len(cs.Candles) -1
	}
	ret.Candles = cs.Candles[f:t]
	return &ret, nil
}
