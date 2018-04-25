package neoutils_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func TestMintTokensFromMobile(t *testing.T) {
	scriptHash := "a6f331f6d7d8b331a95ef8513c573e938c073c6b"

	//both of these are whitelisted

	//ALL PRIVATE NET TEST ADDRESSES
	// wif := "KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr" // contract owner AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y
	wif := "L2W3eBvPYMUaxDZGEb395HZf26tLPZgU5qN351HpyVSAG1DWgDtx" // Adm9ER3UwdJfimFtFhHq1L5MQ5gxLLTUes
	// wif := "L5h6cTh45egotcxFZ2rkF1gv7rLxx9rScfuja9kEVEE9mEj9Uwtv" //AQaZPqcv9Kg2x1eSrF8UBYXLK4WQoTSLH5
	//this addresss is not whitelisted
	// wif := "L5gmcoaetU6YGSzg4wNqvKBEEAfwCAxWseuL3pxvLvEMZB9WyUYp"
	wallet, _ := neoutils.GenerateFromWIF(wif)
	log.Printf("address = %v\n address hash = %v %x", wallet.Address, neoutils.NEOAddressToScriptHash(wallet.Address), wallet.HashedSignature)

	neo := string(smartcontract.NEO)
	// gas := string(smartcontract.GAS)
	amount := float64(1)
	remark := "o3x"
	utxoEndpoint := "http://localhost:5000/"
	networkFeeAmountInGAS := float64(0.001)

	tx, err := neoutils.MintTokensRawTransactionMobile(utxoEndpoint, scriptHash, wif, neo, amount, remark, networkFeeAmountInGAS)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	log.Printf("txID =%v", tx.TXID)
	log.Printf("tx = %x", tx.Data)
}
