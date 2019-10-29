package cli

import (
	"context"
	"fmt"
	"log"
	"time"
	"strings"

	"github.com/atoyr/goflyer/configs"
	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/executor"
	"github.com/atoyr/goflyer/models"
	"github.com/atoyr/goflyer/util"
	urfavecli "github.com/urfave/cli"
)

func fetchCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "fetch"
	command.Aliases = []string{"f"}
	command.Action = fetchAction
	command.Flags = []urfavecli.Flag{
		urfavecli.StringFlag{
			Name:     "target , t",
			Usage:    "target choose ticker ,executoin ...",
			Value:    "fetch target",
			Required: true,
		},
		urfavecli.BoolFlag{
			Name:  "save , s",
			Usage: "save db",
		}}

	return command
}

func fetchAction(c *urfavecli.Context) error {
	target := c.String("target")
	switch strings.ToLower(target) {
	case "ticker":
		return fetchTickerAction(c)
	case "execution":
		return fetchExecutionAction(c)
	default:
		return fmt.Errorf("target not found")
	} 
}

func fetchTickerAction(c *urfavecli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config, err := configs.GetGeneralConfig()
	if err != nil {
		log.Fatal(err)
		return err
	}
	boltdb, err := db.GetBolt(config.DBFile())
	if err != nil {
		return err
	}
	exe := executor.GetExecutor(&boltdb)
	f := make([]func(beforeticker, ticker models.Ticker), 0)
	f = append(f, printFetchTicker)
	if c.Bool("save") {
		f = append(f, exe.SaveTicker)
	}
	exe.FetchTickerAsync(ctx, f)

	return nil
}

func fetchExecutionAction(c *urfavecli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config, err := configs.GetGeneralConfig()
	if err != nil {
		log.Fatal(err)
		return err
	}
	boltdb, err := db.GetBolt(config.DBFile())
	if err != nil {
		return err
	}
	exe := executor.GetExecutor(&boltdb)
	f := make([]func(beforeexecution, execution models.Execution), 0)
	f = append(f, printFetchExecution)
	if c.Bool("save") {
		f = append(f, exe.SaveExecution)
	}
	exe.FetchExecutionAsync(ctx, f)

	return nil
}

func printFetchTicker(beforeticker, ticker models.Ticker) {
	var status, ltp, ask, bid string
	okAtt := util.GetMultiColorAttribute(47, false)
	ngAtt := util.GetMultiColorAttribute(160, false)
	upAtt := util.GetMultiColorAttribute(20, false)
	stayAtt := util.GetMultiColorAttribute(188, false)
	downAtt := util.GetMultiColorAttribute(160, false)

	if ticker.Message == "" {
		status = util.ApplyAttribute("[  OK  ]", okAtt)
	} else {
		status = util.ApplyAttribute("[ FAIL ]", ngAtt)
	}

	ltp = fmt.Sprintf("%.2f", ticker.Ltp)
	if beforeticker.Ltp < ticker.Ltp {
		ltp = util.ApplyAttribute(ltp, upAtt)
	} else if beforeticker.Ltp > ticker.Ltp {
		ltp = util.ApplyAttribute(ltp, downAtt)
	} else {
		ltp = util.ApplyAttribute(ltp, stayAtt)
	}

	ask = fmt.Sprintf("%.2f", ticker.BestAsk)
	if beforeticker.BestAsk < ticker.BestAsk {
		ask = util.ApplyAttribute(ask, upAtt)
	} else if beforeticker.BestAsk > ticker.BestAsk {
		ask = util.ApplyAttribute(ask, downAtt)
	} else {
		ask = util.ApplyAttribute(ask, stayAtt)
	}

	bid = fmt.Sprintf("%.2f", ticker.BestBid)
	if beforeticker.BestBid < ticker.BestBid {
		bid = util.ApplyAttribute(bid, upAtt)
	} else if beforeticker.BestBid > ticker.BestBid {
		bid = util.ApplyAttribute(bid, downAtt)
	} else {
		bid = util.ApplyAttribute(bid, stayAtt)
	}

	fmt.Printf("\r%s  %s  |  ASK : %s  |  BID : %s  |  LTP : %s", status, ticker.DateTime().Format(time.RFC3339), ask, bid, ltp)
}

func printFetchExecution(beforeexecution, execution models.Execution) {
	fmt.Printf("%s | %.0f |  %s | Price %.2f \n",execution.DateTime().Format(time.RFC3339),execution.ID,execution.Side, execution.Price)
}
