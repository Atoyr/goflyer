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
	command.Subcommands = append(command.Subcommands, exportCandlesCommand())

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
	executor.FetchTickerAsync(ctx, make([]func(beforeticker, ticker models.Ticker), 0))

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
	executions,err := executor.GetExecution(0,0,0)
	if err != nil {
		return err
	}
	err = util.SaveJsonMarshalIndent(executions,path)
	if err != nil {
		return err
	}
	return nil
}

func exportCandlesCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "candles"
	command.Action = exportCandlesAction
	command.Flags = []urfavecli.Flag{
		urfavecli.StringFlag{
			Name: "path, p",
			Value: "",
		},
	}

	return command
}

func exportCandlesAction(c *urfavecli.Context) error {
		
	path := c.String("path")
	if path == "" {
		return fmt.Errorf("export file path not found")
	}
	candles := executor.GetCandles(models.GetDuration(models.Duration_1m))
	err := util.SaveJsonMarshalIndent(candles,path)
	if err != nil {
	fmt.Println(err)
		return err
	}
	return nil
}
