package models

import (
	"path/filepath"

	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/util"
)

type config struct {
	apppath     string
	apikey      string
	timeoutmsec int64
}

// TODO move config
// GetConfig is Getting config
// if path is empty this use default config path
func GetConfig(path string) (config, error) {
	var c config
	apppath := path
	if apppath == "" {

		dirPath, err := util.CreateConfigDirectoryIfNotExists("goflyer")
		if err != nil {
			return c, err
		}
		dbfile := filepath.Join(dirPath, "goflyer.db")
		boltdb, err := db.GetBolt(dbfile)
		if err != nil {
			return c, err
		}
		apppath = ""
	}
	c.apppath = apppath
	return c
}
