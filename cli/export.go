package cli

import urfavecli "github.com/urfave/cli"

func exportCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "export"
	command.Aliases = []string{"e"}
	command.Subcommands = append(command.Subcommands,exportTickersCommand())

	
	return command
}

func exportTickersCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "tickers"
	command.Action = exportTickersAction

	return command
}

func exportTickersAction(c *urfavecli.Context) error {
	return nil 
}
