package api

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/atoyr/goflyer/model"
	"github.com/gorilla/websocket"
)

type APIClient struct {
	key        string
	secret     string
	httpClient *http.Client
	config     *apiConfig
}

type JsonRPC2 struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Id      *int        `json:"id,omitempty"`
}

type SubscriveParams struct {
	Channel string `json:"channel"`
}

func New(key, secret string) *APIClient {
	client := new(APIClient)
	client.key = key
	client.secret = secret
	client.httpClient = new(http.Client)
	client.config = newAPIConfig()

	return client
}

func (api *APIClient) header(method, endpoint string, body []byte) map[string]string {
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

func (api *APIClient) doRequest(method, urlPath string, query map[string]string, data []byte) (body []byte, err error) {
	endpoint, err := api.config.GetEndpoint(urlPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
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

func (api *APIClient) doWebsocketRequest(ctx context.Context, jsonRPC2 JsonRPC2, ch chan<- interface{}) {
	c, _, err := websocket.DefaultDialer.Dial(api.config.websocket.String(), nil)
	if err != nil {
		log.Fatalf("function=APIClient.doWebsocketRequest, action=Websocket Dial, argslen=3, args=%v , %v , %v err=%s \n", ctx, jsonRPC2, ch, err.Error())
	}

	defer c.Close()
	if err := c.WriteJSON(&jsonRPC2); err != nil {
		log.Fatalf("function=APIClient.doWebsocketRequest, action=Write Json, argslen=3, args=%v , %v , %v err=%s \n", ctx, jsonRPC2, ch, err.Error())
	}
	c.SetWriteDeadline(time.Now().Add(10 * time.Second))

	for {
		select {
		case <-ctx.Done():
			return
		default:
			message := new(JsonRPC2)
			if err := c.ReadJSON(message); err != nil {
				log.Fatalln("read:", err)
			}

			if message.Method == "channelMessage" {
				switch params := message.Params.(type) {
				case map[string]interface{}:
					for k, v := range params {
						if k == "message" {
							ch <- v
						}
					}
				}
			}
		}
	}
}

func (api *APIClient) GetTicker(productCode string) (ticker *model.Ticker, err error) {
	url := "getticker"
	query := map[string]string{}
	query["product_code"] = productCode

	resp, err := api.doRequest("GET", url, query, nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp, ticker)
	if err != nil {
		return nil, err
	}
	return ticker, nil
}

func (api *APIClient) GetRealtimeTicker(ctx context.Context, ch chan<- model.Ticker, productCode string) {
	jsonRPC2 := new(JsonRPC2)
	jsonRPC2.Version = "2.0"
	jsonRPC2.Method = "subscribe"

	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	jsonRPC2.Params = SubscriveParams{Channel: fmt.Sprintf("lightning_ticker_%s", productCode)}

	var paramCh = make(chan interface{})
	go api.doWebsocketRequest(childctx, *jsonRPC2, paramCh)

OUTER:
	for {
		select {
		case <-ctx.Done():
			return

		default:
			param := <-paramCh
			marchalTick, err := json.Marshal(param)
			if err != nil {
				continue OUTER
			}
			ticker := new(model.Ticker)
			if err := json.Unmarshal(marchalTick, &ticker); err != nil {
				continue OUTER
			}
			ch <- *ticker
		}

	}
}

func (api *APIClient) GetExecutions(productCode string, beforeID, afterID string, count int) (executions []model.Execution, err error) {
	url := "getexecutions"
	query := map[string]string{}
	query["product_code"] = productCode
	if beforeID != "" {
		query["before"] = beforeID
	}
	if afterID != "" {
		query["before"] = afterID
	}
	if count < 0 {
		query["count"] = string(count)
	}
	resp, err := api.doRequest("GET", url, query, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &executions)
	if err != nil {
		return nil, err
	}

	return executions, nil
}

func (api *APIClient) GetBalance() (balances []model.Balance, err error) {
	url := "me/getbalance"
	resp, err := api.doRequest("GET", url, map[string]string{}, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp, &balances)
	if err != nil {
		return nil, err
	}

	return balances, nil
}
