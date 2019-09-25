package backend

import (
	"github.com/atoyr/goflyer/models"
	"github.com/labstack/echo"
)

func GetEcho(ccs models.DataFrames) *echo.Echo {
	e := echo.New()
	e.Static("/", "public/")
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(c)
		}
	})

	return e
}
