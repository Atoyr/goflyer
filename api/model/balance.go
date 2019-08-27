package model


type Balance struct {
	CurrencyCode string `json:"currency_code"` 
	Amount       int    `json:"amount"` 
	Available    int    `json:"available"` 
}
