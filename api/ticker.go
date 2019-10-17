package api

import (
	"net/http"

	"github.com/atoyr/goflyer/executor"
	"github.com/atoyr/goflyer/models"
	"github.com/labstack/echo"
)

func handleTicker(c echo.Context) error {
	context := c.(*Context)

	count := 100

	config, err := models.GetConfig()
	db := config.GetDB()
	exe := executor.GetExecutor(&db)
	tickers, err := exe.GetTicker(0, 0, 0)

	return c.JSON(http.StatusOK, tickers)
}
