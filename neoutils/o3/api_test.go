package o3_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils/o3"
)

func TestO3NEOUTXO(t *testing.T) {
	client := o3.DefaultO3APIClient()
	response := client.GetNEOUTXO("ANk325vGG5kcc6Dcnk6zkoEBHY4E6es2nY")
	if response.Code != 200 {
		t.Fail()
	}
	log.Printf("%+v", response)
}

func TestO3NEOUTXOTestnet(t *testing.T) {
	client := o3.APIClientWithNEOTestnet()
	response := client.GetNEOUTXO("ANk325vGG5kcc6Dcnk6zkoEBHY4E6es2nY")
	if response.Code != 200 {
		t.Fail()
	}
	log.Printf("%+v", response)
}

func TestO3GetClaimableGAS(t *testing.T) {
	client := o3.DefaultO3APIClient()
	response := client.GetNEOClimableGAS("ANk325vGG5kcc6Dcnk6zkoEBHY4E6es2nY")
	if response.Code != 200 {
		t.Fail()
	}
	log.Printf("%+v", response.Result.Data.Gas)
	for _, v := range response.Result.Data.Claims {
		log.Printf("%+v", v)
	}
}
