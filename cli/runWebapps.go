package cli

import (
	"github.com/atoyr/goflyer/api"
	"github.com/atoyr/goflyer/models"
	urfavecli "github.com/urfave/cli"
)

func runWebappsCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "runweb"
	command.Action = runWebappsAction

	return command
}

func runWebappsAction(c *urfavecli.Context) error {
	dfs := make(map[string]models.DataFrame,0)
	e := api.AppendHandler(api.GetEcho(dfs))
	e.Start(":8080")
	return nil
}
