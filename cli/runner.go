package cli

import (
	"context"

	"github.com/atoyr/goflyer/executor"
	urfavecli "github.com/urfave/cli"
)

func runCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "run"
	command.Aliases = []string{"r"}
	command.Action = runAction

	return command
}

func runAction(c *urfavecli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	executor.RunAsync(ctx)
	return nil
}


