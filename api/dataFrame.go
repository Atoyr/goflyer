package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/atoyr/goflyer/executor"
	"github.com/atoyr/goflyer/models"
	"github.com/labstack/echo"
)

func handleCandlestick(c echo.Context) error {
	context := c.(*Context)
	duration := context.Param("duration")

	if duration == "" {
		duration = "1m"
	}

	count := 100

	if countparam := context.QueryParam("count"); countparam != "" {
		c, err := strconv.Atoi(countparam)
		if err != nil {
			return err
		}
		count = c
	}
	cs := executor.GetCandles(models.GetDuration(duration))
	fmt.Println(count)
	return c.JSON(http.StatusOK, cs)
}

func handleSma(c echo.Context) error {
	context := c.(*Context)
	duration := context.Param("duration")

	if duration == "" {
		duration = "1m"
	}

	period := 4
	count := 100

	if periodparam := context.QueryParam("period"); periodparam != "" {
		p, err := strconv.Atoi(periodparam)
		if err != nil {
			return err
		}
		period = p
	}
	if countparam := context.QueryParam("count"); countparam != "" {
		c, err := strconv.Atoi(countparam)
		if err != nil {
			return err
		}
		count = c
	}
	df := executor.DataFrame(models.GetDuration(duration))
	sma := df.GetSma(period)
	if l := len(sma.Values) ; l > count {
		sma.Values = sma.Values[l - count :]
	}

	return c.JSON(http.StatusOK, sma)
}
