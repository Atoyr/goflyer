package configs

import (
"net/url"
)
const base_url = "https://api.bitflyer.com/v1/"
const websocketScheme = "wss"
const websocketHost = "ws.lightstream.bitflyer.com"
const websocketPath = "/json-rpc"

type ClientConfig struct {
	baseURL           string
	websocket url.URL
}

func NewClientConfig() *ClientConfig {
	c := new(ClientConfig)
	c.baseURL = base_url
	c.websocket = url.URL{Scheme: websocketScheme, Host: websocketHost, Path: websocketPath}

	return c
}

func (c *ClientConfig) GetEndpoint(urlPath string) (endpoint string, err error) {
	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}
	apipath, err := url.Parse(urlPath)
	if err != nil {
		return "", err
	}
	return   baseURL.ResolveReference(apipath).String(),nil
}

func (c *ClientConfig) GetWebsocketString() string{
	return c.websocket.String()
}
