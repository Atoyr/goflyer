package cli

import (
	"context"
	"fmt"
	"log"

	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/executor"
	"github.com/atoyr/goflyer/models"
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
			Usage:    "target choose ticker ...",
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
	return fetchTickerAction(c)
}

func fetchTickerAction(c *urfavecli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config, err := models.GetConfig()
	if err != nil {
		log.Fatal(err)
		return err
	}
	boltdb, err := db.GetBolt(config.DBFile())
	if err != nil {
		return err
	}
	exe := executor.GetExecutor(&boltdb)
	f := make([]func(models.Ticker), 0)
	f = append(f, printFetchTicker)
	if c.Bool("save") {
		f = append(f, exe.SaveTicker)
	}
	exe.FetchTickerAsync(ctx, f)

	return nil
}

func printFetchTicker(ticker models.Ticker) {
	fmt.Printf("\r[  OK  ] ASK : %f\t BID : %f", ticker.BestAsk, ticker.BestBid)
}
