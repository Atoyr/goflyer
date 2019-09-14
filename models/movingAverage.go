package models

import (
	"fmt"
	"log"
)

type MovingAverage struct {
	Period      int
	Values      []float64
	periodTotal []float64
}

type Sma struct {
	*MovingAverage
}

type Ema struct {
	*MovingAverage
}

func NewSma(inReal []float64, inTimePeriod int) Sma {
	var sma Sma
	sma.MovingAverage = new(MovingAverage)
	var values []float64
	var periodTotal []float64

	if len(inReal) < inTimePeriod {
		values = make([]float64, len(inReal))
		periodTotal = make([]float64, len(inReal))
	} else {
		if inTimePeriod < 0 {
			inTimePeriod = 1
		}
		values = make([]float64, inTimePeriod)
		total := 0.0
		head := 0
		start := inTimePeriod - 1

		for i := 0; i < start; i++ {
			total += inReal[i]
		}

		for i := start; i < len(inReal); i++ {
			total += inReal[head]
			values[i] = total / float64(inTimePeriod)
			periodTotal[i] = total
			total -= inReal[i]
			head++
		}
	}

	sma.Period = inTimePeriod
	sma.Values = values
	sma.periodTotal = periodTotal
	return sma
}

// Sma - Simple Moving Average
func (sma *Sma) UupdateSma(inReal []float64) {
	var values []float64
	if difflength := len(inReal) - len(sma.Values); difflength > 0 {
		if len(inReal) < sma.Period {
			values := make([]float64, difflength)
			sma.Values = append(sma.Values, values...)
			sma.periodTotal = append(sma.periodTotal, values...)
		} else {
			values = make([]float64, difflength)
			periodTotal := 0.0
			tail := len(sma.Values)
			head := tail - sma.Period + 1
			if head < 0 {
				difflength = difflength + head
				tail = tail - head
				head = 0
				if tail > len(inReal) {
					return
				}
			}

			for i := head; i < tail; i++ {
				periodTotal += inReal[i]
			}

			for i := 0; i < difflength; i++ {
				log.Printf("in : %d, values : %d, period : %d,  head : %d, tail : %d, diff : %d", len(inReal), len(sma.Values), sma.Period, head, tail, difflength)
				periodTotal += inReal[tail]
				values[i] = periodTotal / float64(sma.Period)
				periodTotal -= inReal[head]
				head++
				tail++
			}
		}
		var s string
		for _, v := range values {
			s = s + fmt.Sprintf(", %f ", v)
		}
		log.Println(s)
		sma.Values = append(sma.Values, values...)
	}
}
