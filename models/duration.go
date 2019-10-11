package models

import (
	"time"
)

const (
	Duration_1m  = "1m"
	Duration_3m  = "3m"
	Duration_5m  = "5m"
	Duration_10m = "10m"
	Duration_15m = "15m"
	Duration_30m = "30m"
	Duration_1h  = "1h"
	Duration_2h  = "2h"
	Duration_4h  = "4h"
	Duration_6h  = "6h"
	Duration_12h = "12h"
	Duration_24h = "24h"
)

func GetDuration(duration string) time.Duration {
	d := 24 * time.Hour
	switch duration {
	case "1m":
		d = 1 * time.Minute
	case "3m":
		d = 3 * time.Minute
	case "5m":
		d = 5 * time.Minute
	case "10m":
		d = 10 * time.Minute
	case "15m":
		d = 15 * time.Minute
	case "30m":
		d = 30 * time.Minute
	case "1h":
		d = 1 * time.Hour
	case "2h":
		d = 2 * time.Hour
	case "4h":
		d = 4 * time.Hour
	case "6h":
		d = 6 * time.Hour
	case "12h":
		d = 12 * time.Hour
	case "24h":
		d = 24 * time.Hour
	}
	return d
}
