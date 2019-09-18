package models

import "log"

type MovingAverageConvergenceDivergence struct {
	FastPeriod   int
	SlowPeriod   int
	SignalPeriod int
	Macd         []float64
	MacdSignal   []float64
	MacdHist     []float64
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
	if difflength := len(m.Macd) - len(inReal); difflength > 0 {
		m.fastEma.Update(inReal)
		m.slowEma.Update(inReal)
		macd := make([]float64, difflength)
		for i := range macd {
			macd[i] = m.fastEma.Values[difflength+i] - m.slowEma.Values[difflength+i]
		}
		m.Macd = append(m.Macd, macd...)
		m.signalEma.Update(m.Macd)
		m.MacdSignal = m.signalEma.Values
		hist := make([]float64, difflength)
		for i := range hist {
			hist[i] = m.Macd[difflength+1] - m.MacdSignal[difflength+i]
		}
		m.MacdHist = append(m.MacdHist, hist...)
		log.Printf("f:%d  s:%d signal:%d length macd : %d macdHist : %d macdSignal : %d", m.FastPeriod, m.SlowPeriod, m.SignalPeriod, len(m.Macd), len(m.MacdHist), len(m.MacdSignal))
	}
}
