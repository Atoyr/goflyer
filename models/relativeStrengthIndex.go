package models

type RelativeStrengthIndex struct {
	Period int
	Values []float64
	diff   []float64
}

func NewRelativeStrengthIndex(inReal []float64, inTimePeriod int) RelativeStrengthIndex {
	var rsi RelativeStrengthIndex
	rsi.Period = inTimePeriod

	rsi.Values = make([]float64, len(inReal))
	rsi.diff = make([]float64, len(inReal))

	if len(inReal) > 0 {
		beforeValue := inReal[0]
		for i := 1; i < len(inReal); i++ {
			rsi.diff[i] = inReal[i] - beforeValue
			beforeValue = inReal[i]
		}
		// TODO
		increase := 0.0
		decrease := 0.0
		for i := inTimePeriod - 1; i < len(inReal); i++ {
			if rsi.diff[i] > 0 {
				increase += rsi.diff[i]
			} else {
				decrease -= rsi.diff[i]
			}
		}
	}
	return rsi
}
