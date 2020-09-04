package client

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/atoyr/goflyer/client/bitflyer"
)

func (api *APIClient) GetPermissions() (permissions map[string]bool, err error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	url, err := config.GetWebapiUrl("me/getpermissions")
	if err != nil {
		return nil, err
	}
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
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	url, err := config.GetWebapiUrl("getboardstate")
	if err != nil {
		return nil, err
	}
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
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	url, err := config.GetWebapiUrl("gethealth")
	if err != nil {
		return nil, err
	}
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
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	url, err := config.GetWebapiUrl("getticker")
	if err != nil {
		return nil, err
	}
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

func (api *APIClient) GetBoard(productCode string) (board *bitflyer.Board, err error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	url, err := config.GetWebapiUrl("getboard")
	if err != nil {
		return nil, err
	}
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

func (api *APIClient) GetExecutions(productCode string, beforeID, afterID string, count int) (executions []bitflyer.Execution, err error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	url, err := config.GetWebapiUrl("getexecutions")
	if err != nil {
		return nil, err
	}
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

func (api *APIClient) GetBalance() (balances []bitflyer.Balance, err error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	url, err := config.GetWebapiUrl("me/getbalance")
	if err != nil {
		return nil, err
	}
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
