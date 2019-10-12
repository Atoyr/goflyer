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

func getterCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "print"
	command.Aliases = []string{"p"}

	return command
}

func getterTickerCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "getter"
	command.Action = getterTickerAction

	return command
}

func getterTickerAction(c *urfavecli.Context) error {
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
