package main

import (
  "os"
	"net/http"
	"time"
  "log"
  "fmt"
  "context"
  "strconv"

  "github.com/urfave/cli/v2"
  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
  "github.com/atoyr/goflyer/controller"
)

const APP_NAME = "goflyer"

var duration = 30 * time.Second

func main() {
  app := &cli.App{
    Action: func(c *cli.Context) error {
      clr := controller.New(APP_NAME)
      clr.SaveConfig()

      ctx := context.Background()
      childctx , cancel := context.WithCancel(ctx)
      clr.FetchExecuter(childctx)


      api := echo.New()
      api.Use(middleware.CORS())
      api.GET("/candles/:duration",
        func (ec echo.Context) error {
          ds := ec.Param("duration")
          d, err := strconv.Atoi(ds)
          if err != nil {
            return ec.JSON(http.StatusNotFound, nil)
          }
          c ,err := clr.Candles(time.Duration(d) * time.Minute)
          if err != nil {
            return ec.JSON(http.StatusNotFound, nil)
          }else {
            return ec.JSON(http.StatusOK, c)
          }
        })
      api.POST("/set_duration",
        func (c echo.Context) error {
          fmt.Println()
          // fmt.Println("call set_duration")
          // durationString := c.FormValue("duration")
          // d := models.GetDuration(durationString)
          // df.SetDuration(d)
          // fmt.Println(durationString)
          // for i := range df.Datetimes {
          // fmt.Printf("Time : %s , Open : %7.0f , High : %7.0f , Low : %7.0f , Close : %7.0f , Volume : %f",df.Datetimes[i], df.Opens[i], df.Highs[i], df.Lows[i], df.Closes[i], df.Volumes[i])
          // fmt.Println()
          // }
          return c.String(http.StatusOK, "")
        })
      api.POST("/cancel",
        func (c echo.Context) error {
          fmt.Println("cancel call")
          cancel()
          return c.String(http.StatusOK, "")
        })
      api.Start(":8080")
      return nil
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
