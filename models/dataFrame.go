package models

import (
	"fmt"
	"time"

	"github.com/atoyr/go-talib"
)

// DataFrame is goflyer chart data framework
type DataFrame struct {
	ProductCode string
	Duration    time.Duration

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
	Macd          []MovingAverageConvergenceDivergence
}

// NewDataFrame is getting CreateDataFrame
func NewDataFrame(productCode string, duration time.Duration) DataFrame {
	df := DataFrame{ProductCode: productCode, Duration: duration}
	return df
}

// Name is getting productCode_duration
func (df *DataFrame) Name() string {
	return fmt.Sprintf("%s_%s", df.ProductCode, df.Duration)
}

// Add is Add value
func (df *DataFrame) Add(datetime time.Time, price, volume float64) {
	dt := datetime.Truncate(df.Duration)
	for i := range df.Datetimes {
		index := len(df.Datetimes) - i - 1
		if df.Datetimes[index].Equal(dt) {
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
				df.Datetimes = append(df.Datetimes, dt)
				df.Opens = append(df.Opens, price)
				df.Closes = append(df.Closes, price)
				df.Highs = append(df.Highs, price)
				df.Lows = append(df.Lows, price)
				df.Volumes = append(df.Volumes, volume)
			} else {
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
	df.updateChart()
}

// GetCandle is Export Candles data
func (df *DataFrame) GetCandles() Candles {
	cs := NewCandles(df.ProductCode, df.Duration)
	for i := range df.Datetimes {
		c := Candle{Time: df.Datetimes[i], Open: df.Opens[i], Close: df.Closes[i], High: df.Highs[i], Low: df.Lows[i]}
		cs.Candles = append(cs.Candles, c)
	}
	return cs
}

func (df *DataFrame) updateChart() {
	df.updateSmas()
	df.updateEmas()
	df.updateMacd()
	df.updateRsis()
}

func (df *DataFrame) refreshChart() {
	df.refreshSmas()
	df.refreshEmas()
}

// SMA

// AddSmas is added Smas setting
func (df *DataFrame) AddSmas(period int) {
	sma := NewSma(df.Closes, period)
	df.Smas = append(df.Smas, sma)
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
	ema := NewEma(df.Closes, period)
	df.Emas = append(df.Emas, ema)
}

func (df *DataFrame) updateEmas() {
	for i := range df.Emas {
		df.Emas[i].Update(df.Closes)
	}
}

func (df *DataFrame) refreshEmas() {
	for i, ema := range df.Emas {
		if len(df.Datetimes) > ema.Period {
			df.Emas[i].Values = talib.Ema(df.Closes, ema.Period)
		} else {
			df.Emas[i].Values = make([]float64, len(df.Datetimes))
		}
	}
}

// BollingerBand

// AddBollingerBand is added setting
func (df *DataFrame) AddBollingerBand(n int, k1, k2 float64) {
	bb := new(BollingerBand)
	bb.N = n
	bb.K1 = k1
	bb.K2 = k2
	if n <= len(df.Datetimes) {
		Closes := df.Closes
		up1, center, down1 := talib.BBands(Closes, n, k1, k1, 0)
		up2, center, down2 := talib.BBands(Closes, n, k2, k2, 0)
		bb.Up2 = up2
		bb.Up1 = up1
		bb.Center = center
		bb.Down1 = down1
		bb.Down2 = down2
	} else {
		bb.Up2 = make([]float64, len(df.Datetimes))
		bb.Up1 = make([]float64, len(df.Datetimes))
		bb.Center = make([]float64, len(df.Datetimes))
		bb.Down1 = make([]float64, len(df.Datetimes))
		bb.Down2 = make([]float64, len(df.Datetimes))
	}
	df.BollingerBand = bb
}

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
	macd := NewMovingAverageConvergenceDivergence(Closes, fastPeriod, slowPeriod, signalPeriod)
	df.Macd = append(df.Macd, macd)
}

func (df *DataFrame) updateMacd() {
	for i := range df.Macd {
		df.Macd[i].Update(df.Closes)
	}
}
