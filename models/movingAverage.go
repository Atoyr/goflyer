package models

type MovingAverage struct {
	Period      int
	Values      []float64
	periodTotal []float64
}

type Sma struct {
	MovingAverage
}

func NewSma(inReal []float64, inTimePeriod int) Sma {
	var sma Sma
	var values []float64
	var periodTotal []float64

	if len(inReal) < inTimePeriod {
		values := make([]float64, len(inReal))
		periodTotal := make([]float64, len(inReal))
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
	length := len(inReal) - len(sma.Values)
	if len(inReal) < sma.Period {
		values := make([]float64, length)
		sma.Values = append(sma.Values, values...)
		sma.periodTotal = append(sma.periodTotal, values...)
	} else {
		values = make([]float64, length)
		periodTotal := 0.0
		head := len(inReal) - sma.Period - length
		tail := len(inReal) - length
		for i := 0; i < tail; i++ {

		}

		for i := 0; i < start; i++ {
			total += inReal[i]
		}

		// TODO
		for i := 0; i < length; i++ {
			periodTotal += inReal[head]
			values[i] = periodTotal / float64(inTimePeriod)
			periodTotal -= inReal[tail]
		}

	}
}
