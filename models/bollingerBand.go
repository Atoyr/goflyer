package models

type BollingerBand struct {
	N int
	K1 float64
	K2 float64
	Up2 []float64
	Up1 []float64
	Center [] float64
	Down1 []float64
	Down2 []float64
}

