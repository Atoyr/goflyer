package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
  "path"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"net/url"

	"github.com/gorilla/websocket"
)

const (
  BTC_JPY = "BTC_JPY"
  timeoutmsec = 5000
  retrymsec = 60000
	webapiUrl       = "https://api.bitflyer.com/v1/"
	websocketScheme = "wss"
	websocketHost   = "ws.lightstream.bitflyer.com"
	websocketPath   = "/json-rpc"
)


// Client is bitflyer api client
type Client struct {
	key        string
	secret     string
	httpClient *http.Client


	timeoutmsec     int64
	retrymsec       int64
	webapiUrl       string
	websocketScheme string
	websocketHost   string
	websocketPath   string
}

// JsonRPC2 is Json rpc 2 struct
type JsonRPC2 struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Id      *int        `json:"id,omitempty"`
}

// SubscribeParams is Json rpc 2 Params
type SubscribeParams struct {
	Channel string `json:"channel"`
}

// New is create Client
func New(key, secret string) *Client {
	client := new(Client)
	client.key = key
	client.secret = secret
	client.httpClient = new(http.Client)

  client.timeoutmsec = timeoutmsec
  client.retrymsec = retrymsec

	return client
}

func NewJsonRPC2Subscribe() *JsonRPC2 {
	jsonRPC2 := new(JsonRPC2)
	jsonRPC2.Version = "2.0"
	jsonRPC2.Method = "subscribe"
	return jsonRPC2
}

func (c *Client) SetTimeoutmsec(t int64) {
  if t < 0 {
    c.timeoutmsec = timeoutmsec
  }else {
    c.timeoutmsec = t
  }
}

func (c *Client) SetRetrymsec(t int64) {
  if t < 0 {
    c.retrymsec = retrymsec
  }else {
    c.retrymsec = t
  }
}

func (c *Client) SetWebApiUrl(url string) {
  c.webapiUrl = url
}

func (c *Client) SetWebsocket(scheme, host, path string) {
  c.websocketScheme = scheme
  c.websocketHost = host
  c.websocketPath = path
}

func (c *Client) getWebapiUrl(urlPath string) (string, error) {
  baseUrl, err := url.Parse(c.webapiUrl)
	if err != nil {
		return "", err
	}
	baseUrl.Path = path.Join(baseUrl.Path, urlPath)
	return baseUrl.String(), nil
}

func (c *Client) getWebsocketString() string {
	websocket := url.URL{Scheme: c.websocketScheme, Host: c.websocketHost, Path: c.websocketPath}
	return websocket.String()
}

// header is create api call header
func (api *Client) header(method, endpoint string, body []byte) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	message := timestamp + method + endpoint + string(body)

	mac := hmac.New(sha256.New, []byte(api.secret))
	mac.Write([]byte(message))
	sign := hex.EncodeToString(mac.Sum(nil))

	return map[string]string{
		"ACCESS-KEY":       api.key,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      sign,
		"Context-Type":     "application/json",
	}
}

// doRequest is request api
func (api *Client) doRequest(method, url string, query map[string]string, data []byte) (body []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	for k, v := range api.header(method, req.URL.RequestURI(), body) {
		req.Header.Add(k, v)
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

func (api *Client) doWebsocketRequest(ctx context.Context, jsonRPC2 JsonRPC2, ch chan<- interface{}) error {
	c, _, err := websocket.DefaultDialer.Dial(api.getWebsocketString(), nil)
	if err != nil {
    e := fmt.Errorf("function=Client.doWebsocketRequest, action=Websocket Dial, err=%v \n", err)
    logf("error : %v", e)
		return e
	}

	defer c.Close()
	if err := c.WriteJSON(&jsonRPC2); err != nil {
    e := fmt.Errorf("function=Client.doWebsocketRequest, action=Write Json, err=%v \n", err)
    logf("error : %v", e)
		return e
	}
	c.SetWriteDeadline(time.Now().Add(10 * time.Second))

	if err != nil {
    e := fmt.Errorf("function=Client.doWebsocketRequest, action=Get Config, err=%v \n", err)
    logf("error : %v", e)
		return e
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			message := new(JsonRPC2)
			if err := c.ReadJSON(message); err != nil {
				if api.retrymsec > 0 {
					time.Sleep(time.Duration(api.retrymsec) * time.Millisecond)
				} else {
					 logf("function=Client.doWebsocketRequest, action=Read Json, message=%v, err=%v \n", message, err)
				}
			}

			if message.Method == "channelMessage" {
				switch params := message.Params.(type) {
				case map[string]interface{}:
					if v, ok := params["message"]; ok {
						ch <- v
					}
				}
			}
		}
	}
}
