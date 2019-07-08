package neoutils

import (
	"encoding/json"
	"fmt"
	"log"
	"math"

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

func BuildOntologyInvocationTransaction(contract string, method string, args string, gasPrice int, gasLimit int, wif string) (string, error) {
	raw, err := ontmobile.BuildInvocationTransaction(contract, method, args, uint(gasPrice), uint(gasLimit), wif, "")
	if err != nil {
		return "", err
	}

	return raw, nil
}

// OntologyInvoke : Invoke a neovm contract in Ontology
func OntologyInvoke(endpoint string, contract string, method string, args string, gasPrice int, gasLimit int, wif string) (string, error) {
	raw, err := ontmobile.BuildInvocationTransaction(contract, method, args, uint(gasPrice), uint(gasLimit), wif, "")
	if err != nil {
		return "", err
	}

	txid, err := ontmobile.SendRawTransaction(endpoint, raw)
	if err != nil {
		return "", err
	}

	return txid, nil
}

func OEP4Transfer(endpoint string, contract string, fromAddress string, toAddress string, amount float64, tokenDecimals int, gasPrice int, gasLimit int, wif string) (string, error) {

	transferringAmount := uint(ontmobile.RoundFixed(float64(amount), tokenDecimals) * float64(math.Pow10(tokenDecimals)))
	payer := fromAddress
	fromAddressParam := ontmobile.ParameterJSONForm{T: "Address", V: fromAddress}

	toAddressParam := ontmobile.ParameterJSONForm{T: "Address", V: toAddress}
	amountParam := ontmobile.ParameterJSONForm{T: "Integer", V: transferringAmount}

	jsonData := &ontmobile.ParameterJSONArrayForm{A: []ontmobile.ParameterJSONForm{fromAddressParam,
		toAddressParam,
		amountParam}}

	argData, _ := json.Marshal(jsonData)
	argString := string(argData)

	raw, err := ontmobile.BuildInvocationTransaction(contract, "transfer", argString, uint(gasPrice), uint(gasLimit), wif, payer)
	if err != nil {
		return "", err
	}

	txid, err := ontmobile.SendRawTransaction(endpoint, raw)
	if err != nil {
		return "", err
	}

	return txid, nil
}
