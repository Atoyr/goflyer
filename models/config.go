package models

import (
	"io/ioutil"
	"path/filepath"

	"encoding/json"

	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/util"
)

type config struct {
	appPath     string
	apikey      string
	dbtype      string
	dbfile      string
	timeoutmsec int64
	retrymsec   int64
}

type outconfig struct {
	AppPath     string `json:"app_path"`
	Apikey      string `json:"apikey"`
	Dbtype      string `json:"dbtype"`
	Dbfile      string `json:"dbfile"`
	Timeoutmsec int64  `json:"timeoutmsec"`
	Retrymsec   int64  `json:"retrymsec"`
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
	configFile := filepath.Join(appPath, configName)
	if util.FileExists(configFile) {
		raw, err := ioutil.ReadFile(configFile)
		if err != nil {
			return config{}, err
		}
		var out outconfig
		err = json.Unmarshal(raw, &out)
		if err != nil {
			return config{}, err
		}
		c.appPath = out.AppPath
		c.apikey = out.Apikey
		c.dbtype = out.Dbtype
		c.dbfile = out.Dbfile
		c.timeoutmsec = out.Timeoutmsec
		c.retrymsec = out.Retrymsec
	} else {
		c.appPath = appPath
		c.dbtype = "bolt"
		c.dbfile = dbName
		c.timeoutmsec = 5000
		c.timeoutmsec = 60000
		err := c.Save()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func (c *config) Save() error {
	out := outconfig{
		AppPath:     c.appPath,
		Apikey:      c.apikey,
		Dbtype:      c.dbtype,
		Dbfile:      c.dbfile,
		Timeoutmsec: c.timeoutmsec,
		Retrymsec:   c.retrymsec,
	}
	return util.SaveJsonMarshalIndent(out, filepath.Join(c.appPath, configName))
}

func (c *config) AppPath() string {
	return c.appPath
}

func (c *config) DBFile() string {
	return filepath.Join(c.appPath, c.dbfile)
}

func (c *config) GetDB() db.DB {
	switch c.dbtype {
	case "bolt":
		dbfile = filepath.Join(c.appPath, c.dbfile)
		db, err := db.GetBolt(dbfilr)
		if err != nil {
			return nil
		}
		return db
	default:
		return nil
	}
}

func (c *config) Timeoutmsec() int64 {
	return c.timeoutmsec
}

func (c *config) Retrymsec() int64 {
	return c.retrymsec
}
