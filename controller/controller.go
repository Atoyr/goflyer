package controller

import (
  "fmt"
	"time"
  "context"

  "github.com/mattn/go-pubsub"

	"github.com/atoyr/goflyer/client"
	"github.com/atoyr/goflyer/client/bitflyer"
	"github.com/atoyr/goflyer/config"
	"github.com/atoyr/goflyer/models"
)

type Controller struct {
  appName   string
	client    *client.Client
	config    config.Config
	dataframe models.DataFrame

	duration time.Duration

  ps *pubsub.PubSub
}

func New(appName string) *Controller {
	c := new(Controller)
  conf, err := config.Load(appName)
  if err != nil {
    logf("%v",err)
    return nil
  }
  c.config = conf
	c.client = client.New(c.config.Apikey, c.config.Secretkey)
  c.dataframe = models.NewDataFrame("BTC_JPY", time.Duration(c.config.DataFrameDurationMinute) * time.Minute)
  c.ps = pubsub.New()

  c.applyConfig()
	return c
}

func (c *Controller) applyConfig() {
	c.client.SetTimeoutmsec(c.config.Timeoutmsec)
	c.client.SetRetrymsec(c.config.Retrymsec)
	c.client.SetWebApiUrl(c.config.WebapiUrl)
	c.client.SetWebsocket(c.config.WebsocketScheme, c.config.WebsocketHost, c.config.WebsocketPath)
	c.duration = time.Duration(c.config.DataFrameUpdateDuration) * time.Second
}

func (c *Controller) SaveConfig() {
  c.config.Save()
}

func (c *Controller) FetchExecuter(ctx context.Context) {
  ch := make(chan []bitflyer.Execution)
  ticker := time.NewTicker(c.duration)
  go c.client.GetRealtimeExecutions(ctx,ch,"BTC_JPY")
  go func() {
    for {
      select {
      case <-ctx.Done():
        fmt.Println("Done")
        return
      case <-ticker.C:
        c.dataframe.ApplyExecution()
        i := len(c.dataframe.Datetimes) - 1
        fmt.Printf("Time : %s , Open : %7.0f , High : %7.0f , Low : %7.0f , Close : %7.0f , Volume : %f",c.dataframe.Datetimes[i], c.dataframe.Opens[i], c.dataframe.Highs[i], c.dataframe.Lows[i], c.dataframe.Closes[i], c.dataframe.Volumes[i])
        fmt.Println()
      case param := <-ch:
        for i := range  param {
          e := models.Execution{ Side: param[i].Side, Price : param[i].Price, Size : param[i].Size, Time : param[i].DateTime()}
          c.dataframe.AddExecution(e)
        }
      }
    }
  }()
}

func (c *Controller) Candles() []models.Candle {
  return c.dataframe.GetCandles()
}
