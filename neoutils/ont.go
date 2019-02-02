package neoutils

import (
	"fmt"
	"log"

	"github.com/o3labs/ont-mobile/ontmobile"
	"github.com/o3labs/ont-mobile/ontmobile/ontrpc"
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
	raw, err := ontmobile.BuildInvocationTransaction(contract, method, args, uint(gasPrice), uint(gasLimit), wif)
	if err != nil {
		return "", err
	}

	return raw, nil
}

// OntologyInvoke : Invoke a neovm contract in Ontology
func OntologyInvoke(endpoint string, contract string, method string, args string, gasPrice int, gasLimit int, wif string) (string, error) {
	raw, err := ontmobile.BuildInvocationTransaction(contract, method, args, uint(gasPrice), uint(gasLimit), wif)
	if err != nil {
		return "", err
	}

	txid, err := ontmobile.SendRawTransaction(endpoint, raw)
	if err != nil {
		return "", err
	}

	return txid, nil
}

func ontGetBlockCount(endpoint string) (ontrpc.GetBlockCountResponse, error) {
	client := ontrpc.NewRPCClient(endpoint)
	response, err := client.GetBlockCount()
	if err != nil {
		return response, err
	}
	return response, nil
}

func ontGetBalance(endpoint string, address string) (ontrpc.GetBalanceResponse, error) {
	client := ontrpc.NewRPCClient(endpoint)
	response, err := client.GetBalance(address)
	if err != nil {
		return response, err
	}
	return response, nil
}

func ontGetSmartCodeEvent(endpoint string, txHash string) (ontrpc.GetSmartCodeEventResponse, error) {
	client := ontrpc.NewRPCClient(endpoint)
	response, err := client.GetSmartCodeEvent(txHash)
	if err != nil {
		return response, err
	}
	return response, nil
}

func ontSendRawTransaction(endpoint string, raw string) (string, error) {
	txid, err := ontmobile.SendRawTransaction(endpoint, raw)
	if err != nil {
		return "", err
	}

	return txid, nil
}
