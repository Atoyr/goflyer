package models

import (
	"encoding/json"
)

type MovingAverage struct {
	Period int       `json:"period"`
	Values []float64 `json:"values"`
}

type Sma struct {
	*MovingAverage
}

type Ema struct {
	*MovingAverage
}

func  JsonUnmarshalSma(row []byte)  (*Sma,error) {
	var sma = new(Sma)
	err := json.Unmarshal(row,sma)
	if err != nil {
		return nil, err
	}
	return sma ,nil
}
func  JsonUnmarshalEma(row []byte)  (*Ema,error) {
	var ema = new(Ema)
	err := json.Unmarshal(row,ema)
	if err != nil {
		return nil, err
	}
	return ema ,nil
}

func NewSma(inReal []float64, inTimePeriod int) Sma {
	var sma Sma
	sma.MovingAverage = new(MovingAverage)
	values := make([]float64, len(inReal))

	if inTimePeriod <= 0 {
		inTimePeriod = 1
	}

	if len(inReal) >= inTimePeriod {
		total := 0.0
		start := inTimePeriod - 1

		for i := 0; i < start; i++ {
			total += inReal[i]
		}

		for i := start; i < len(inReal); i++ {
			total += inReal[i]
			values[i] = total / float64(inTimePeriod)
			total -= inReal[i-inTimePeriod+1]
		}
	}

	sma.Period = inTimePeriod
	sma.Values = values
	return sma
}

// Sma - Simple Moving Average
func (sma *Sma) Update(inReal []float64) {
	if difflength := len(inReal) - len(sma.Values); difflength > 0 {
		values := make([]float64, difflength)
		if len(inReal) >= sma.Period {
			total := 0.0
			start := len(sma.Values) - 1
			if start < sma.Period-1 {
				start = sma.Period - 1
			}
			sma.Values = append(sma.Values, values...)

			for i := start - sma.Period + 1; i < start; i++ {
				total += inReal[i]
			}

			for i := start; i < len(inReal); i++ {
				total += inReal[i]
				sma.Values[i] = total / float64(sma.Period)
				total -= inReal[i-sma.Period+1]
			}
		} else {
			sma.Values = append(sma.Values, values...)
		}
	}
}

func NewEma(inReal []float64, inTimePeriod int) Ema {
	var ema Ema
	ema.MovingAverage = new(MovingAverage)
	values := make([]float64, len(inReal))

	if inTimePeriod < 1 {
		inTimePeriod = 1
	}

	if len(inReal) >= inTimePeriod {
		periodTotal := 0.0
		k := 2.0 / float64(inTimePeriod+1)
		for i := 0; i < inTimePeriod; i++ {
			periodTotal += inReal[i]
		}
		values[inTimePeriod-1] = periodTotal / float64(inTimePeriod)
		for i := inTimePeriod; i < len(inReal); i++ {
			values[i] = values[i-1] + k*(inReal[i]-values[i-1])
		}
	}

	ema.Period = inTimePeriod
	ema.Values = values
	return ema
}

func (ema *Ema) Update(inReal []float64) {
	if difflength := len(inReal) - len(ema.Values); difflength > 0 {
		if len(inReal) < ema.Period {
			values := make([]float64, difflength)
			ema.Values = append(ema.Values, values...)
		} else if len(inReal) == ema.Period {
			periodTotal := 0.0
			values := make([]float64, len(inReal))
			for i := 0; i < len(inReal); i++ {
				periodTotal += inReal[i]
			}
			values[len(inReal)-1] = periodTotal / float64(ema.Period)
			ema.Values = values
		} else {
			k := 2.0 / float64(ema.Period+1)

			for i := len(ema.Values); i < len(inReal); i++ {
				value := ema.Values[i-1] + k*(inReal[i]-ema.Values[i-1])
				ema.Values = append(ema.Values, value)
			}
		}
	}
}
