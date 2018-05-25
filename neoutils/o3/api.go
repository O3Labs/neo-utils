package o3

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const apiEndpoint = "https://platform.o3.network/api"

type NEONetWork string

var NEOMainNet = "main"
var NEOTestNet = "test"

type O3APIInterface interface {
	GetNEOUTXO(address string) UTXOResponse
	GetNEOClimableGAS(address string) ClaimableGASResponse
}

type O3Client struct {
	APIBaseEndpoint url.URL
	neoNetwork      string
}

func DefaultO3APIClient() *O3Client {
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		return nil
	}
	return &O3Client{APIBaseEndpoint: *u}
}

func APIClientWithNEOTestnet() *O3Client {
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		return nil
	}
	return &O3Client{APIBaseEndpoint: *u, neoNetwork: "test"}
}

func APIClientWithNEOPrivateNet() *O3Client {
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		return nil
	}
	return &O3Client{APIBaseEndpoint: *u, neoNetwork: "private"}
}

//make sure all method interface is implemented
var _ O3APIInterface = (*O3Client)(nil)

func (n *O3Client) makeGETRequest(endpoint string, out interface{}) error {

	fullEndpointString := fmt.Sprintf("%v%v", n.APIBaseEndpoint.String(), endpoint)
	fullEndpoint, _ := url.Parse(fullEndpointString)

	if n.neoNetwork == "test" || n.neoNetwork == "private" {
		q := fullEndpoint.Query()
		q.Set("network", n.neoNetwork)
		fullEndpoint.RawQuery = q.Encode()
	}

	log.Printf("%v", fullEndpoint.String())

	req, err := http.NewRequest("GET", fullEndpoint.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Set("Connection", "close")
	req.Close = true
	res, err := http.DefaultClient.Do(req)
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

func (o *O3Client) GetNEOUTXO(address string) UTXOResponse {
	response := UTXOResponse{}
	err := o.makeGETRequest(fmt.Sprintf("/v1/neo/%v/utxo", address), &response)
	if err != nil {
		return response
	}
	return response
}
func (o *O3Client) GetNEOClimableGAS(address string) ClaimableGASResponse {
	response := ClaimableGASResponse{}

	err := o.makeGETRequest(fmt.Sprintf("/v1/neo/%v/claimablegas", address), &response)
	if err != nil {
		return response
	}
	return response
}
