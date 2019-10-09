package api

import (
	"net/http"
	"fmt"

	"github.com/labstack/echo"
)

func handleDataFrame(c echo.Context) error {
	context := c.(*Context)
	key := context.Param("productCode")

	duration := c.QueryParam("duration")
	if duration == "" {
		return fmt.Errorf("duration is required")
	}

	return c.JSON(http.StatusOK, context.DataFrames[key])
}

