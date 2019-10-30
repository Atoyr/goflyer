package cli

import (
	"context"
	"fmt"

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
	command.Subcommands = append(command.Subcommands, exportExecutionsCommand())

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
	exe := executor.GetExecutor()
	exe.FetchTickerAsync(ctx, make([]func(beforeticker, ticker models.Ticker), 0))

	return nil
}

func exportExecutionsCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "executions"
	command.Action = exportExecutionsAction
	command.Flags = []urfavecli.Flag{
		urfavecli.StringFlag{
			Name: "path, p",
			Value: "export file path",
		},
	}

	return command
}

func exportExecutionsAction(c *urfavecli.Context) error {
		
	path := c.String("path")
	if path == "" {
		return fmt.Errorf("export file path not found")
	}
	exe := executor.GetExecutor()
	executions,err := exe.GetExecution(0,0,0)
	if err != nil {
		return err
	}
	err = util.SaveJsonMarshalIndent(executions,path)
	if err != nil {
		return err
	}
	return nil
}
