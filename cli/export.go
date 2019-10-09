package cli

import (
	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/executor"
	urfavecli "github.com/urfave/cli"
)

func exportCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "export"
	command.Aliases = []string{"e"}
	command.Subcommands = append(command.Subcommands, exportTickersCommand())

	return command
}

func exportCommandAction(c *urfavecli.Context) error {
	d, _ := db.GetJsonDB()
	executor.GetExecutor(&d)
	return nil
}
