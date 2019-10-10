package api

import (
	"net/http"

	"github.com/labstack/echo"
)

type Context struct {
	echo.Context
}

func GetEcho() *echo.Echo {
	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			context := &Context{c}
			return h(context)
		}
	})

	return e
}

func AppendHandler(e *echo.Echo) *echo.Echo {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})
	e.GET("/v1/DataFrame/BTC_JPY/:duration", handleDataFrame)
	return e
} 
