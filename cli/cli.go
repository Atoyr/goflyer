package cli

import (
	urfavecli "github.com/urfave/cli"
)

func NewCli() *urfavecli.App {
	app := urfavecli.NewApp()
	app.Name = "goflyer"
	app.Author = "atoyr"

	app.Commands = append(app.Commands, exportCommand())
	app.Commands = append(app.Commands, runWebappsCommand())
	app.Commands = append(app.Commands, fetchCommand())
	app.Commands = append(app.Commands, migrationDBCommand())
	app.Commands = append(app.Commands, runCommand())

	return app
}
