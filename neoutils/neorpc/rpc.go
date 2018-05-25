package neorpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/o3labs/neo-utils/neoutils/btckey"
)

type NEORPCInterface interface {
	GetContractState(scripthash string) GetContractStateResponse
	SendRawTransaction(rawTransactionInHex string) SendRawTransactionResponse
	GetRawTransaction(txID string) GetRawTransactionResponse
	makeRequest(method string, params []interface{}, out interface{}) error
	GetBlockCount() GetBlockCountResponse
	GetBlock(blockHash string) GetBlockResponse
	GetBlockByIndex(index int) GetBlockResponse
	GetAccountState(address string) GetAccountStateResponse
	InvokeScript(scriptInHex string) InvokeScriptResponse
	GetTokenBalance(tokenHash string, adddress string) TokenBalanceResponse
}

type NEORPCClient struct {
	Endpoint   url.URL
	httpClient *http.Client
}

//make sure all method interface is implemented
var _ NEORPCInterface = (*NEORPCClient)(nil)

func NewClient(endpoint string) *NEORPCClient {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil
	}
	// var netTransport = &http.Transport{
	// 	Dial: (&net.Dialer{
	// 		Timeout: 8 * time.Second,
	// 	}).Dial,
	// 	TLSHandshakeTimeout: 8 * time.Second,
	// }

	var netClient = &http.Client{
		Timeout: time.Second * 60,
		// Transport: netTransport,
	}

	return &NEORPCClient{Endpoint: *u, httpClient: netClient}
}

func (n *NEORPCClient) makeRequest(method string, params []interface{}, out interface{}) error {
	request := NewRequest(method, params)

	jsonValue, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", n.Endpoint.String(), bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Set("Connection", "close")
	req.Close = true
	res, err := n.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return err
	}

	return nil
}

func (n *NEORPCClient) GetContractState(scripthash string) GetContractStateResponse {
	response := GetContractStateResponse{}
	params := []interface{}{scripthash, 1}
	err := n.makeRequest("getcontractstate", params, &response)
	if err != nil {
		return response
	}
	return response
}

func (n *NEORPCClient) SendRawTransaction(rawTransactionInHex string) SendRawTransactionResponse {
	response := SendRawTransactionResponse{}
	params := []interface{}{rawTransactionInHex, 1}
	err := n.makeRequest("sendrawtransaction", params, &response)
	if err != nil {
		return response
	}
	return response
}

func (n *NEORPCClient) GetRawTransaction(txID string) GetRawTransactionResponse {
	response := GetRawTransactionResponse{}
	params := []interface{}{txID, 1}
	err := n.makeRequest("getrawtransaction", params, &response)
	if err != nil {
		return response
	}
	return response
}

func (n *NEORPCClient) GetBlock(blockHash string) GetBlockResponse {
	response := GetBlockResponse{}
	params := []interface{}{blockHash, 1}
	err := n.makeRequest("getblock", params, &response)
	if err != nil {
		return response
	}
	return response
}

func (n *NEORPCClient) GetBlockByIndex(index int) GetBlockResponse {
	response := GetBlockResponse{}
	params := []interface{}{index, 1}
	err := n.makeRequest("getblock", params, &response)
	if err != nil {
		return response
	}
	return response
}

func (n *NEORPCClient) GetBlockCount() GetBlockCountResponse {
	response := GetBlockCountResponse{}
	params := []interface{}{}
	err := n.makeRequest("getblockcount", params, &response)
	if err != nil {
		return response
	}
	return response
}

func (n *NEORPCClient) GetAccountState(address string) GetAccountStateResponse {
	response := GetAccountStateResponse{}
	params := []interface{}{address, 1}
	err := n.makeRequest("getaccountstate", params, &response)
	if err != nil {
		return response
	}
	return response
}

func (n *NEORPCClient) GetTokenBalance(tokenHash string, neoAddress string) TokenBalanceResponse {
	response := TokenBalanceResponse{}
	args := []interface{}{}

	v, b, _ := btckey.B58checkdecode(neoAddress)
	if v != 0x17 {
		return TokenBalanceResponse{}
	}
	adddressScriptHash := fmt.Sprintf("%x", b)
	log.Printf("adddressScriptHash %v", adddressScriptHash)
	input := NewInvokeFunctionStackByteArray(adddressScriptHash)
	args = append(args, input)

	params := []interface{}{tokenHash, "balanceOf", args}
	err := n.makeRequest("invokefunction", params, &response)
	if err != nil {
		return response
	}
	return response
}

func (n *NEORPCClient) InvokeScript(scriptInHex string) InvokeScriptResponse {
	response := InvokeScriptResponse{}
	params := []interface{}{scriptInHex, 1}
	err := n.makeRequest("invokescript", params, &response)
	if err != nil {
		return response
	}
	return response
}
