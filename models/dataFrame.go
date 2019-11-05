package models

import (
	"fmt"
	"time"

	"github.com/atoyr/go-talib"
)

// DataFrame is goflyer chart data framework
type DataFrame struct {
	productCode string
	duration    time.Duration

	datetimes []time.Time
	opens     []float64
	closes    []float64
	highs     []float64
	lows      []float64
	volumes   []float64

	Smas          []Sma
	Emas          []Ema
	BollingerBand *BollingerBand
	Rsis          []RelativeStrengthIndex
	Macd          []MovingAverageConvergenceDivergence
}

// NewDataFrame is getting CreateDataFrame
func NewDataFrame(productCode string, duration time.Duration) DataFrame {
	df := DataFrame{productCode: productCode, duration: duration}
	return df
}

// func JsonUnmarshalDataFrame(row []byte) (*DataFrame, error) {
// 	var dataFrame = new(DataFrame)
// 	err := json.Unmarshal(row, dataFrame)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return dataFrame, nil
// }

// Duration is duration getter
func (df *DataFrame) Duration() time.Duration {
	return df.duration
}

// Name is getting productCode_duration
func (df *DataFrame) Name() string {
	return fmt.Sprintf("%s_%s", df.productCode, df.duration)
}

// Add is Add value 
func (df *DataFrame) Add(datetime time.Time, price, volume float64) {
	dt := datetime.Truncate(df.duration)
	for i := range df.datetimes {
		index := len(df.datetimes) - i -1
		if df.datetimes[index].Equal(dt) {
			df.closes[index] = price
			if df.highs[index] < price {
				df.highs[index] = price
			}else if df.lows[index] > price {
				df.lows[index] = price
			}
			df.volumes[index] += volume
		}else if df.datetimes[index].Before(dt){
			if i == 0 {
				df.datetimes = append(df.datetimes, dt)
				df.opens = append(df.opens, price)
				df.closes = append(df.closes,price)
				df.highs = append(df.highs, price)
				df.lows  = append(df.lows, price ) 
				df.volumes  = append(df.volumes, volume) 
			}else {
				tdates := df.datetimes[index +1 :]
				df.datetimes , df.datetimes = append(df.datetimes[:index],dt), append(df.datetimes,tdates...)
				topens := df.opens[index +1 :]
				df.opens , df.opens = append(df.opens[:index],price), append(df.opens,topens...)
				tcloses := df.closes[index +1 :]
				df.closes , df.closes = append(df.closes[:index],price), append(df.closes,tcloses...)
				thighs := df.highs[index +1 :]
				df.highs , df.highs = append(df.highs[:index],price), append(df.highs,thighs...)
				tlows := df.lows[index +1 :]
				df.lows , df.lows = append(df.lows[:index],price), append(df.lows,tlows...)
				tvolumes := df.volumes[index +1 :]
				df.volumes , df.volumes = append(df.volumes[:index],price), append(df.volumes,tvolumes...)
			}
		}else if index == 0 {
			// append HEAD
			df.datetimes , df.datetimes[0] = append(df.datetimes[:1],df.datetimes[0:]...), dt
			df.opens , df.opens[0] = append(df.opens[:1],df.opens[0:]...), price
			df.closes , df.closes[0] = append(df.closes[:1],df.closes[0:]...), price
			df.highs , df.highs[0] = append(df.highs[:1],df.highs[0:]...), price
			df.lows , df.lows[0] = append(df.lows[:1],df.lows[0:]...), price 
			df.volumes , df.volumes[0] = append(df.volumes[:1],df.volumes[0:]...), volume
		}
	}
	if len(df.datetimes) == 0 {
		df.datetimes = append(df.datetimes, dt)
		df.opens = append(df.opens, price)
		df.closes = append(df.closes,price)
		df.highs = append(df.highs, price)
		df.lows  = append(df.lows, price ) 
		df.volumes  = append(df.volumes, volume) 
	}
	df.updateChart()
}

// GetCandle is Export Candles data
func (df *DataFrame) GetCandles() Candles {
	cs := NewCandles(df.productCode,df.duration)
	for i := range df.datetimes {
		c := Candle{Time: df.datetimes[i],Open: df.opens[i], Close: df.closes[i], High: df.highs[i],Low : df.lows[i]}
		cs.candles = append(cs.candles,c)
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
	sma := NewSma(df.closes, period)
	df.Smas = append(df.Smas, sma)
}

func (df *DataFrame) updateSmas() {
	for i := range df.Smas {
		df.Smas[i].Update(df.closes)
	}
}

func (df *DataFrame) refreshSmas() {
	for i, sma := range df.Smas {
		if len(df.datetimes) > sma.Period {
			df.Smas[i].Values = NewSma(df.closes, df.Smas[i].Period).Values
		} else {
			df.Smas[i].Values = make([]float64, len(df.datetimes))
		}
	}
}

// EMA

// AddEmas is Added Emas setting
func (df *DataFrame) AddEmas(period int) {
	ema := NewEma(df.closes, period)
	df.Emas = append(df.Emas, ema)
}

func (df *DataFrame) updateEmas() {
	for i := range df.Emas {
		df.Emas[i].Update(df.closes)
	}
}

func (df *DataFrame) refreshEmas() {
	for i, ema := range df.Emas {
		if len(df.datetimes) > ema.Period {
			df.Emas[i].Values = talib.Ema(df.closes, ema.Period)
		} else {
			df.Emas[i].Values = make([]float64, len(df.datetimes))
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
	if n <= len(df.datetimes) {
		closes := df.closes
		up1, center, down1 := talib.BBands(closes, n, k1, k1, 0)
		up2, center, down2 := talib.BBands(closes, n, k2, k2, 0)
		bb.Up2 = up2
		bb.Up1 = up1
		bb.Center = center
		bb.Down1 = down1
		bb.Down2 = down2
	} else {
		bb.Up2 = make([]float64, len(df.datetimes))
		bb.Up1 = make([]float64, len(df.datetimes))
		bb.Center = make([]float64, len(df.datetimes))
		bb.Down1 = make([]float64, len(df.datetimes))
		bb.Down2 = make([]float64, len(df.datetimes))
	}
	df.BollingerBand = bb
}

// RSI

// AddEsis is added Rsis setting
func (df *DataFrame) AddRsis(period int) {
	rsi := NewRelativeStrengthIndex(df.closes, period)
	df.Rsis = append(df.Rsis, rsi)
}

func (df *DataFrame) updateRsis() {
	for i := range df.Rsis {
		df.Rsis[i].Update(df.closes)
	}
}

func (df *DataFrame) refreshRsis() {
}

// MACD

// AddMacd is added Macd setting
func (df *DataFrame) AddMacd(fastPeriod, slowPeriod, signalPeriod int) {
	closes := df.closes
	macd := NewMovingAverageConvergenceDivergence(closes, fastPeriod, slowPeriod, signalPeriod)
	df.Macd = append(df.Macd, macd)
}

func (df *DataFrame) updateMacd() {
	for i := range df.Macd {
		df.Macd[i].Update(df.closes)
	}
}
