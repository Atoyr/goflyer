package client

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"path"
	"path/filepath"
	"sync"

	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/util"
)

type ClientConfig struct {
	AppPath         string `json:"app_path"`
	Apikey          string `json:"apikey"`
	Secretkey       string `json:"secret_key"`
	Timeoutmsec     int64  `json:"timeoutmsec"`
	Retrymsec       int64  `json:"retrymsec"`
	WebapiUrl       string `json:"webapi_url"`
	WebsocketScheme string `json:"websocket_scheme"`
	WebsocketHost   string `json:"websocket_host"`
	WebsocketPath   string `json:"websocket_path"`
	Dbtype          string `json:"dbtype"`
	Dbfile          string `json:"dbfile"`
}

const (
	appName         = "goflyer"
	configFileName  = "client.config"
	webapiUrl       = "https://api.bitflyer.com/v1/"
	websocketScheme = "wss"
	websocketHost   = "ws.lightstream.bitflyer.com"
	websocketPath   = "/json-rpc"
	dbFileName      = "goflyer.db"
)

var (
	once sync.Once

	config    generalConfig
	configerr error
)

// GetConfig is Getting ClientConfig.
// if path is empty this use default generalConfig path
func GetConfig() (ClientConfig, error) {
	once.Do(func() {
		configerr = config.Load()
	})
	return config, configerr
}

func (c *ClientConfig) Load() error {
	appPath, err := util.CreateConfigDirectoryIfNotExists(appName)
	if err != nil {
		return err
	}
	configFile := filepath.Join(appPath, configFileName)
	if util.FileExists(configFile) {
		raw, err := ioutil.ReadFile(configFile)
		if err != nil {
			return err
		}
		err = json.Unmarshal(raw, &c)
		if err != nil {
			return err
		}
	} else {
		c.AppPath = appPath
		c.Timeoutmsec = 5000
		c.Retrymsec = 60000
		c.WebapiUrl = webapiUrl
		c.WebsocketScheme = websocketScheme
		c.WebsocketHost = websocketHost
		c.WebsocketPath = websocketPath
		c.Dbtype = "bolt"
		c.Dbfile = dbFileName
		err := c.Save()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ClientConfig) Save() error {
	return util.SaveJsonMarshalIndent(c, filepath.Join(c.AppPath, configFileName))
}

func (c *ClientConfig) DBFilePath() string {
	return filepath.Join(c.AppPath, c.Dbfile)
}

func (c *ClientConfig) GetWebapiUrl(urlPath string) (string, error) {
	baseUrl, err := url.Parse(c.WebapiUrl)
	if err != nil {
		return "", err
	}
	baseUrl.Path = path.Join(baseUrl.Path, urlPath)
	return baseUrl.String(), nil
}

func (c *ClientConfig) GetDB() db.DB {
	switch c.Dbtype {
	case "bolt":
		dbfile := c.DBFilePath()
		db, err := db.GetBolt(dbfile)
		if err != nil {
			return nil
		}
		return &db
	default:
		return nil
	}
}

func (c *ClientConfig) GetWebsocketString() string {
	websocket := url.URL{Scheme: c.WebsocketScheme, Host: c.WebsocketHost, Path: c.WebsocketPath}
	return websocket.String()
}
