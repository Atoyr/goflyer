package models

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
	fastEma := NewEma(inReal, fastPeriod)
	slowEma := NewEma(inReal, slowPeriod)
	macd := make([]float64, len(fastEma.Values))
	for i := range macd {
		macd[i] = fastEma.Values[i] - slowEma.Values[i]
	}
	signalEma := NewEma(macd, signalPeriod)
	hist := make([]float64, len(signalEma.Values))
	for i := range hist {
		hist[i] = macd[i] - signalEma.Values[i]
	}

	movingAverageConvergenceDivergence.FastPeriod = fastPeriod
	movingAverageConvergenceDivergence.SlowPeriod = slowPeriod
	movingAverageConvergenceDivergence.SignalPeriod = signalPeriod
	movingAverageConvergenceDivergence.Macd = macd
	movingAverageConvergenceDivergence.MacdSignal = signalEma.Values
	movingAverageConvergenceDivergence.MacdHist = hist
	movingAverageConvergenceDivergence.fastEma = fastEma
	movingAverageConvergenceDivergence.slowEma = slowEma
	movingAverageConvergenceDivergence.signalEma = signalEma

	return movingAverageConvergenceDivergence
}
