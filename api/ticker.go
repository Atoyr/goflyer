package api

import (
	"net/http"

	"github.com/atoyr/goflyer/executor"
	"github.com/atoyr/goflyer/configs"
	"github.com/labstack/echo"
)

func handleTicker(c echo.Context) error {

	config, err := configs.GetGeneralConfig()
	if err != nil {
		return err
	}
	db := config.GetDB()
	exe := executor.GetExecutor(db)
	tickers, err := exe.GetTicker(0, 0, 0)

	return c.JSON(http.StatusOK, tickers)
}
