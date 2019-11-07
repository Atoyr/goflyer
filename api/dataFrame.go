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
	jsondb ,_ := db.GetJsonDB()
	executor.ChangeDB(&jsondb)
  d, err := strconv.ParseInt(duration,10,64)
  if err != nil {
  	return err
  }
	df := executor.DataFrame(time.Duration(d))
	fmt.Println(count)
	
	return c.JSON(http.StatusOK, df)
}
