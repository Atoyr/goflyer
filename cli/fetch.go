package cli

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/executor"
	"github.com/atoyr/goflyer/models"
	"github.com/atoyr/goflyer/util"
	"github.com/atoyr/goflyer/configs"
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
	f := make([]func(beforeticker,ticker models.Ticker), 0)
	f = append(f, printFetchTicker)
	if c.Bool("save") {
		f = append(f, exe.SaveTicker)
	}
	exe.FetchTickerAsync(ctx, f)

	return nil
}

func printFetchTicker(beforeticker,ticker models.Ticker) {
	var status, ltp string
	okAtt := util.GetMultiColorAttribute(47,false)
	ngAtt := util.GetMultiColorAttribute(160,false)
	upAtt := util.GetMultiColorAttribute(20,false)
	stayAtt := util.GetMultiColorAttribute(188,false)
	downAtt := util.GetMultiColorAttribute(160,false)

	if ticker.Message == "" {
		status = util.ApplyAttribute("[  OK  ]",okAtt)
	}else {
		status = util.ApplyAttribute("[ FAIL ]",ngAtt )
	}

	ltp = fmt.Sprintf("%.2f",ticker.Ltp)
	if beforeticker.Ltp < ticker.Ltp {
	ltp = util.ApplyAttribute(ltp,upAtt)
} else if beforeticker.Ltp > ticker.Ltp {
	ltp = util.ApplyAttribute(ltp,downAtt)
} else {
	ltp = util.ApplyAttribute(ltp,stayAtt)
}

	fmt.Printf("\r%s  %s  |  LTP : %s", status,ticker.DateTime().Format(time.RFC3339),ltp)
}
