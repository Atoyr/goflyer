package cli

import (
	"context"
	"fmt"
	"path/filepath"

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
			Usage:    "target choose ticker ...",
			Value:    "fetch target",
			Required: true,
		}}

	return command
}

func fetchAction(c *urfavecli.Context) error {
	return fetchTickerAction(c)
}

func fetchTickerAction(c *urfavecli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dirPath, err := util.CreateConfigDirectoryIfNotExists("goflyer")
	if err != nil {
		return err
	}
	dbfile := filepath.Join(dirPath, "goflyer.db")
	boltdb, err := db.GetBolt(dbfile)
	if err != nil {
		return err
	}
	exe := executor.GetExecutor(&boltdb)
	f := make([]func(models.Ticker), 0)
	f = append(f, func(ticker models.Ticker) { fmt.Println(ticker) })
	exe.FetchTickerAsync(ctx, f)

	return nil
}
