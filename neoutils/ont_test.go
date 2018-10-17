package neoutils_test

import (
	"log"
	"math"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
)

func TestONTTransfer(t *testing.T) {

	for i := 1; i <= 100; i++ {

		endpoint := "http://polaris1.ont.io:20336"
		//pass := ""
		//wif, _ := neoutils.NEP2Decrypt("", pass)
		wif := ""
		asset := "ong"
		to := "AcWfHYbPDt4ysz7s5WQtkGvcFsfTsM6anm"

		amount := float64(float64(i) / math.Pow10(9))
		gasPrice := int(500)
		gasLimit := int(20000)
		txid, err := neoutils.OntologyTransfer(endpoint, gasPrice, gasLimit, wif, asset, to, amount)
		if err != nil {
			log.Printf("err %v", err)
			return
		}
		log.Printf("tx id =%v", txid)
	}
}

func TestClaimONG(t *testing.T) {
	endpoint := "http://dappnode2.ont.io:20336"
	wif, _ := neoutils.NEP2Decrypt("", "")

	gasPrice := int(500)
	gasLimit := int(20000)
	txid, err := neoutils.ClaimONG(endpoint, gasPrice, gasLimit, wif)
	if err != nil {
		log.Printf("err %v", err)
		return
	}
	log.Printf("tx id =%v", txid)
}
