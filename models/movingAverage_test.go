package models

import (
	"testing"
)

func TestSimple(t *testing.T) {
	in := make([]float64, 10)
	in[0] = 1
	in[1] = 2
	in[2] = 3
	in[3] = 4
	in[4] = 5
	in[5] = 6
	in[6] = 7
	in[7] = 8
	in[8] = 9
	in[9] = 10

	period := 3

	out := make([]float64, 10)
	for i, v := range in {
		if i >= period-1 {
			sum := 0.0
			for _, x := range in[i-period+1 : i] {
				sum += x
			}
			out[i] = sum
		}
	}

	sma := NewSma(in, 3)
	for i := range out {
		if out[i] != sma.Values[i] {
			t.Fatalf("No %v want %v, but %v:", i, out[i], sma.Values[i])
		}

	}
}
