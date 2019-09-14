package models

type MovingAverage struct {
	Period int
	Values []float64
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
	values := make([]float64, len(inReal))

	if len(inReal) >= inTimePeriod {
	} else {
		if inTimePeriod < 0 {
			inTimePeriod = 1
		}
		total := 0.0
		head := 0
		start := inTimePeriod - 1

		for i := 0; i < start; i++ {
			total += inReal[i]
		}

		for i := start; i < len(inReal); i++ {
			total += inReal[head]
			values[i] = total / float64(inTimePeriod)
			total -= inReal[i]
			head++
		}
	}

	sma.Period = inTimePeriod
	sma.Values = values
	return sma
}

// Sma - Simple Moving Average
func (sma *Sma) UupdateSma(inReal []float64) {
	var values []float64
	if difflength := len(inReal) - len(sma.Values); difflength > 0 {
		if len(inReal) < sma.Period {
			values := make([]float64, difflength)
			sma.Values = append(sma.Values, values...)
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
				periodTotal += inReal[tail]
				values[i] = periodTotal / float64(sma.Period)
				periodTotal -= inReal[head]
				head++
				tail++
			}
		}
		sma.Values = append(sma.Values, values...)
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
		k := 2 / float64(inTimePeriod+1)
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
