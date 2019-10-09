package api

import (
	"net/http"
	"fmt" 
	"os"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/atoyr/goflyer/models"
)

func handleDataFrame(c echo.Context) error {
	context := c.(*Context)
	duration := context.Param("duration")

	if duration == "" {
		return fmt.Errorf("duration is required")
	}
	if df, ok := context.DataFrames[duration] ; !ok{
		jsonFile, err := os.Open("./testdata/tickers.json")
		if err != nil {
			return err
		}
		defer jsonFile.Close()
		raw, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			return err
		}
	df = models.NewDataFrame(models.BTC_JPY,models.GetDuration(duration))
		tickers, err := models.JsonUnmarshalTickers(raw)
	if err != nil {
		return err
	}
	for i := range tickers {
		df.AddTicker(tickers[i])
	}
	df.AddEmas(6)
	context.DataFrames[duration] = df
	}

	return c.JSON(http.StatusOK, context.DataFrames[duration])
}

