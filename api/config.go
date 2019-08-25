package api

const base_url = "https://api.bitflyer.com/v1/"
const websocket_endpoint = "wss://ws.lightstream.bitflyer.com/json-rpc"

type apiConfig struct {
	baseURL           string
	websocketEndpoint string
}

func newAPIConfig() *apiConfig {
	c := new(apiConfig)

	return c
}
