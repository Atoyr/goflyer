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

	for _, bid := range argBoard.Bids {
		if bid.Size == 0 {
			for i, boardBid := range b.Bids {
				if bid.Price == boardBid.Price && bid.Size == boardBid.Size {
					b.Bids = append(b.Bids[:i], b.Bids[i+1:]...)
				}
			}
		} else {
			b.Bids = append(b.Bids, bid)
		}
	}

	for _, ask := range argBoard.Asks {
		if ask.Size == 0 {
			for i, boardAsk := range b.Asks {
				if ask.Price == boardAsk.Price && ask.Size == boardAsk.Size {
					b.Asks = append(b.Asks[:i], b.Asks[i+1:]...)
				}
			}
		} else {
			b.Asks = append(b.Asks, ask)
		}
	}
	return b, nil
}
