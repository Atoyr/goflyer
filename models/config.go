package models

import (
	"path/filepath"
	"io/ioutil"

	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/util"
)

type config struct {
	appPath     string `json:"app_path"`
	apikey      string `json:"apikey"`
	dbfile      string `json:"dbfile"`
	timeoutmsec int64  `json:"timeoutmsec"`
}

const (
	appName    = "goflyer"
	configName = "config"
	dbName     = "goflyer.db"
)

// GetConfig is Getting config
// if path is empty this use default config path
func GetConfig() (config, error) {
	var c config
	appPath, err := util.CreateConfigDirectoryIfNotExists(appName)
	if err != nil {
		return config{}, err
	}
	configFile := filepath.Join(appPath,configName)
	if util.FileExists(configFile) {
		raw, err := ioutil.ReadFile(configFile)
		if err != nil {
			return config{}, err
		}
		err := json.Unmarshal(row, &c)
		if err != nil {
			return config{}, err
		}
	}else { 
		c.appPath = appPath
		c.dbfile = dbName
		err := c.Save()
		if err != nil {
			return c, err
		}
	}
	return c
}

func (c *config) Save() error {
	return util.SaveJsonMarshalIndent(c,c.appPath)
}

