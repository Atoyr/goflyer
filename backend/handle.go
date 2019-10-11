package backend

import (
	"github.com/labstack/echo"
)

func GetEcho() *echo.Echo {
	e := echo.New()
	e.Static("/", "public/")

	return e
}
