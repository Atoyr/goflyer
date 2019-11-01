package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/executor"
	"github.com/labstack/echo"
)

func handleCandlestick(c echo.Context) error {
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
	jsondb ,_ := db.GetJsonDB()
	exe := executor.GetExecutor()
	exe.ChangeDB(&jsondb)
  d, err := strconv.ParseInt(duration,10,64)
  if err != nil {
  	return err
  }
	cs := exe.GetCandleOHLCs(time.Duration(d))
	start := len(cs) - count
	if start < 0 {
		start = 0
	}
	
	return c.JSON(http.StatusOK, cs[start:])
}

