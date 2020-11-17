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
      candle := new(models.Candle)
      candle.Time = time.Now().Truncate(duration)
      ticker := time.NewTicker(30 * time.Second)
      df := models.NewDataFrame("BTC_JPY", 5 * time.Minute)
      for {
        select {
        case <-ctx.Done():
          return nil
        case <-ticker.C:
          df.ApplyExecution()
          fmt.Println()
          fmt.Println("dataframe")
          for i := range df.Datetimes {
            fmt.Printf("  time : %s , Open : %7.0f , High : %7.0f , Low : %7.0f , Close : %7.0f , Volume : %f",df.Datetimes[i], df.Opens[i], df.Highs[i], df.Lows[i], df.Closes[i], df.Volumes[i])
            fmt.Println()
          }
        case param := <-ch:
          for i := range  param {
            t := candle.Time
            mergeCandle(candle, param[i].DateTime(), param[i].Price)
            e := models.Execution{ Side: param[i].Side, Price : param[i].Price, Size : param[i].Size, Time : param[i].DateTime()}
            df.AddExecution(e)
            if t.Equal(candle.Time) {
              fmt.Printf("\r%s",candle.String())
            }else {
              fmt.Printf("\r%s",candle.String())
              fmt.Println()
            }
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

func mergeCandle(c *models.Candle, execDate time.Time, price float64) *models.Candle {
  dt := execDate.Truncate(duration)
  if c.Time.Equal(dt) {
			c.Close = price
			if c.High < price {
				c.High = price
			} else if c.Low > price {
				c.Low = price
			}
  }else {
    c.Time = dt
    c.Open = price
    c.High = price
    c.Low = price
    c.Close = price
  }
  return c
}
