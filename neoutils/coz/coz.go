package coz

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type CozClientInterface interface {
	GetUnspentByAddress(address string) (*UnspentBalance, error)
	GetClaims(address string) (*ClaimResponse, error)
}

type CozClient struct {
	Endpoint   url.URL
	httpClient *http.Client
}

//make sure all method interface is implemented
var _ CozClientInterface = (*CozClient)(nil)

func NewClient(endpoint string) *CozClient {
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

	return &CozClient{Endpoint: *u, httpClient: netClient}
}

func (c *CozClient) GetUnspentByAddress(address string) (*UnspentBalance, error) {
	req, _ := http.NewRequest("GET", c.Endpoint.String()+"/v2/address/balance/"+address, nil)

	res, _ := c.httpClient.Do(req)

	unspent := UnspentBalance{}
	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(&unspent)
	if err != nil {
		return nil, err
	}
	return &unspent, nil
}

func (c *CozClient) GetClaims(address string) (*ClaimResponse, error) {
	req, _ := http.NewRequest("GET", c.Endpoint.String()+"/v2/address/claims/"+address, nil)
	req.Header.Set("Connection", "close")
	req.Close = true
	res, _ := http.DefaultClient.Do(req)

	response := ClaimResponse{}
	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
