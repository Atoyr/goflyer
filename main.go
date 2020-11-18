package main

import (
  "os"
	"net/http"
	"time"
	"context"
  "log"
  "fmt"

  "github.com/urfave/cli/v2"
  "github.com/labstack/echo"
  "github.com/atoyr/goflyer/config"
	"github.com/atoyr/goflyer/models"
  "github.com/atoyr/goflyer/client"
	"github.com/atoyr/goflyer/client/bitflyer"
)

const APP_NAME = "goflyer"

var duration = 30 * time.Second

func main() {
  app := &cli.App{
    Action: func(c *cli.Context) error {
      conf, _ := config.Load(APP_NAME)
      conf.Save()
      cli := client.New("","")

      // apply config
      cli.SetTimeoutmsec(conf.Timeoutmsec)
      cli.SetRetrymsec(conf.Retrymsec)
      cli.SetWebApiUrl(conf.WebapiUrl)
      cli.SetWebsocket(conf.WebsocketScheme,conf.WebsocketHost,conf.WebsocketPath)
      duration = time.Duration(conf.DataFrameUpdateDuration) * time.Second

      ctx := context.Background()
      ch := make(chan []bitflyer.Execution)
      go cli.GetRealtimeExecutions(ctx,ch,"BTC_JPY")

      fmt.Println("execute")
      ticker := time.NewTicker(30 * time.Second)
      df := models.NewDataFrame("BTC_JPY", 1 * time.Minute)
      // df.SetLogger( log.New(os.Stdout, "myapp", log.LstdFlags))

      t := time.Now().Truncate(1 * time.Hour)

      api := echo.New()
      api.GET("/last_candle",
        func (c echo.Context) error {
          if index := len(df.Datetimes) - 1; index < 0 {
            return c.String(http.StatusServiceUnavailable, "fooo")
          }else {
            return c.JSON(http.StatusOK, df.GetCandles())
          }
        })
      api.POST("/set_duration",
        func (c echo.Context) error {
          fmt.Println()
          fmt.Println("call set_duration")
          durationString := c.FormValue("duration")
          d := models.GetDuration(durationString)
          df.SetDuration(d)
          fmt.Println(durationString)
          for i := range df.Datetimes {
          fmt.Printf("Time : %s , Open : %7.0f , High : %7.0f , Low : %7.0f , Close : %7.0f , Volume : %f",df.Datetimes[i], df.Opens[i], df.Highs[i], df.Lows[i], df.Closes[i], df.Volumes[i])
          fmt.Println()
          }
          return c.String(http.StatusOK, string(d))
        })
      go api.Start(":8080")
      for {
        select {
        case <-ctx.Done():
          return nil
        case <-ticker.C:
          df.ApplyExecution()
          i := len(df.Datetimes) - 1
          if df.Datetimes[i].Equal(t) {
            fmt.Printf("\r")
          }else {
            fmt.Println()
            t = df.Datetimes[i]
          }
          fmt.Printf("Time : %s , Open : %7.0f , High : %7.0f , Low : %7.0f , Close : %7.0f , Volume : %f",df.Datetimes[i], df.Opens[i], df.Highs[i], df.Lows[i], df.Closes[i], df.Volumes[i])
        case param := <-ch:
          for i := range  param {
            e := models.Execution{ Side: param[i].Side, Price : param[i].Price, Size : param[i].Size, Time : param[i].DateTime()}
            df.AddExecution(e)
          }
        }
      }
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
