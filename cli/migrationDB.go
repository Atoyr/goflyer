package cli

import (
	"fmt"
	urfavecli "github.com/urfave/cli"
)

func migrationDBCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "migration"
	command.Action = migrationDBAction

	return command
}

func migrationDBAction(c *urfavecli.Context) error {
	fmt.Println("not implement")
	return nil
}
