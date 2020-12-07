package controller

import (
  "fmt"
	"time"
  "context"
  "sort"

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
    fmt.Println(c.config.CanUsedDataFrameDurationMinute[i])
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
  ticker := time.NewTicker(1 * time.Second)
  defer ticker.Stop()
  fetchContext, cancel := context.WithCancel(ctx)
  go c.fetchExecuter(fetchContext)
  for {
    select {
    case <-ctx.Done():
      cancel()
      fmt.Println("Done")
      return
    case <-ticker.C:
      now := time.Now().Truncate(1 * time.Second)
      if now.Equal(now.Truncate(c.duration)) {
        c.dataframeSet.ApplyExecution()
        c.dataframeSet.UpdateTechnicalChartData()
        df, err := c.dataframeSet.GetDataFrame(1 * time.Minute)
        if err != nil {
          fmt.Println(err)
        }
        if i := len(df.Datetimes) - 1; i >= 0 {
          t := df.Datetimes[i].In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format(time.RFC3339)
          fmt.Printf("%d  Time : %s , Open : %7.0f , High : %7.0f , Low : %7.0f , Close : %7.0f , Volume : %f",i,t, df.Opens[i], df.Highs[i], df.Lows[i], df.Closes[i], df.Volumes[i])
          fmt.Println()
        }
      }
    }
  }
}

func (c *Controller) fetchExecuter(ctx context.Context) {
  ch := make(chan []bitflyer.Execution)
  childCtx, cancel := context.WithCancel(ctx)
  go c.client.GetRealtimeExecutions(childCtx, ch, "BTC_JPY")
  for {
    select {
    case <-ctx.Done():
      cancel()
      fmt.Println("Done")
      return
    case param := <-ch:
      sort.Slice(param, func(i, j int) bool { return param[i].ID < param[j].ID })
      for i := range  param {
        e := models.Execution{ Side: param[i].Side, Price : param[i].Price, Size : param[i].Size, Time : param[i].DateTime()}
        c.dataframeSet.AddExecution(e)
      }
    }
  }
}

func (c *Controller) Candles(duration time.Duration) ([]models.Candle, error) {
  if df, err := c.dataframeSet.GetDataFrame(duration); err != nil {
    return nil, err
  }else {
    fmt.Println(df.GetCandles())
    return df.GetCandles(), nil
  }
}
