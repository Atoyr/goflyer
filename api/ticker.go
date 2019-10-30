package api

import (
	"net/http"

	"github.com/atoyr/goflyer/executor"
	"github.com/labstack/echo"
)

func handleTicker(c echo.Context) error {

	exe := executor.GetExecutor()
	tickers, err := exe.GetTicker(0, 0, 0)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, tickers)
}
