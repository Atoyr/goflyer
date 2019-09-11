package client

import (
"net/url"
)
const base_url = "https://api.bitflyer.com/v1/"
const websocketScheme = "wss"
const websocketHost = "ws.lightstream.bitflyer.com"
const websocketPath = "/json-rpc"

type apiConfig struct {
	baseURL           string
	websocket url.URL
}

func newAPIConfig() *apiConfig {
	c := new(apiConfig)
	c.baseURL = base_url
	c.websocket = url.URL{Scheme: websocketScheme, Host: websocketHost, Path: websocketPath}

	return c
}

func (api *apiConfig) GetEndpoint(urlPath string) (endpoint string, err error) {
	baseURL, err := url.Parse(api.baseURL)
	if err != nil {
		return "", err
	}
	apipath, err := url.Parse(urlPath)
	if err != nil {
		return "", err
	}
	return   baseURL.ResolveReference(apipath).String(),nil
}
