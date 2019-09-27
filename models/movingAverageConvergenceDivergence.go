package models

type MovingAverageConvergenceDivergence struct {
	FastPeriod   int       `json:"fast_period"`
	SlowPeriod   int       `json:"slow_period"`
	SignalPeriod int       `json:"signal_period"`
	Macd         []float64 `json:"macd"`
	MacdSignal   []float64 `json:"macd_signal"`
	MacdHist     []float64 `json:"macd_hist"`
	fastEma      Ema
	slowEma      Ema
	signalEma    Ema
}

func NewMovingAverageConvergenceDivergence(inReal []float64, fastPeriod, slowPeriod, signalPeriod int) MovingAverageConvergenceDivergence {
	var movingAverageConvergenceDivergence MovingAverageConvergenceDivergence
	fastEma := NewEma([]float64{}, fastPeriod)
	slowEma := NewEma([]float64{}, slowPeriod)
	signalEma := NewEma([]float64{}, signalPeriod)
	macd := make([]float64, 0)
	hist := make([]float64, 0)

	movingAverageConvergenceDivergence.FastPeriod = fastPeriod
	movingAverageConvergenceDivergence.SlowPeriod = slowPeriod
	movingAverageConvergenceDivergence.SignalPeriod = signalPeriod
	movingAverageConvergenceDivergence.Macd = macd
	movingAverageConvergenceDivergence.MacdSignal = signalEma.Values
	movingAverageConvergenceDivergence.MacdHist = hist
	movingAverageConvergenceDivergence.fastEma = fastEma
	movingAverageConvergenceDivergence.slowEma = slowEma
	movingAverageConvergenceDivergence.signalEma = signalEma
	movingAverageConvergenceDivergence.Update(inReal)

	return movingAverageConvergenceDivergence
}

func (m *MovingAverageConvergenceDivergence) Update(inReal []float64) {
	if difflength := len(inReal) - len(m.Macd); difflength > 0 {
		baselength := len(m.Macd)
		m.fastEma.Update(inReal)
		m.slowEma.Update(inReal)
		macd := make([]float64, difflength)
		for i := range macd {
			macd[i] = m.fastEma.Values[baselength+i] - m.slowEma.Values[baselength+i]
		}
		m.Macd = append(m.Macd, macd...)
		m.signalEma.Update(m.Macd)
		m.MacdSignal = m.signalEma.Values
		hist := make([]float64, difflength)
		for i := range hist {
			hist[i] = m.Macd[baselength+i] - m.MacdSignal[baselength+i]
		}
		m.MacdHist = append(m.MacdHist, hist...)
	}
}
