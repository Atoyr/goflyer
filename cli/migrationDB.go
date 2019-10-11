package cli

import (
	"fmt"
	"path/filepath"

	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/executor"
	"github.com/atoyr/goflyer/util"
	urfavecli "github.com/urfave/cli"
)

func migrationDBCommand() urfavecli.Command {
	var command urfavecli.Command
	command.Name = "migration"
	command.Action = migrationDBAction

	return command
}

func migrationDBAction(c *urfavecli.Context) error {
	jsonDB, err := db.GetJsonDB()
	if err != nil {
		return err
	}
	dirPath, err := util.CreateConfigDirectoryIfNotExists("goflyer")
	if err != nil {
		return err
	}
	dbfile := filepath.Join(dirPath, "goflyer.db")
	boltdb, err := db.GetBolt(dbfile)
	if err != nil {
		return err
	}
	exe := executor.GetExecutor(&jsonDB)
	fmt.Println("Start")
	exe.MigrationDB(&boltdb)
	fmt.Println("end")
	return nil
}
