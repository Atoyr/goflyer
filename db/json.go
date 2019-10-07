package db

import (
	"github.com/atoyr/goflyer/util"
)

func ExportJsonForTickers(db DB ,path string) error {
	tickers, err := db.GetTickerAll()
	if err != nil {
		return err
	}
	return util.JsonMarshalIndent(tickers,path)
}
