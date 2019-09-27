package models

type RelativeStrengthIndex struct {
	Period int       `json:"period"`
	Values []float64 `json:"values"`
	diff   []float64
}

func NewRelativeStrengthIndex(inReal []float64, inTimePeriod int) RelativeStrengthIndex {
	var rsi RelativeStrengthIndex
	if inTimePeriod < 1 {
		inTimePeriod = 1
	}
	rsi.Period = inTimePeriod
	rsi.setInitializeValues(0)
	rsi.Update(inReal)
	return rsi
}

func (rsi *RelativeStrengthIndex) Update(inReal []float64) {
	if difflength := len(inReal) - len(rsi.Values); difflength > 0 {
		inlength := len(inReal)
		valuelength := len(rsi.Values)
		rsi.appendValues(difflength)
		rsi.setDiffValue(inReal)

		if inlength >= rsi.Period {
			if valuelength < rsi.Period {
				rsi.setFirstPeriodValue(inReal)
			}
			for i := rsi.Period; i < inlength; i++ {
				a := rsi.Values[i-1] * float64(rsi.Period-1)
				b := rsi.Values[i-1] * float64(rsi.Period-1)
				if rsi.diff[i] > 0 {
					a = a + rsi.diff[i]
				} else {
					b = b - rsi.diff[i]
				}
				if a+b == 0 {
					rsi.Values[i] = 0
				} else {
					rsi.Values[i] = a / (a + b) * 100
				}
			}
		}
	}
}

func (rsi *RelativeStrengthIndex) setInitializeValues(length int) {
	rsi.Values = make([]float64, length)
	rsi.diff = make([]float64, length)
}

func (rsi *RelativeStrengthIndex) appendValues(length int) {
	rsi.Values = append(rsi.Values, make([]float64, length)...)
	rsi.diff = append(rsi.diff, make([]float64, length)...)
}

func (rsi *RelativeStrengthIndex) setDiffValue(inReal []float64) {
	if len(inReal) == 0 {
		return
	}
	beforeValue := 0.0
	for i := 1; i < len(inReal); i++ {
		if rsi.diff[i] == 0 {
			rsi.diff[i] = inReal[i] - beforeValue
		}
		beforeValue = inReal[i]
	}
}

func (rsi *RelativeStrengthIndex) setFirstPeriodValue(inReal []float64) {
	if len(inReal) < rsi.Period {
		return
	}
	increase := 0.0
	decrease := 0.0
	for i := 1; i < rsi.Period; i++ {
		if i < rsi.Period {
			if rsi.diff[i] > 0 {
				increase += rsi.diff[i]
			} else {
				decrease -= rsi.diff[i]
			}
		}
		if i == rsi.Period-1 && increase+decrease != 0 {
			rsi.Values[i] = increase / (increase + decrease) * 100
		}
	}
}
