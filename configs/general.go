package configs

import (
	"io/ioutil"
	"path/filepath"
	"net/url"

	"encoding/json"

	"github.com/atoyr/goflyer/db"
	"github.com/atoyr/goflyer/util"
)

type generalConfig struct {
	appPath     string
	apikey      string
	dbtype      string
	dbfile      string
	timeoutmsec int64
	retrymsec   int64
	baseURL         string
	websocketScheme string
	websocketHost   string
	websocketPath   string
}

type outGeneralConfig struct {
	AppPath     string `json:"app_path"`
	Apikey      string `json:"apikey"`
	Dbtype      string `json:"dbtype"`
	Dbfile      string `json:"dbfile"`
	Timeoutmsec int64  `json:"timeoutmsec"`
	Retrymsec   int64  `json:"retrymsec"`
	BaseURL         string `json:"base_url"`
	WebsocketScheme string `json:"websocket_scheme"`
	WebsocketHost   string `json:"websocket_host"`
	WebsocketPath   string `json:"websocket_path"`
}

const (
	appName               = "goflyer"
	generalConfigFileName = "general.config"
	dbFileName            = "goflyer.db"
	base_url = "https://api.bitflyer.com/v1/"
	websocketScheme = "wss"
	websocketHost = "ws.lightstream.bitflyer.com"
	websocketPath = "/json-rpc"
)

// GetConfig is Getting generalConfig
// if path is empty this use default generalConfig path
func GetGeneralConfig() (generalConfig, error) {
	var c generalConfig
	appPath, err := util.CreateConfigDirectoryIfNotExists(appName)
	if err != nil {
		return generalConfig{}, err
	}
	generalConfigFile := filepath.Join(appPath, generalConfigFileName)
	if util.FileExists(generalConfigFile) {
		raw, err := ioutil.ReadFile(generalConfigFile)
		if err != nil {
			return generalConfig{}, err
		}
		var out outGeneralConfig
		err = json.Unmarshal(raw, &out)
		if err != nil {
			return generalConfig{}, err
		}
		c.appPath = out.AppPath
		c.apikey = out.Apikey
		c.dbtype = out.Dbtype
		c.dbfile = out.Dbfile
		c.timeoutmsec = out.Timeoutmsec
		c.retrymsec = out.Retrymsec
		c.baseURL = out.BaseURL
		c.websocketScheme = out.WebsocketScheme
		c.websocketHost = out.WebsocketHost
		c.websocketPath = out.WebsocketPath
	} else {
		c.appPath = appPath
		c.dbtype = "bolt"
		c.dbfile = dbFileName
		c.timeoutmsec = 5000
		c.timeoutmsec = 60000
		c.baseURL = base_url
		c.websocketScheme = websocketScheme
		c.websocketHost = websocketHost
		c.websocketPath= websocketPath
		err := c.Save()
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func (c *generalConfig) Save() error {
	out := outGeneralConfig{
		AppPath:     c.appPath,
		Apikey:      c.apikey,
		Dbtype:      c.dbtype,
		Dbfile:      c.dbfile,
		Timeoutmsec: c.timeoutmsec,
		Retrymsec:   c.retrymsec,
	BaseURL: c.baseURL,
WebsocketScheme: c.websocketScheme,
WebsocketHost: c.websocketHost,
WebsocketPath : c.websocketPath,
	}
	return util.SaveJsonMarshalIndent(out, filepath.Join(c.appPath, generalConfigFileName))
}

func (c *generalConfig) AppPath() string {
	return c.appPath
}

func (c *generalConfig) DBFile() string {
	return filepath.Join(c.appPath, c.dbfile)
}

func (c *generalConfig) GetDB() db.DB {
	switch c.dbtype {
	case "bolt":
		dbfile := filepath.Join(c.appPath, c.dbfile)
		db, err := db.GetBolt(dbfile)
		if err != nil {
			return nil
		}
		return &db
	default:
		return nil
	}
}

func (c *generalConfig) Timeoutmsec() int64 {
	return c.timeoutmsec
}

func (c *generalConfig) Retrymsec() int64 {
	return c.retrymsec
}

func (c *generalConfig) GetEndpoint(urlPath string) (endpoint string, err error) {
	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}
	apipath, err := url.Parse(urlPath)
	if err != nil {
		return "", err
	}
	return baseURL.ResolveReference(apipath).String(), nil
}

func (c *generalConfig) GetWebsocketString() string {
	websocket := url.URL{Scheme: c.websocketScheme, Host: c.websocketHost, Path: c.websocketPath}
	return websocket.String()
}
