package models

import (
	"time"
  "sync"
  "sort"
  "log"
)

// DataFrame is goflyer chart data framework
type DataFrame struct {
	ProductCode string
	Duration    time.Duration

  // main data
	Datetimes []time.Time
	Opens     []float64
	Closes    []float64
	Highs     []float64
	Lows      []float64
	Volumes   []float64

	Smas          []Sma
	Emas          []Ema
	BollingerBand *BollingerBand
	Rsis          []RelativeStrengthIndex
	Macd          []MACD

  // inner data
  executionPool []Execution
  // candle duration is 1m
  candles []Candle
  volumes []float64

  m *sync.Mutex

  logger *log.Logger
}

// NewDataFrame is getting CreateDataFrame
func NewDataFrame(productCode string, duration time.Duration) DataFrame {
	df := DataFrame{ProductCode: productCode, Duration: duration}
  df.clearMainData()

  df.Smas = make([]Sma, 0)
  df.Emas = make([]Ema, 0)
  df.BollingerBand = new(BollingerBand)
  df.Rsis = make([]RelativeStrengthIndex, 0)
  df.Macd = make([]MACD, 0)

  df.executionPool = make([]Execution, 0)
  df.candles = make([]Candle, 0)
  df.volumes = make([]float64, 0)

  df.m = new(sync.Mutex)
	return df
}

func (df *DataFrame) AddExecution(e Execution) {
  df.m.Lock()
  df.executionPool = append(df.executionPool, e)
  df.m.Unlock()
}

func (df *DataFrame) ApplyExecution() {
  df.m.Lock()
  defer df.m.Unlock()
  sort.Slice(df.executionPool, func(i, j int) bool { return df.executionPool[i].Time.Before(df.executionPool[j].Time) })
  df.logf("apply execution : execution count is %d", len(df.executionPool))

  for i := range df.executionPool {
    df.Add(df.executionPool[i].Time, df.executionPool[i].Price, df.executionPool[i].Size)
    last := len(df.candles) - 1
    if last >= 0 && df.candles[last].Time.Equal(df.executionPool[i].Time.Truncate(1 * time.Minute)) {
      df.candles[last].Close = df.executionPool[i].Price
      if df.candles[last].High < df.executionPool[i].Price {
        df.candles[last].High = df.executionPool[i].Price
      }else if df.candles[last].Low > df.executionPool[i].Price {
        df.candles[last].Low = df.executionPool[i].Price
      }
      df.volumes[len(df.volumes) - 1] = df.volumes[len(df.volumes) - 1] + df.executionPool[i].Size
    }else {
      df.candles = append(df.candles, NewCandle(1 * time.Minute, df.executionPool[i].Time, df.executionPool[i].Price))
      df.volumes = append(df.volumes, df.executionPool[i].Size)
    }
  }
  df.executionPool = make([]Execution, 0)
}

// Add is Add value
func (df *DataFrame) Add(datetime time.Time, price, volume float64) {
	dt := datetime.Truncate(df.Duration)
	for i := range df.Datetimes {
		index := len(df.Datetimes) - i - 1
		if df.Datetimes[index].Equal(dt) {
      // add last
			df.Closes[index] = price
			if df.Highs[index] < price {
				df.Highs[index] = price
			} else if df.Lows[index] > price {
				df.Lows[index] = price
			}
			df.Volumes[index] += volume
			break
		} else if df.Datetimes[index].Before(dt) {
			if i == 0 {
        // append last
				df.Datetimes = append(df.Datetimes, dt)
				df.Opens = append(df.Opens, price)
				df.Closes = append(df.Closes, price)
				df.Highs = append(df.Highs, price)
				df.Lows = append(df.Lows, price)
				df.Volumes = append(df.Volumes, volume)
			} else {
        // add target
				tdates := df.Datetimes[index+1:]
				df.Datetimes, df.Datetimes = append(df.Datetimes[:index], dt), append(df.Datetimes, tdates...)
				tOpens := df.Opens[index+1:]
				df.Opens, df.Opens = append(df.Opens[:index], price), append(df.Opens, tOpens...)
				tCloses := df.Closes[index+1:]
				df.Closes, df.Closes = append(df.Closes[:index], price), append(df.Closes, tCloses...)
				tHighs := df.Highs[index+1:]
				df.Highs, df.Highs = append(df.Highs[:index], price), append(df.Highs, tHighs...)
				tLows := df.Lows[index+1:]
				df.Lows, df.Lows = append(df.Lows[:index], price), append(df.Lows, tLows...)
				tVolumes := df.Volumes[index+1:]
				df.Volumes, df.Volumes = append(df.Volumes[:index], price), append(df.Volumes, tVolumes...)
			}
			break
		} else if index == 0 {
			// append HEAD
			df.Datetimes, df.Datetimes[0] = append(df.Datetimes[:1], df.Datetimes[0:]...), dt
			df.Opens, df.Opens[0] = append(df.Opens[:1], df.Opens[0:]...), price
			df.Closes, df.Closes[0] = append(df.Closes[:1], df.Closes[0:]...), price
			df.Highs, df.Highs[0] = append(df.Highs[:1], df.Highs[0:]...), price
			df.Lows, df.Lows[0] = append(df.Lows[:1], df.Lows[0:]...), price
			df.Volumes, df.Volumes[0] = append(df.Volumes[:1], df.Volumes[0:]...), volume
			break
		}
	}
	if len(df.Datetimes) == 0 {
		df.Datetimes = append(df.Datetimes, dt)
		df.Opens = append(df.Opens, price)
		df.Closes = append(df.Closes, price)
		df.Highs = append(df.Highs, price)
		df.Lows = append(df.Lows, price)
		df.Volumes = append(df.Volumes, volume)
	}
}

// GetCandles is Export Candles data
func (df *DataFrame) GetCandles() []Candle {
  cs := make([]Candle,len(df.Datetimes))
	for i := range df.Datetimes {
		c := Candle{Time: df.Datetimes[i], Open: df.Opens[i], Close: df.Closes[i], High: df.Highs[i], Low: df.Lows[i]}
		cs[i] = c
	}
	return cs
}

func (df *DataFrame) SetLogger(l *log.Logger) {
  df.logger = l
}

func (df *DataFrame) logf(format string, v ...interface{}) {
  if df.logger == nil {
    return
  }
  df.logger.Printf(format, v...)
}

func (df *DataFrame) SetDuration(d time.Duration) {
  df.m.Lock()
  defer df.m.Unlock()
  df.Duration = d
  df.clearMainData()
  if len(df.candles) == 0 {
    return
  }

  for i := range df.candles {
    t := df.candles[i].Time.Truncate(df.Duration)
    last := len(df.Datetimes) - 1
    if last >= 0 && df.Datetimes[last].Equal(t) {
      if df.candles[i].High > df.Highs[last] {
        df.Highs[last] = df.candles[i].High
      }
      if df.candles[i].Low < df.Lows[last] {
        df.Lows[last] = df.candles[i].Low
      }
      df.Closes[last] = df.candles[i].Close
      df.Volumes[last] = df.Volumes[last] + df.volumes[i]
    }else {
      df.Datetimes = append(df.Datetimes, t)
      df.Opens = append(df.Opens, df.candles[i].Open)
      df.Highs = append(df.Highs, df.candles[i].High)
      df.Lows = append(df.Lows, df.candles[i].Low)
      df.Closes = append(df.Closes, df.candles[i].Close)
      df.Volumes = append(df.Volumes, df.volumes[i])
    }
  }
}

func (df *DataFrame) clearMainData() {
  df.Datetimes = make([]time.Time, 0)
  df.Opens = make([]float64, 0)
  df.Closes = make([]float64, 0)
  df.Highs = make([]float64, 0)
  df.Lows = make([]float64, 0)
  df.Volumes = make([]float64, 0)
}

func (df *DataFrame) updateChart() {
	df.updateSmas()
	df.updateEmas()
	df.updateMacd()
	df.updateRsis()
}

// SMA
// GetSmas is added Smas setting
func (df *DataFrame) GetSma(period int) Sma {
	for i := range df.Smas {
		if df.Smas[i].Period == period {
			return df.Smas[i]
		}
	}
	index := len(df.Smas)
	sma := NewSma(df.Closes, period)
	df.Smas = append(df.Smas, sma)
	return df.Smas[index]
}

func (df *DataFrame) updateSmas() {
	for i := range df.Smas {
		df.Smas[i].Update(df.Closes)
	}
}

func (df *DataFrame) refreshSmas() {
	for i, sma := range df.Smas {
		if len(df.Datetimes) > sma.Period {
			df.Smas[i].Values = NewSma(df.Closes, df.Smas[i].Period).Values
		} else {
			df.Smas[i].Values = make([]float64, len(df.Datetimes))
		}
	}
}

// EMA

// AddEmas is Added Emas setting
func (df *DataFrame) AddEmas(period int) {
	for i := range df.Emas {
		if df.Emas[i].Period == period {
			return
		}
	}
	ema := NewEma(df.Closes, period)
	df.Emas = append(df.Emas, ema)
}

func (df *DataFrame) updateEmas() {
	for i := range df.Emas {
		df.Emas[i].Update(df.Closes)
	}
}

// func (df *DataFrame) refreshEmas() {
// 	for i, ema := range df.Emas {
// 		if len(df.Datetimes) > ema.Period {
// 			df.Emas[i].Values = talib.Ema(df.Closes, ema.Period)
// 		} else {
// 			df.Emas[i].Values = make([]float64, len(df.Datetimes))
// 		}
// 	}
// }

// BollingerBand

// AddBollingerBand is added setting
// func (df *DataFrame) AddBollingerBand(n int, k1, k2 float64) {
// 	bb := new(BollingerBand)
// 	bb.N = n
// 	bb.K1 = k1
// 	bb.K2 = k2
// 	if n <= len(df.Datetimes) {
// 		Closes := df.Closes
// 		up1, center, down1 := talib.BBands(Closes, n, k1, k1, 0)
// 		up2, center, down2 := talib.BBands(Closes, n, k2, k2, 0)
// 		bb.Up2 = up2
// 		bb.Up1 = up1
// 		bb.Center = center
// 		bb.Down1 = down1
// 		bb.Down2 = down2
// 	} else {
// 		bb.Up2 = make([]float64, len(df.Datetimes))
// 		bb.Up1 = make([]float64, len(df.Datetimes))
// 		bb.Center = make([]float64, len(df.Datetimes))
// 		bb.Down1 = make([]float64, len(df.Datetimes))
// 		bb.Down2 = make([]float64, len(df.Datetimes))
// 	}
// 	df.BollingerBand = bb
// }

// RSI

// AddEsis is added Rsis setting
func (df *DataFrame) AddRsis(period int) {
	rsi := NewRelativeStrengthIndex(df.Closes, period)
	df.Rsis = append(df.Rsis, rsi)
}

func (df *DataFrame) updateRsis() {
	for i := range df.Rsis {
		df.Rsis[i].Update(df.Closes)
	}
}

func (df *DataFrame) refreshRsis() {
}

// MACD

// AddMacd is added Macd setting
func (df *DataFrame) AddMacd(fastPeriod, slowPeriod, signalPeriod int) {
	Closes := df.Closes
	macd := NewMACD(Closes, fastPeriod, slowPeriod, signalPeriod)
	df.Macd = append(df.Macd, macd)
}

func (df *DataFrame) updateMacd() {
	for i := range df.Macd {
		df.Macd[i].Update(df.Closes)
	}
}
