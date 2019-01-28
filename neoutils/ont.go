package neoutils

import (
	"fmt"
	"log"

	"github.com/o3labs/ont-mobile/ontmobile"
)

func OntologyTransfer(endpoint string, gasPrice int, gasLimit int, wif string, asset string, to string, amount float64) (string, error) {
	raw, err := ontmobile.Transfer(uint(gasPrice), uint(gasLimit), wif, asset, to, amount)
	if err != nil {
		return "", err
	}
	log.Printf("raw = %x", raw.Data)
	txid, err := ontmobile.SendRawTransaction(endpoint, fmt.Sprintf("%x", raw.Data))
	if err != nil {
		return "", err
	}

	return txid, nil
}

func ClaimONG(endpoint string, gasPrice int, gasLimit int, wif string) (string, error) {
	raw, err := ontmobile.WithdrawONG(uint(gasPrice), uint(gasLimit), endpoint, wif)
	if err != nil {
		return "", err
	}

	txid, err := ontmobile.SendRawTransaction(endpoint, fmt.Sprintf("%x", raw.Data))
	if err != nil {
		return "", err
	}

	return txid, nil
}

func BuildOntologyInvocationTransaction(contractHex string, operation string, args []ontmobile.Parameter, gasPrice uint, gasLimit uint, wif string) (string, error) {
	raw, err := ontmobile.BuildInvocationTransaction(contractHex, operation, args, gasPrice, gasLimit, wif)
	if err != nil {
		return "", err
	}

	return raw, nil
}

func OntologyInvoke(endpoint string, contractHex string, operation string, args []Parameter, gasPrice uint, gasLimit uint, wif string) (string, error) {
	raw, err := ontmobile.BuildInvocationTransaction(contractHex, operation, args, gasPrice, gasLimit, wif)
	if err != nil {
		return "", err
	}

	txid, err := ontmobile.SendRawTransaction(endpoint, raw)
	if err != nil {
		return "", err
	}

	return txid, nil
}

type Parameter = ontmobile.Parameter
type ParameterType = ontmobile.ParameterType

const (
  Address     ParameterType = 0
  String      ParameterType = 1
  Integer     ParameterType = 2
  Fixed8      ParameterType = 3
  Array       ParameterType = 4
)
