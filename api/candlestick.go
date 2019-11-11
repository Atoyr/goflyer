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
	cs  := executor.GetCandles(models.GetDuration(duration))
	fmt.Println(count)
	return c.JSON(http.StatusOK, cs)
}

