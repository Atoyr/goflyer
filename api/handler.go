package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/atoyr/goflyer/models"
) 

func GetEcho() (*echo.Echo, error ){
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})
	e.GET("/candleCollection/:key", handleCandleCollection)
	return e, nil
}

func handleCandleCollection(c echo.Context) error {
	key := c.Param("key")
	ccs := models.NewCandleCollections()

	return c.JSON(http.StatusOK,ccs[key])
}
