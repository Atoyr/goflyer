package cli

import (
	"context"
	"path/filepath"

	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/executor"
	"github.com/atoyr/goflyer/models"
	"github.com/atoyr/goflyer/util"
	urfavecli "github.com/urfave/cli"
)

func exportCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "export"
	command.Aliases = []string{"e"}
	command.Subcommands = append(command.Subcommands, exportTickersCommand())

	return command
}

func exportTickersCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "tickers"
	command.Action = exportTickersAction

	return command
}

func exportTickersAction(c *urfavecli.Context) error {
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
	exe.FetchTickerAsync(ctx, make([]func(beforeticker, ticker models.Ticker), 0))

	return nil
}
