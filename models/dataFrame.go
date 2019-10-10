package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/atoyr/go-talib"
)

type DataFrames map[string]DataFrame

type DataFrame struct {
	ProductCode string
	Duration    time.Duration
	Candles     []Candle

	Smas          []Sma
	Emas          []Ema
	BollingerBand *BollingerBand
	Rsis          []RelativeStrengthIndex
	Macd          []MovingAverageConvergenceDivergence
}

const Open = "Open"
const Close = "Close"
const High = "High"
const Low = "Low"
const Volume = "Volume"

func NewDataFrame(productCode string, duration time.Duration) DataFrame {
	df := DataFrame{ProductCode: productCode, Duration: duration}
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

func (df *DataFrame) GetCountDataFrame(count int) *DataFrame{
	start  := len(df.Candles) - count
	if count < 0 {
		count = 0
	}
	ret := new(DataFrame) 
	ret.ProductCode = df.ProductCode
	ret.Duration = df.Duration
	ret.Candles = df.Candles[start:]
	ret.refreshChart()

	return ret
}

func (df *DataFrame) Name() string {
	fmt.Printf("%s_%s", df.ProductCode, df.Duration)
	return fmt.Sprintf("%s_%s", df.ProductCode, df.Duration)
}

func (df *DataFrame) AddTicker(ticker Ticker) error {
	dt := ticker.TruncateDateTime(df.Duration)
	if tail := len(df.Candles) - 1; tail < 0 {
		df.Candles = append(df.Candles, *NewCandle(df.ProductCode, df.Duration, ticker))
	} else if dt.Equal(df.Candles[tail].Time) {
		df.Candles[tail].AddTicker(ticker)
	} else if dt.After(df.Candles[tail].Time) {
		df.Candles = append(df.Candles, *NewCandle(df.ProductCode, df.Duration, ticker))
	} else if tail == 0 {
		c := []Candle{*NewCandle(df.ProductCode, df.Duration, ticker)}
		df.Candles = append(c, df.Candles...)
	} else {
		beforeTime := df.Candles[tail].Time
		for i := tail - 1; i >= 0; i-- {
			if dt.Equal(df.Candles[i].Time) {
				df.Candles[i].AddTicker(ticker)
			}
			if dt.After(df.Candles[i].Time) && dt.Before(beforeTime) {
				df.Candles = append(df.Candles[:i], *NewCandle(df.ProductCode, df.Duration, ticker))
				df.Candles = append(df.Candles, df.Candles[i+1:]...)
			}
		}
	}
	return nil
}

func (df *DataFrame) MergeCandle(candle Candle) error {
	if candle.GetDuration() != df.Duration {
		// TODO return error
		return nil
	}
	if len(df.Candles) == 0 {
		df.Candles = []Candle{candle}
		df.updateChart()
		return nil
	}

	max := len(df.Candles) - 1
	beforeTime := df.Candles[max].Time
	if candle.Time.Equal(beforeTime) {
		df.Candles[max] = candle
		df.refreshChart()
	} else if candle.Time.Before(df.Candles[max].Time) {
		for i := range df.Candles {
			if candle.Time.Equal(df.Candles[max-i].Time) {
				df.Candles[len(df.Candles)-1-i] = candle
				df.refreshChart()
				break
			} else if candle.Time.Before(beforeTime) && candle.Time.After(df.Candles[max-i].Time) {
				before := df.Candles[:max-i]
				after := df.Candles[max-i+1:]
				df.Candles = append(before, candle)
				df.Candles = append(df.Candles, after...)
				df.refreshChart()
				break
			}
		}
	} else {
		df.Candles = append(df.Candles, candle)
		df.updateChart()
	}
	return nil
}

func (df *DataFrame) Alls() (opens, closes, highs, lows, volumes []float64) {
	opens = make([]float64, len(df.Candles))
	closes = make([]float64, len(df.Candles))
	highs = make([]float64, len(df.Candles))
	lows = make([]float64, len(df.Candles))
	volumes = make([]float64, len(df.Candles))

	for i, v := range df.Candles {
		opens[i] = v.Open
		closes[i] = v.Close
		highs[i] = v.High
		lows[i] = v.Low
		volumes[i] = v.Volume
	}
	return
}

func (df *DataFrame) Values(valueType string) []float64 {
	ret, _ := df.LastOfValues(valueType, 0)
	return ret
}

func (df *DataFrame) LastOfValues(valueType string, from int) ([]float64, error) {
	if len(df.Candles) <= from {
		// TODO return error
		return nil, nil
	}
	// 123456 012345
	ret := make([]float64, len(df.Candles)-from)
	switch valueType {
	case Open:
		for i, v := range df.Candles[from:] {
			ret[i] = v.Open
		}
	case Close:
		for i, v := range df.Candles[from:] {
			ret[i] = v.Close
		}
	case High:
		for i, v := range df.Candles[from:] {
			ret[i] = v.High
		}
	case Low:
		for i, v := range df.Candles[from:] {
			ret[i] = v.Low
		}
	case Volume:
		for i, v := range df.Candles[from:] {
			ret[i] = v.Volume
		}
	default:
	}
	return ret, nil
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
	sma := NewSma(df.Values(Close), period)
	df.Smas = append(df.Smas, sma)
}

func (df *DataFrame) updateSmas() {
	for i := range df.Smas {
		df.Smas[i].Update(df.Values(Close))
	}
}

func (df *DataFrame) refreshSmas() {
	for i, sma := range df.Smas {
		if len(df.Candles) > sma.Period {
			df.Smas[i].Values = talib.Sma(df.Values(Close), sma.Period)
		} else {
			df.Smas[i].Values = make([]float64, len(df.Candles))
		}
	}
}

// EMA
func (df *DataFrame) AddEmas(period int) {
	ema := NewEma(df.Values(Close), period)
	df.Emas = append(df.Emas, ema)
}

func (df *DataFrame) updateEmas() {
	for i := range df.Emas {
		df.Emas[i].Update(df.Values(Close))
	}
}

func (df *DataFrame) refreshEmas() {
	for i, ema := range df.Emas {
		if len(df.Candles) > ema.Period {
			df.Emas[i].Values = talib.Ema(df.Values(Close), ema.Period)
		} else {
			df.Emas[i].Values = make([]float64, len(df.Candles))
		}
	}
}

// BollingerBand
func (df *DataFrame) AddBollingerBand(n int, k1, k2 float64) {
	bb := new(BollingerBand)
	bb.N = n
	bb.K1 = k1
	bb.K2 = k2
	if n <= len(df.Candles) {
		closes := df.Values(Close)
		up1, center, down1 := talib.BBands(closes, n, k1, k1, 0)
		up2, center, down2 := talib.BBands(closes, n, k2, k2, 0)
		bb.Up2 = up2
		bb.Up1 = up1
		bb.Center = center
		bb.Down1 = down1
		bb.Down2 = down2
	} else {
		bb.Up2 = make([]float64, len(df.Candles))
		bb.Up1 = make([]float64, len(df.Candles))
		bb.Center = make([]float64, len(df.Candles))
		bb.Down1 = make([]float64, len(df.Candles))
		bb.Down2 = make([]float64, len(df.Candles))
	}
	df.BollingerBand = bb
}

// RSI
func (df *DataFrame) AddRsis(period int) {
	rsi := NewRelativeStrengthIndex(df.Values(Close), period)
	df.Rsis = append(df.Rsis, rsi)
}

func (df *DataFrame) updateRsis() {
	for i := range df.Rsis {
		df.Rsis[i].Update(df.Values(Close))
	}
}

func (df *DataFrame) refreshRsis() {
}

// MACD
func (df *DataFrame) AddMacd(fastPeriod, slowPeriod, signalPeriod int) {
	closes := df.Values(Close)
	macd := NewMovingAverageConvergenceDivergence(closes, fastPeriod, slowPeriod, signalPeriod)
	df.Macd = append(df.Macd, macd)
}

func (df *DataFrame) updateMacd() {
	for i := range df.Macd {
		df.Macd[i].Update(df.Values(Close))
	}
}
