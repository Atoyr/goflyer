package config

import (
  "os"
  "path"
  "path/filepath"
	"io/ioutil"
	"encoding/json"
	"net/url"
)

// Config is goflyer config file
type Config struct {
  appName         string
	Apikey          string `json:"apikey"`
	Secretkey       string `json:"secret_key"`
	Timeoutmsec     int64  `json:"timeoutmsec"`
	Retrymsec       int64  `json:"retrymsec"`
	WebapiUrl       string `json:"webapi_url"`
	WebsocketScheme string `json:"websocket_scheme"`
	WebsocketHost   string `json:"websocket_host"`
	WebsocketPath   string `json:"websocket_path"`
  // DataFrameUpdateDuration unit is second
  DataFrameUpdateDuration int `json:"data_frame_update_duration"`
}

const (
  ConfigFileName = "app.config"
	WebapiUrl       = "https://api.bitflyer.com/v1/"
	WebsocketScheme = "wss"
	WebsocketHost   = "ws.lightstream.bitflyer.com"
	WebsocketPath   = "/json-rpc"
)

func Load(appName string) (Config, error) {
  c := new(Config)
  c.loadDefaultValue()
  c.appName = appName
  configDir, err := createConfigDirectoryIfNotExists(appName)
  if err == nil {
    return *c, err
  }
  configFilePath := filepath.Join(configDir, ConfigFileName)
  if _, err := os.Stat(configFilePath); os.IsExist(err) {
    raw, err := ioutil.ReadFile(configFilePath)
    if err != nil {
      return *c, err
    }
    json.Unmarshal(raw, c)
  }

  return *c, nil
}

func (c *Config) loadDefaultValue() {
		c.Timeoutmsec = 5000
		c.Retrymsec = 60000
		c.WebapiUrl = WebapiUrl
		c.WebsocketScheme = WebsocketScheme
		c.WebsocketHost = WebsocketHost
		c.WebsocketPath = WebsocketPath
    c.DataFrameUpdateDuration = 30
}

func (c *Config) Save() error {
  configDir, err := createConfigDirectoryIfNotExists(c.appName)
  if err != nil {
    return err
  }
  configFilePath := filepath.Join(configDir, ConfigFileName)
  err = saveJsonMarshalIndent(*c, configFilePath)
  if err != nil {
    return err
  }
  return nil
}

func (c *Config) GetWebapiUrl(urlPath string) (string, error) {
	baseUrl, err := url.Parse(c.WebapiUrl)
	if err != nil {
		return "", err
	}
	baseUrl.Path = path.Join(baseUrl.Path, urlPath)
	return baseUrl.String(), nil
}

func (c *Config) GetWebsocketString() string {
	websocket := url.URL{Scheme: c.WebsocketScheme, Host: c.WebsocketHost, Path: c.WebsocketPath}
	return websocket.String()
}
