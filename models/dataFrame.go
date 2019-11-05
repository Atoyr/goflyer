package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/atoyr/go-talib"
)

type DataFrame struct {
	productCode string
	duration    time.Duration
	candles     Candles

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

func NewDataFrame(productCode string, duration time.Duration) DataFrame {
	df := DataFrame{productCode: productCode, duration: duration}
	df.candles = NewCandles(productCode, duration)

	return df
}

func JsonUnmarshalDataFrame(row []byte) (*DataFrame, error) {
	var dataFrame = new(DataFrame)
	err := json.Unmarshal(row, dataFrame)
	if err != nil {
		return nil, err
	}
	return dataFrame, nil
}

func (df *DataFrame) Duration() time.Duration {
	return df.duration
}

func (df *DataFrame) Name() string {
	return fmt.Sprintf("%s_%s", df.productCode, df.duration)
}

func (df *DataFrame) AddValue(datetime time.Time, price, volume float64) {
	df.candles.Add(datetime,price)
	// TODO UPDATE open ~ close data
	// TODO UPDATE volumes
	df.updateChart()
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
		if len(df.candles.candles) > sma.Period {
			df.Smas[i].Values = NewSma(df.closes, df.Smas[i].Period).Values
		} else {
			df.Smas[i].Values = make([]float64, len(df.candles.candles))
		}
	}
}

// EMA
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
		if len(df.candles.candles) > ema.Period {
			df.Emas[i].Values = talib.Ema(df.closes, ema.Period)
		} else {
			df.Emas[i].Values = make([]float64, len(df.candles.candles))
		}
	}
}

// BollingerBand
func (df *DataFrame) AddBollingerBand(n int, k1, k2 float64) {
	bb := new(BollingerBand)
	bb.N = n
	bb.K1 = k1
	bb.K2 = k2
	if n <= len(df.candles.candles) {
		closes := df.closes
		up1, center, down1 := talib.BBands(closes, n, k1, k1, 0)
		up2, center, down2 := talib.BBands(closes, n, k2, k2, 0)
		bb.Up2 = up2
		bb.Up1 = up1
		bb.Center = center
		bb.Down1 = down1
		bb.Down2 = down2
	} else {
		bb.Up2 = make([]float64, len(df.candles.candles))
		bb.Up1 = make([]float64, len(df.candles.candles))
		bb.Center = make([]float64, len(df.candles.candles))
		bb.Down1 = make([]float64, len(df.candles.candles))
		bb.Down2 = make([]float64, len(df.candles.candles))
	}
	df.BollingerBand = bb
}

// RSI
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
