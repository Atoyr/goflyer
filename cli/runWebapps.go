package cli

import (
	"github.com/atoyr/goflyer/api"
	"github.com/atoyr/goflyer/backend"
	urfavecli "github.com/urfave/cli"
)

func runWebappsCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "runweb"
	command.Action = runWebappsAction

	return command
}

func runWebappsAction(c *urfavecli.Context) error {
	go func(){
		b := backend.GetEcho()
		b.Start(":3000")
	}()
	e := api.AppendHandler(api.GetEcho())
	e.Start(":8080")
	return nil
}
