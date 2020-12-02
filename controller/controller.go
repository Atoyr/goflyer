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
	dataframeSet models.DataFrameSet
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
  c.dataframeSet = models.NewDataFrameSet("BTC_JPY")
  for i := range c.config.CanUsedDataFrameDurationMinute {
    c.dataframeSet.AddDataFrame( time.Duration( c.config.CanUsedDataFrameDurationMinute[i]) * time.Minute)
  }
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

func (c *Controller) Run(ctx context.Context) {

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
        c.dataframeSet.ApplyExecution()
        df, err := c.dataframeSet.GetDataFrame(15 * time.Minute)
        if err != nil {
          fmt.Println(err)
        }
        i := len(df.Datetimes) - 1
        fmt.Printf("Time : %s , Open : %7.0f , High : %7.0f , Low : %7.0f , Close : %7.0f , Volume : %f",df.Datetimes[i], df.Opens[i], df.Highs[i], df.Lows[i], df.Closes[i], df.Volumes[i])
        fmt.Println()
      case param := <-ch:
        for i := range  param {
          e := models.Execution{ Side: param[i].Side, Price : param[i].Price, Size : param[i].Size, Time : param[i].DateTime()}
          c.dataframeSet.AddExecution(e)
        }
      }
    }
  }()
}

func (c *Controller) Candles(duration time.Duration) ([]models.Candle, error) {
  if df, err := c.dataframeSet.GetDataFrame(duration); err != nil {
    return nil, err
  }else {
    fmt.Println(df.GetCandles())
    return df.GetCandles(), nil
  }
}
