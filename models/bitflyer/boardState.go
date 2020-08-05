package bitflyer

type BoardState struct {
	Health string `json:"health"`
	State  string `json:"state"`
	Data   struct {
		SpecialQuotation int `json:"special_quotation,omitempty"`
	} `json:"data,omitempty"`
}
