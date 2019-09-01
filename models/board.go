package models

type Board struct {
	MidPrice float64 `json:"mid_price,omitempty"`
	Bids     []Order `json:"bids,omitempty"`
	Asks     []Order `json:"asks,omitempty"`
}
type Order struct {
	Price float64 `json:"price,omitempty"`
	Size  float64 `json:"size,omitempty"`
}

func (b *Board) Merge(argBoard Board) (*Board, error) {
	if argBoard.MidPrice > 0 {
		b.MidPrice = argBoard.MidPrice
	}

	return b, nil
}
