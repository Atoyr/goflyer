package configs

type ChartConfig struct {
	DataFrameConfigs []DataFrameConfig `json:"data_frame_configs"`
}

type DataFrameConfig struct {
	Duration string                    `json:"duration"`
	Smas     []MovingAveragePreference `json:"smas"`
	Emas     []MovingAveragePreference `json:"emas"`
	Macd     []MACDPreference          `json:"macd"`
}

type MovingAveragePreference struct {
	Period int `json:"period"`
}

type MACDPreference struct {
	FastPeriod   int       `json:"fast_period"`
	SlowPeriod   int       `json:"slow_period"`
	SignalPeriod int       `json:"signal_period"` 
}

func (CC *ChartConfig) Save(){

}
