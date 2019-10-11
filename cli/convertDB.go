package cli

import (
	urfavecli "github.com/urfave/cli"
)

func convertCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "convert"
	command.Action = convertDBAction

	return command
}

func convertDBAction(c *urfavecli.Context) error {
	return nil
}
