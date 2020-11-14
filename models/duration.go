package models

import (
	"time"
	"sync"
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

var (
	once sync.Once
	durations []string
)

func Durations() []string{
	once.Do(func() {
		durations = make([]string,12)
		durations[0] = Duration_1m
		durations[1] = Duration_3m
		durations[2] = Duration_5m
		durations[3] = Duration_10m
		durations[4] = Duration_15m
		durations[5] = Duration_30m
		durations[6] = Duration_1h
		durations[7] = Duration_2h
		durations[8] = Duration_4h
		durations[9] = Duration_6h
		durations[10] = Duration_12h
		durations[11] = Duration_24h
	})
	return durations
}

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

func GetDurationString(duration time.Duration) string {
	d := ""
	switch duration {
	case 1 * time.Minute:
		d = "1m"
	case 3 * time.Minute:
		d = "3m"
	case 5 * time.Minute:
		d = "5m"
	case 10 * time.Minute:
		d = "10m"
	case 15 * time.Minute:
		d = "15m"
	case 30 * time.Minute:
		d = "30m"
	case 1 * time.Hour:
		d = "1h"
	case 2 * time.Hour:
		d = "2h"
	case 4 * time.Hour:
		d = "4h"
	case 6 * time.Hour:
		d = "6h"
	case 12 * time.Hour:
		d = "12h"
	case 24 * time.Hour:
		d = "24h"
	}
	return d
}
