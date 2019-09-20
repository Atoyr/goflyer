package models

type RelativeStrengthIndex struct {
	Period int
	Values []float64
	diff   []float64
}

func NewRelativeStrengthIndex(inReal []float64, inTimePeriod int) RelativeStrengthIndex {
	var rsi RelativeStrengthIndex
	if inTimePeriod < 1 {
		inTimePeriod = 1
	}
	rsi.Period = inTimePeriod

	rsi.Values = make([]float64, len(inReal))
	rsi.diff = make([]float64, len(inReal))

	if len(inReal) > 0 {
		increase := 0.0
		decrease := 0.0
		beforeValue := inReal[0]
		for i := 1; i < len(inReal); i++ {
			rsi.diff[i] = inReal[i] - beforeValue
			if i < inTimePeriod {
				if rsi.diff[i] > 0 {
					increase += rsi.diff[i]
				} else {
					decrease -= rsi.diff[i]
				}
			}
			beforeValue = inReal[i]
		}
		if len(inReal) >= inTimePeriod {
			rsi.Values[inTimePeriod - 1] = increase / (increase + decrease) * 100
			for i :=inTimePeriod; i < len(inReal); i++ {
				a := rsi.Values[i -1 ] * float64(inTimePeriod - 1)
				b := rsi.Values[i -1 ] * float64(inTimePeriod - 1)
				if rsi.diff[i] > 0 {
					a += rsi.diff[i]
				}else {
					b -= rsi.diff[i] 
				}
				a = a / float64(inTimePeriod)
				b = b / float64(inTimePeriod)
				rsi.Values[i] = a / (a + b) * 100
			}
		}
	}
	return rsi
}

func (rsi *RelativeStrengthIndex) Update(inReal []float64) {
	if difflength := len(inReal) - len(rsi.Values); difflength > 0 {
		if inlength := len(inReal); inlength< rsi.Period {
			rsi.Values = make([]float64, len(inReal))
			rsi.diff = make([]float64, len(inReal)) 
		}else if inlength == rsi.Period {
			increase := 0.0
			decrease := 0.0
			rsilength := len(rsi.Values)
			rsi.Values = append(rsi.Values, make([]float64,difflength)...)
			rsi.diff = append(rsi.diff, make([]float64,difflength)...)
			beforeValue := 0.0
			if rsilength > 0 {
				beforeValue = inReal[rsilength - 1] 
			}
			if rsilength == 0 {
				rsilength = 1
			}
			for i := rsilength; i < len(inReal); i++ {
				rsi.diff[i] = inReal[i] - beforeValue
				beforeValue = inReal[i]
			}
			for i := 0; i < len(rsi.diff) ; i++ {
				if rsi.diff[i] > 0 {
					increase += rsi.diff[i]
				} else {
					decrease -= rsi.diff[i]
				}
			}
			rsi.Values[rsi.Period - 1] = increase / (increase + decrease) * 100 
		}else {
			// TODO
			if len(rsi.Values) > rsi.Period {

			}else {
				increase := 0.0
				decrease := 0.0
				rsilength := len(rsi.Values)
				rsi.Values = append(rsi.Values, make([]float64,difflength)...)
				rsi.diff = append(rsi.diff, make([]float64,difflength)...)
				beforeValue := 0.0
				if rsilength > 0 {
					beforeValue = inReal[rsilength - 1] 
				}
				if rsilength == 0 {
					rsilength = 1
				}
				for i := rsilength; i < len(inReal); i++ {
					rsi.diff[i] = inReal[i] - beforeValue
					beforeValue = inReal[i]
				}
				for i := 0; i < len(rsi.diff) ; i++ {
					if rsi.diff[i] > 0 {
						increase += rsi.diff[i]
					} else {
						decrease -= rsi.diff[i]
					}
				}
				rsi.Values[rsi.Period - 1] = increase / (increase + decrease) * 100 
			} 
		}
	}
}
