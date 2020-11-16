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
      for {
        select {
        case <-ctx.Done():
          return nil
        case <-ticker.C:
          fmt.Println()
          fmt.Println(time.Now())
        case param := <-ch:
          for i := range  param {
            t := candle.Time
            mergeCandle(candle, param[i].DateTime(), param[i].Price)
            if t.Equal(candle.Time) {
              fmt.Printf("\r%s",candle.String())
            }else {
              fmt.Printf("%s",candle.String())
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
