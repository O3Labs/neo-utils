package coz_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils/coz"
)

func TestGetUnspent(t *testing.T) {
	client := coz.NewClient("http://127.0.0.1:5000")
	unspent, err := client.GetUnspentByAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	log.Printf("%+v", unspent)
	for _, v := range unspent.GAS.Unspent {
		log.Printf("%.8f", v.Value)
	}
}

func TestGetClaims(t *testing.T) {
	client := coz.NewClient("http://api.wallet.cityofzion.io/")
	response, err := client.GetClaims("AJaohf6jmPyPcSF1JQLjC1RspT7F74mhBP")
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	log.Printf("%+v", response)
}
