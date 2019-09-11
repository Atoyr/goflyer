package api

import (
	"net/http"

	"github.com/atoyr/goflyer/models"
	"github.com/labstack/echo"
)

type Context struct {
	echo.Context
	CandleCollections models.CandleCollections
}

func GetEcho(ccs models.CandleCollections) *echo.Echo {
	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			context := &Context{c, ccs}
			return h(context)
		}
	})

	return e
}

func AppendHandler(e *echo.Echo) *echo.Echo {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})
	e.GET("/candleCollection/:key", handleCandleCollection)
	return e
}

func handleCandleCollection(c echo.Context) error {
	context := c.(*Context)
	key := context.Param("key")

	return c.JSON(http.StatusOK, context.CandleCollections[key])
}
