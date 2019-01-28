package neoutils_test

import (
	"log"
	"math"
	"testing"
	"encoding/json"

	"github.com/o3labs/neo-utils/neoutils"
)

type parameterJSONArrayForm struct {
	A []parameterJSONForm `json:"array"`
}

type parameterJSONForm struct {
	T string `json:"type"`
	V interface{} `json:"value"`
}

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
	if wif == "" {
		log.Printf("No wif")
		return
	}

	gasPrice := int(500)
	gasLimit := int(20000)
	txid, err := neoutils.ClaimONG(endpoint, gasPrice, gasLimit, wif)
	if err != nil {
		log.Printf("err %v", err)
		return
	}
	log.Printf("tx id =%v", txid)
}

func TestBuildOntologyInvocation(t *testing.T) {
	wif := ""
	if wif == "" {
		log.Printf("No wif")
		return
	}

  account, _ := neoutils.GenerateFromWIF(wif)
  address := account.Address

	addr := parameterJSONForm{T: "Address", V: address}
  val := parameterJSONForm{T: "String", V: "Hi there"}

  jsondat := &parameterJSONArrayForm{A: []parameterJSONForm{addr, val}}
  argData, _ := json.Marshal(jsondat)
	argString := string(argData)

  gasPrice := int(500)
  gasLimit := int(20000)

  txData, err := neoutils.BuildOntologyInvocationTransaction("c168e0fb1a2bddcd385ad013c2c98358eca5d4dc", "put", argString, gasPrice, gasLimit, wif)
  if err != nil {
    log.Printf("Error creating invocation transaction: %s", err)
    t.Fail()
  } else {
    log.Printf("Raw transaction: %s", txData)
  }
}

func TestOntologyInvoke(t *testing.T) {
	wif := ""
	if wif == "" {
		log.Printf("No wif")
		return
	}

  account, _ := neoutils.GenerateFromWIF(wif)
  address := account.Address

	addr := parameterJSONForm{T: "Address", V: address}
  val := parameterJSONForm{T: "String", V: "Hi there"}

  jsondat := &parameterJSONArrayForm{A: []parameterJSONForm{addr, val}}
  argData, _ := json.Marshal(jsondat)
	argString := string(argData)

  gasPrice := int(500)
  gasLimit := int(20000)

	endpoint := "http://polaris2.ont.io:20336"

  txid, err := neoutils.OntologyInvoke(endpoint, "c168e0fb1a2bddcd385ad013c2c98358eca5d4dc", "put", argString, gasPrice, gasLimit, wif)
  if err != nil {
    log.Printf("Error creating invocation transaction: %s", err)
    t.Fail()
  } else {
		log.Printf("tx id = %s", txid)
  }
}
