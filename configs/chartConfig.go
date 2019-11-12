package configs

type ChartConfig struct {
	DataFrameConfigs []DataFrameConfig
}

type DataFrameConfig struct {
	Duration string
	Smas []MovingAveragePreference
	Emas []MovingAveragePreference
	Macd []MACDPreference
}

type MovingAveragePreference struct{ 
	Period int 
}

type MACDPreference struct {
	FastPeriod   int       `json:"fast_period"`
	SlowPeriod   int       `json:"slow_period"`
	SignalPeriod int       `json:"signal_period"` 
}
