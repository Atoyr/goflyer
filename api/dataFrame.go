package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/atoyr/goflyer/models"
	"github.com/labstack/echo"
)

func handleDataFrame(c echo.Context) error {
	context := c.(*Context)
	duration := context.Param("duration")

	if duration == "" {
		return fmt.Errorf("duration is required")
	}

	count := 100

	if countparam := context.QueryParam("count"); countparam != "" {
		c, err := strconv.Atoi(countparam)
		if err != nil {
			return err
		}
		count = c
	}

	if df, ok := context.DataFrames[duration]; !ok {
		jsonFile, err := os.Open("./testdata/tickers.json")
		if err != nil {
			return err
		}
		defer jsonFile.Close()
		raw, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			return err
		}
		df = models.NewDataFrame(models.BTC_JPY, models.GetDuration(duration))
		tickers, err := models.JsonUnmarshalTickers(raw)
		if err != nil {
			return err
		}
		start := len(tickers) - count
		if start < 0 {
			start = 0
		}
		for i := range tickers[start:] {
			df.AddTicker(tickers[i])
		}
		df.AddEmas(6)
		context.DataFrames[duration] = df
	}

	return c.JSON(http.StatusOK, context.DataFrames[duration])
}
