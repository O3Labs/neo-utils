package neorpc

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type NEORPCInterface interface {
	GetContractState(scripthash string) GetContractStateResponse
	SendRawTransaction(rawTransactionInHex string) SendRawTransactionResponse
	makeRequest()
}

type NEORPC struct {
	Endpoint url.URL
}

func NewNEORPC(endpoint string) *NEORPC {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil
	}
	return &NEORPC{Endpoint: *u}
}

func (n *NEORPC) makeRequest(method string, params []interface{}, out interface{}) error {
	request := NewRequest(method, params)

	jsonValue, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", n.Endpoint.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("error request %v", err)
		return err
	}
	req.Header.Add("content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		log.Printf("error %v", err)
		return err
	}

	return nil

}
func (n *NEORPC) GetContractState(scripthash string) GetContractStateResponse {
	response := GetContractStateResponse{}
	params := []interface{}{scripthash, 1}
	err := n.makeRequest("getcontractstate", params, &response)
	if err != nil {
		return response
	}
	return response
}

func (n *NEORPC) SendRawTransaction(rawTransactionInHex string) SendRawTransactionResponse {
	response := SendRawTransactionResponse{}
	params := []interface{}{rawTransactionInHex, 1}
	err := n.makeRequest("sendrawtransaction", params, &response)
	if err != nil {
		return response
	}
	return response
}
