package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/atoyr/goflyer/models"
) 

func handleCandleCollection(c echo.Context) error {
	ccs := models.NewCandleCollections()

	return c.JSON(http.StatusOK,ccs[""])
}
