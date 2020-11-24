package models

type DataFrameSet struct {
  DataFrames []DataFrame

  // inner data
  executionPool []Execution
  // candle duration is 1m
  candles []Candle
  volumes []float64

  m *sync.Mutex

  logger *log.Logger
}
