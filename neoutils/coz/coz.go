package coz

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type CozClientInterface interface {
	GetUnspentByAddress(address string) (*UnspentBalance, error)
}

type CozClient struct {
	Endpoint url.URL
}

//make sure all method interface is implemented
var _ CozClientInterface = (*CozClient)(nil)

func NewClient(endpoint string) *CozClient {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil
	}
	return &CozClient{Endpoint: *u}
}

func (c *CozClient) GetUnspentByAddress(address string) (*UnspentBalance, error) {
	req, _ := http.NewRequest("GET", c.Endpoint.String()+"/v2/address/balance/"+address, nil)

	res, _ := http.DefaultClient.Do(req)

	unspent := UnspentBalance{}
	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(&unspent)
	if err != nil {
		return nil, err
	}
	return &unspent, nil
}
