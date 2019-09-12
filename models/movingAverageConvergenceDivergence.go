package models

type MovingAverageConvergenceDivergence struct {
	FastPeriod int
	SlowPeriod int
	SignalPeriod int
	Macd []float64
	MacdSignal []float64
	MacdHist []float64
}
