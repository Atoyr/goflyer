package models

import (
	"fmt"
	"time"
)

func ConvertDurationToString(d time.Duration) string {
	if h := d / time.Hour; h > 0 {
		return fmt.Sprintf("%dH", h)
	}
	return ""
}
