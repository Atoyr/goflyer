package main

import (
  "os"
	"time"
	"context"
  "log"
  "fmt"

  "github.com/urfave/cli/v2"
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
      df := models.NewDataFrame("BTC_JPY", 5 * time.Minute)
      t := time.Now().Truncate(1 * time.Minute)
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
