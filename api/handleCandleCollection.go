package handler

import (
	"net/http"

	"github.com/labstack/echo"
) 

func handleCandleCollection(c echo.Context) error {
	return c.String(http.StatusOK, "handleCandleCollection") 
}
