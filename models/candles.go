package models

type Candles struct {
	productCode   string    
	duration      int64     
	candles []Candle
}

func NewCandles(productCode string, duration int64) Candles{
	var c Candles
	c.productCode = productCode
	c.duration = duration
	return c
}

func (cs Candles) AddExecution(executions ...Execution){

}

func (cs Candles) GetCandleOHLCs() []CandleOHLC {
	ohlcs := make([]CandleOHLC, len(cs.candles))
	for i := range cs.candles {
		ohlcs[i] = cs.candles[i].GetCandleOHLC()
	}

	return ohlcs
}
