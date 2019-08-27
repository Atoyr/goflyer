package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strconv"
	"time"


	"github.com/atoyr/goflyer/api/model"
	"github.com/gorilla/websocket"

)

type APIClient struct {
	key        string
	secret     string
	httpClient *http.Client
	config     *apiConfig
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
	endpoint ,err :=  api.config.GetEndpoint(urlPath)
	if err != nil {
		return nil, err
	}

	req , err := http.NewRequest(method,endpoint,bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for k,v := range query {
		q.Add(k,v)
	}
	req.URL.RawQuery = q.Encode()
for k,v := range api.header(method, req.URL.RequestURI(), body) {
		req.Header.Add(k,v)
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

func (api *APIClient) GetTicker() (ticker *model.Ticker, err error){
	url:= "getticker"
	resp, err := api.doRequest("GET", url, map[string]string{}, nil)
	if err != nil {
		return nil , err
	}
	err = json.Unmarshal(resp, ticker)
	if err != nil {
		return nil , err
	}
	return ticker, nil
}

func (api *APIClient) GetRealtimeTicker(){
	//c , _ ,err := websocket
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

