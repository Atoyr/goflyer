package cli

import (
	urfavecli "github.com/urfave/cli"
)

func printCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "print"
	command.Aliases = []string{"p"}

	return command
}

func printTickerAction() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "tickers"
	command.Action = exportTickersAction

	return command
}
