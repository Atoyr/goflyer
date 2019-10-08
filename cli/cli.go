package cli

import (
	urfavecli "github.com/urfave/cli"
)

func NewCli() *urfavecli.App {
	app := urfavecli.NewApp()
	app.Name = "goflyer"
	app.Author = "atoyr"

	app.Commands = append(app.Commands, exportCommand())

	return app
}
