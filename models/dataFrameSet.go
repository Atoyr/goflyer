package models

import (
  "fmt"
	"time"
  "sync"
  "log"
  "sort"
)

type DataFrameSet struct {
	ProductCode string

  DataFrames []DataFrame

  // inner data
  executionPool []Execution
  // candle duration is 1m
  candles []Candle
  volumes []float64

  m *sync.Mutex

  logger *log.Logger
}

// NewDataFrameSet is getting DataFrameSet
func NewDataFrameSet(productCode string) DataFrameSet {
  dfs := DataFrameSet{ProductCode: productCode}

  dfs.executionPool = make([]Execution, 0)
  dfs.candles = make([]Candle, 0)
  dfs.volumes = make([]float64, 0)
  dfs.m = new(sync.Mutex)
	return dfs
}

func (dfs *DataFrameSet) GetDataFrame(duration time.Duration) (DataFrame, error) {
  for i := range dfs.DataFrames {
    if dfs.DataFrames[i].Duration == duration {
      return dfs.DataFrames[i], nil
    }
  }
  return *new(DataFrame), fmt.Errorf("DataFrame Not Found")
}

func (dfs *DataFrameSet) AddDataFrame(duration time.Duration) {
  if _, err := dfs.GetDataFrame(duration); err != nil {
    dfs.m.Lock()
    defer dfs.m.Unlock()

    df := NewDataFrame(dfs.ProductCode, duration)
    for i := range dfs.candles {
      t := dfs.candles[i].Time.Truncate(df.Duration)
      last := len(df.Datetimes) - 1
      if last >= 0 && df.Datetimes[last].Equal(t) {
        if dfs.candles[i].High > df.Highs[last] {
          df.Highs[last] = dfs.candles[i].High
        }
        if dfs.candles[i].Low < df.Lows[last] {
          df.Lows[last] = dfs.candles[i].Low
        }
        df.Closes[last] = dfs.candles[i].Close
        df.Volumes[last] = df.Volumes[last] + dfs.volumes[i]
      }else {
        df.Datetimes = append(df.Datetimes, t)
        df.Opens = append(df.Opens, dfs.candles[i].Open)
        df.Highs = append(df.Highs, dfs.candles[i].High)
        df.Lows = append(df.Lows, dfs.candles[i].Low)
        df.Closes = append(df.Closes, dfs.candles[i].Close)
        df.Volumes = append(df.Volumes, dfs.volumes[i])
      }
    }

    dfs.DataFrames = append(dfs.DataFrames, df)
  }
}

func (dfs *DataFrameSet) AddExecution(e Execution) {
  dfs.m.Lock()
  dfs.executionPool = append(dfs.executionPool, e)
  dfs.m.Unlock()
}

func (dfs *DataFrameSet) ApplyExecution() {
  dfs.m.Lock()
  defer dfs.m.Unlock()
  sort.Slice(dfs.executionPool, func(i, j int) bool { return dfs.executionPool[i].Time.Before(dfs.executionPool[j].Time) })

  for i := range dfs.executionPool {
    for j := range dfs.DataFrames {
      dfs.DataFrames[j].Add(dfs.executionPool[i].Time, dfs.executionPool[i].Price, dfs.executionPool[i].Size)

      // create one minite
      last := len(dfs.candles) - 1
      if last >= 0 && dfs.candles[last].Time.Equal(dfs.executionPool[i].Time.Truncate(1 * time.Minute)) {
        dfs.candles[last].Close = dfs.executionPool[i].Price
        if dfs.candles[last].High < dfs.executionPool[i].Price {
          dfs.candles[last].High = dfs.executionPool[i].Price
        }else if dfs.candles[last].Low > dfs.executionPool[i].Price {
          dfs.candles[last].Low = dfs.executionPool[i].Price
        }
        dfs.volumes[len(dfs.volumes) - 1] = dfs.volumes[len(dfs.volumes) - 1] + dfs.executionPool[i].Size
      }else {
        dfs.candles = append(dfs.candles, NewCandle(1 * time.Minute, dfs.executionPool[i].Time, dfs.executionPool[i].Price))
        dfs.volumes = append(dfs.volumes, dfs.executionPool[i].Size)
      }
    }
  }
  dfs.executionPool = make([]Execution, 0)
}
