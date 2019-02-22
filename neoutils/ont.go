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

func BuildOntologyInvocationTransaction(contract string, method string, args string, gasPrice int, gasLimit int, wif string, payer string) (string, error) {
	raw, err := ontmobile.BuildInvocationTransaction(contract, method, args, uint(gasPrice), uint(gasLimit), wif, payer)
	if err != nil {
		return "", err
	}

	return raw, nil
}

// OntologyInvoke : Invoke a neovm contract in Ontology
func OntologyInvoke(endpoint string, contract string, method string, args string, gasPrice int, gasLimit int, wif string, payer string) (string, error) {
	raw, err := ontmobile.BuildInvocationTransaction(contract, method, args, uint(gasPrice), uint(gasLimit), wif, payer)
	if err != nil {
		return "", err
	}

	txid, err := ontmobile.SendRawTransaction(endpoint, raw)
	if err != nil {
		return "", err
	}

	return txid, nil
}

func OntologyGetBlockCount(endpoint string) (int, error) {
	client := ontrpc.NewRPCClient(endpoint)
	response, err := client.GetBlockCount()
	if err != nil {
		return response.Result, err
	}
	return response.Result, nil
}

type OntologyBalances struct {
	Ont string
	Ong string
}

func OntologyGetBalance(endpoint string, address string) (*OntologyBalances, error) {
	client := ontrpc.NewRPCClient(endpoint)
	response, err := client.GetBalance(address)
	if err != nil {
		return nil, err
	}

	balances := &OntologyBalances{
		Ont: response.Result.Ont,
		Ong: response.Result.Ong,
	}

	return balances, nil
}

type SmartCodeEvent struct {
	TxHash      string
	State       int
	GasConsumed int
}

func OntologyGetSmartCodeEvent(endpoint string, txHash string) (*SmartCodeEvent, error) {
	client := ontrpc.NewRPCClient(endpoint)
	response, err := client.GetSmartCodeEvent(txHash)
	if err != nil {
		return nil, err
	}

	event := &SmartCodeEvent{
		TxHash:      response.Result.TxHash,
		State:       response.Result.State,
		GasConsumed: response.Result.GasConsumed,
	}

	return event, nil
}

func OntologySendRawTransaction(endpoint string, raw string) (string, error) {
	txid, err := ontmobile.SendRawTransaction(endpoint, raw)
	if err != nil {
		return "", err
	}

	return txid, nil
}

func OntologySendPreExecRawTransaction(endpoint string, raw string) (string, error) {
	result, err := ontmobile.SendPreExecRawTransaction(endpoint, raw)
	if err != nil {
		return "", err
	}

	return result, nil
}

func OntologyGetUnboundONG(endpoint string, address string) (string, error) {
	client := ontrpc.NewRPCClient(endpoint)
	response, err := client.GetUnboundONG(address)
	if err != nil {
		return "", err
	}

	return response.Result, nil
}

func OntologyGetStorage(endpoint string, scriptHash string, key string) (string, error) {
	client := ontrpc.NewRPCClient(endpoint)
	response, err := client.GetStorage(scriptHash, key)
	if err != nil {
		return response.Result, err
	}
	return response.Result, nil
}

func OntologyGetRawTransaction(endpoint string, txID string) (string, error) {
	client := ontrpc.NewRPCClient(endpoint)
	response, err := client.GetRawTransaction(txID)
	if err != nil {
		return response.Result, err
	}
	return response.Result, nil
}

func OntologyGetBlockWithHash(endpoint string, blockHash string) (string, error) {
	client := ontrpc.NewRPCClient(endpoint)
	response, err := client.GetBlockWithHash(blockHash)
	if err != nil {
		return response.Result, err
	}
	return response.Result, nil
}

func OntologyGetBlockWithHeight(endpoint string, blockHeight int) (string, error) {
	client := ontrpc.NewRPCClient(endpoint)
	response, err := client.GetBlockWithHeight(blockHeight)
	if err != nil {
		return response.Result, err
	}
	return response.Result, nil
}

type ONTAccount struct {
	Address    string //base58
	WIF        string
	PrivateKey []byte
	PublicKey  []byte
}

func ONTCreateAccount() *ONTAccount {
	acc := ontmobile.NewONTAccount()
	if acc == nil {
		return nil
	}
	return &ONTAccount{
		Address:    acc.Address,
		WIF:        acc.WIF,
		PrivateKey: acc.PrivateKey,
		PublicKey:  acc.PublicKey,
	}
}

func ONTAccountFromPrivateKey(privateKeyBytes []byte) *ONTAccount {
	acc := ontmobile.ONTAccountWithPrivateKey(privateKeyBytes)
	if acc == nil {
		return nil
	}
	return &ONTAccount{
		Address:    acc.Address,
		WIF:        acc.WIF,
		PrivateKey: acc.PrivateKey,
		PublicKey:  acc.PublicKey,
	}
}

func ONTAccountFromWIF(wif string) *ONTAccount {
	acc := ontmobile.ONTAccountWithWIF(wif)
	if acc == nil {
		return nil
	}
	return &ONTAccount{
		Address:    acc.Address,
		WIF:        acc.WIF,
		PrivateKey: acc.PrivateKey,
		PublicKey:  acc.PublicKey,
	}
}

func OntologyMakeRegister(gasPrice int, gasLimit int, ontidWif string, payerWif string) (string, error) {
	raw, err := ontmobile.MakeRegister(uint(gasPrice), uint(gasLimit), ontidWif, payerWif)
	if err != nil {
		return "", err
	}
	return raw, nil
}

func OntologyBuildGetDDO(ontid string) (string, error) {
	raw, err := ontmobile.BuildGetDDO(ontid)
	if err != nil {
		return "", err
	}
	return raw, nil
}

func ONTAddressFromPublicKey(publicKey []byte) string {
	return ontmobile.ONTAddressFromPublicKey(publicKey)
}
