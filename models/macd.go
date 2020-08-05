package models
import (
	"encoding/json"
)

type MACD struct {
	FastPeriod   int       `json:"fast_period"`
	SlowPeriod   int       `json:"slow_period"`
	SignalPeriod int       `json:"signal_period"`
	MACD         []float64 `json:"macd"`
	MACDSignal   []float64 `json:"macd_signal"`
	MACDHist     []float64 `json:"macd_hist"`
	fastEma      Ema
	slowEma      Ema
	signalEma    Ema
}

func NewMACD(inReal []float64, fastPeriod, slowPeriod, signalPeriod int) MACD {
	var macd MACD
	fastEma := NewEma([]float64{}, fastPeriod)
	slowEma := NewEma([]float64{}, slowPeriod)
	signalEma := NewEma([]float64{}, signalPeriod)
	m := make([]float64, 0)
	hist := make([]float64, 0)

	macd.FastPeriod = fastPeriod
	macd.SlowPeriod = slowPeriod
	macd.SignalPeriod = signalPeriod
	macd.MACD = m
	macd.MACDSignal = signalEma.Values
	macd.MACDHist = hist
	macd.fastEma = fastEma
	macd.slowEma = slowEma
	macd.signalEma = signalEma
	macd.Update(inReal)

	return macd
}

func  JsonUnmarshalMACD(row []byte)  (*MACD,error) {
	var macd = new(MACD)
	err := json.Unmarshal(row,macd)
	if err != nil {
		return nil, err
	}
	return macd ,nil
}

func (m *MACD) Update(inReal []float64) {
	if difflength := len(inReal) - len(m.MACD); difflength > 0 {
		baselength := len(m.MACD)
		m.fastEma.Update(inReal)
		m.slowEma.Update(inReal)
		macd := make([]float64, difflength)
		for i := range macd {
			macd[i] = m.fastEma.Values[baselength+i] - m.slowEma.Values[baselength+i]
		}
		m.MACD = append(m.MACD, macd...)
		m.signalEma.Update(m.MACD)
		m.MACDSignal = m.signalEma.Values
		hist := make([]float64, difflength)
		for i := range hist {
			hist[i] = m.MACD[baselength+i] - m.MACDSignal[baselength+i]
		}
		m.MACDHist = append(m.MACDHist, hist...)
	}
}
