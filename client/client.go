package client

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
	"strings"
	"time"

	"github.com/atoyr/goflyer/configs"
	"github.com/atoyr/goflyer/models/bitflyer"
	"github.com/gorilla/websocket"
)

const bitflyerURL = "https://api.bitflyer.com/v1/"

// APIClient is bitflyer api client
type APIClient struct {
	key        string
	secret     string
	httpClient *http.Client
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

// SubscriveParams is Json rpc 2 Params
type SubscriveParams struct {
	Channel string `json:"channel"`
}

// New is create APIClient
func New(key, secret string) *APIClient {
	client := new(APIClient)
	client.key = key
	client.secret = secret
	client.httpClient = new(http.Client)

	return client
}

// header is create api call header
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

// doRequest is request api
func (api *APIClient) doRequest(method, url string, query map[string]string, data []byte) (body []byte, err error) {
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

func (api *APIClient) doWebsocketRequest(ctx context.Context, jsonRPC2 JsonRPC2, ch chan<- interface{}) {
	config, err := configs.GetGeneralConfig()
	if err != nil {
		return
	}
	c, _, err := websocket.DefaultDialer.Dial(config.GetWebsocketString(), nil)
	if err != nil {
		log.Fatalf("function=APIClient.doWebsocketRequest, action=Websocket Dial, argslen=3, args=%v , %v , %v err=%s \n", ctx, jsonRPC2, ch, err.Error())
	}

	defer c.Close()
	if err := c.WriteJSON(&jsonRPC2); err != nil {
		log.Fatalf("function=APIClient.doWebsocketRequest, action=Write Json, argslen=3, args=%v , %v , %v err=%s \n", ctx, jsonRPC2, ch, err.Error())
	}
	c.SetWriteDeadline(time.Now().Add(10 * time.Second))

	if err != nil {
		log.Fatalf("function=APIClient,doWebsocketRequest, action=Get Config, argslen=3, args=%v , %v , %v err=%s \n", ctx, jsonRPC2, ch, err.Error())
	}
	retrymsec := config.Retrymsec()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			message := new(JsonRPC2)
			if err := c.ReadJSON(message); err != nil {
				if retrymsec > 0 {
					time.Sleep(time.Duration(retrymsec) * time.Millisecond)
				} else {
					log.Println(message)
					log.Fatalln("read:", err)
				}
			}

			if message.Method == "channelMessage" {
				switch params := message.Params.(type) {
				case map[string]interface{}:
          if v,ok := params["message"]; ok {
            ch <- v
          }
				}
			}
		}
	}
}

func (api *APIClient) GetPermissions() (permissions map[string]bool, err error) {
	url := bitflyerURL + "me/getpermissions"
	query := map[string]string{}
	permissions = make(map[string]bool, 0)

	// HTTP Public API
	permissions["getmarket"] = true
	permissions["getboard"] = true
	permissions["getticker"] = true
	permissions["getexecutions"] = true
	permissions["getboardstate"] = true
	permissions["gethealth"] = true
	permissions["getchats"] = true
	// HTTP Private API
	permissions["me/getpermissions"] = false
	permissions["me/getbalance"] = false
	permissions["me/getcollateral"] = false
	permissions["me/getcollateralaccounts"] = false
	permissions["me/getaddresses"] = false
	permissions["me/getcoinins"] = false
	permissions["me/getcoinouts"] = false
	permissions["me/getbankaccounts"] = false
	permissions["me/getdeposits"] = false
	permissions["me/whthdraw"] = false
	permissions["me/getwithdrawals"] = false
	permissions["me/sendchildorder"] = false
	permissions["me/cancelchildorder"] = false
	permissions["me/sendparentorder"] = false
	permissions["me/cancelparentorder"] = false
	permissions["me/cancelallchildorders"] = false
	permissions["me/getchildorders"] = false
	permissions["me/getparentorders"] = false
	permissions["me/getparentorder"] = false
	permissions["me/getexecutions"] = false
	permissions["me/getbalancehistory"] = false
	permissions["me/getpositions"] = false
	permissions["me/getcollateralhistory"] = false
	permissions["me/gettradingcommisson"] = false

	resp, err := api.doRequest("GET", url, query, nil)
	if err != nil {
		return permissions, err
	}

	ret := new(bitflyer.Permission)
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		log.Print(resp)
		log.Printf("error is %s", err)
		return permissions, err
	}

	for _, v := range ret.Value {
		slice := strings.Split(v, "/")
		key := strings.Join(slice[2:], "/")
		permissions[key] = true
	}
	return permissions, nil

}

func (api *APIClient) GetBoardState(productCode string) (boardState *bitflyer.BoardState, err error) {
	url := bitflyerURL + "getboardstate"
	query := map[string]string{}
	query["product_code"] = productCode

	resp, err := api.doRequest("GET", url, query, nil)
	if err != nil {
		return nil, err
	}
	boardState = new(bitflyer.BoardState)
	err = json.Unmarshal(resp, boardState)
	if err != nil {
		return nil, err
	}
	return boardState, nil
}

func (api *APIClient) GetHealth(productCode string) (health *bitflyer.Health, err error) {
	url := bitflyerURL + "gethealth"
	query := map[string]string{}
	query["product_code"] = productCode

	resp, err := api.doRequest("GET", url, query, nil)
	if err != nil {
		return nil, err
	}
	health = new(bitflyer.Health)
	err = json.Unmarshal(resp, health)
	if err != nil {
		return nil, err
	}
	return health, nil
}

func (api *APIClient) GetTicker(productCode string) (ticker *bitflyer.Ticker, err error) {
	url := bitflyerURL + "getticker"
	query := map[string]string{}
	query["product_code"] = productCode

	resp, err := api.doRequest("GET", url, query, nil)
	if err != nil {
		return nil, err
	}
	ticker = new(bitflyer.Ticker)
	err = json.Unmarshal(resp, ticker)
	if err != nil {
		return nil, err
	}
	return ticker, nil
}

func (api *APIClient) GetRealtimeTicker(ctx context.Context, ch chan<- bitflyer.Ticker, productCode string) {
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
			ticker := new(bitflyer.Ticker)
			marchalTick, err := json.Marshal(param)
			if err != nil {
				ticker.Message = err.Error()
				ch <- *ticker
				continue OUTER
			}
			if err := json.Unmarshal(marchalTick, &ticker); err != nil {
				ticker.Message = err.Error()
				ch <- *ticker
				continue OUTER
			}
			ch <- *ticker
		}

	}
}

func (api *APIClient) GetBoard(productCode string) (board *bitflyer.Board, err error) {
	url := bitflyerURL + "getboard"
	query := map[string]string{}
	query["product_code"] = productCode

	resp, err := api.doRequest("GET", url, query, nil)
	if err != nil {
		return nil, err
	}
	board = new(bitflyer.Board)
	err = json.Unmarshal(resp, board)
	if err != nil {
		return nil, err
	}
	return board, nil
}

func (api *APIClient) GetRealtimeBoard(ctx context.Context, ch chan<- bitflyer.Board, productCode string, isDiff bool) {
	jsonRPC2 := new(JsonRPC2)
	jsonRPC2.Version = "2.0"
	jsonRPC2.Method = "subscribe"

	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if isDiff {
		jsonRPC2.Params = SubscriveParams{Channel: fmt.Sprintf("lightning_board_%s", productCode)}
	} else {
		jsonRPC2.Params = SubscriveParams{Channel: fmt.Sprintf("lightning_board_snapshot_%s", productCode)}

	}

	var paramCh = make(chan interface{})
	go api.doWebsocketRequest(childctx, *jsonRPC2, paramCh)

OUTER:
	for {
		select {
		case <-ctx.Done():
			return

		default:
			param := <-paramCh
			marchalBoard, err := json.Marshal(param)
			if err != nil {
				log.Printf("error : %s", err)
				continue OUTER
			}
			board := new(bitflyer.Board)
			if err := json.Unmarshal(marchalBoard, &board); err != nil {
				log.Printf("error : %s", err)
				continue OUTER
			}
			ch <- *board
		}
	}
}

func (api *APIClient) GetExecutions(productCode string, beforeID, afterID string, count int) (executions []bitflyer.Execution, err error) {
	url := bitflyerURL + "getexecutions"
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

	executions = make([]bitflyer.Execution, 0)
	err = json.Unmarshal(resp, &executions)
	if err != nil {
		return nil, err
	}

	return executions, nil
}

func (api *APIClient) GetRealtimeExecutions(ctx context.Context, ch chan<- []bitflyer.Execution, productCode string) {
	jsonRPC2 := new(JsonRPC2)
	jsonRPC2.Version = "2.0"
	jsonRPC2.Method = "subscribe"

	childctx, cancel := context.WithCancel(ctx)
	defer cancel()
	jsonRPC2.Params = SubscriveParams{Channel: fmt.Sprintf("lightning_executions_%s", productCode)}

	var paramCh = make(chan interface{})
	go api.doWebsocketRequest(childctx, *jsonRPC2, paramCh)

OUTER:
	for {
		select {
		case <-ctx.Done():
			return

		default:
			param := <-paramCh
			marshalExecutions, err := json.Marshal(param)
			if err != nil {
				continue OUTER
			}
			executions := make([]bitflyer.Execution, 0)
			if err := json.Unmarshal(marshalExecutions, &executions); err != nil {
				continue OUTER
			}
			ch <- executions
		}

	}
}

func (api *APIClient) GetBalance() (balances []bitflyer.Balance, err error) {
	url := bitflyerURL + "me/getbalance"
	resp, err := api.doRequest("GET", url, map[string]string{}, nil)
	if err != nil {
		return nil, err
	}

	balances = make([]bitflyer.Balance, 0)
	err = json.Unmarshal(resp, &balances)
	if err != nil {
		return nil, err
	}

	return balances, nil
}
