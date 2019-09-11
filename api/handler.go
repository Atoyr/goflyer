package handler

import (
	"net/http"

	"github.com/labstack/echo"
) 

func GetEcho() (*echo.Echo, error ){
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})
	e.GET("/candleCollection", handleCandleCollection)
	return e, nil
}
