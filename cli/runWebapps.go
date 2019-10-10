package cli

import (
	"github.com/atoyr/goflyer/api"
	urfavecli "github.com/urfave/cli"
)

func runWebappsCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "runweb"
	command.Action = runWebappsAction

	return command
}

func runWebappsAction(c *urfavecli.Context) error {
	e := api.AppendHandler(api.GetEcho())
	e.Start(":8080")
	return nil
}
