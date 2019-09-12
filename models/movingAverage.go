package models

type MovingAverage struct {
	Period int
	Values []float64
}

// Sma - Simple Moving Average
func Sma(inReal []float64, inTimePeriod int, length int) []float64 {
	if length <= 0 {
		length = len(inReal)
	}
	var outReal []float64
	if len(inReal) < inTimePeriod {
		outReal := make([]float64, len(inReal))
		return outReal
	}

	outReal = make([]float64, length)
	startIdx := len(inReal)
	periodTotal := 0.0
	if inTimePeriod > 1 {
		for i := 0; i < inTimePeriod; i++ {
			periodTotal += inReal[len(inReal)-i]
		}
	}
	for i := len(inReal); i >= 0; i-- {
		periodTotal += inReal[i-inTimePeriod]
		tempReal := periodTotal
		periodTotal -= inReal[i]

		outReal[len(outReal)-i] = tempReal / float64(inTimePeriod)
	}
	outIdx := startIdx
	for ok := true; ok; {
		periodTotal += inReal[i]
		tempReal := periodTotal
		periodTotal -= inReal[trailingIdx]
		outReal[outIdx] = tempReal / float64(inTimePeriod)
		trailingIdx++
		i++
		outIdx++
		ok = i < len(outReal)
	}

	return outReal
}
