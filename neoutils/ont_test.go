package neoutils_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
)

func TestONTTransfer(t *testing.T) {
	endpoint := "http://polaris2.ont.io:20336"
	wif := ""
	asset := "ong"
	to := "AafQxV6wQhtGYGYFboEyBjw3eMYNtBFW8J"
	amount := float64(1000)
	txid, err := neoutils.OntologyTransfer(endpoint, wif, asset, to, amount)
	if err != nil {
		log.Printf("err %v", err)
		return
	}
	log.Printf("tx id =%v", txid)
}
