package neoutils_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
)

func TestMintTokensFromMobile(t *testing.T) {
	scriptHash := "a6f331f6d7d8b331a95ef8513c573e938c073c6b"

	//both of these are whitelisted

	// wif := "KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr" // contract owner AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y
	// wif := "Kz3dZoCXU8SsmE67GLoGZKaghD1bG1kbePY72LVKpuchMqmRwKer" // Adm9ER3UwdJfimFtFhHq1L5MQ5gxLLTUes
	wif := "L5h6cTh45egotcxFZ2rkF1gv7rLxx9rScfuja9kEVEE9mEj9Uwtv" //AQaZPqcv9Kg2x1eSrF8UBYXLK4WQoTSLH5
	wallet, _ := neoutils.GenerateFromWIF(wif)
	log.Printf("address = %v\n address hash = %v %x", wallet.Address, neoutils.NEOAddressToScriptHash(wallet.Address), wallet.HashedSignature)

	sendingAssetID := "c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b" //NEO
	amount := float64(1)
	remark := "o3x"
	utxoEndpoint := "http://localhost:5000/"
	networkFeeAmountInGAS := float64(0.0011)

	tx, err := neoutils.MintTokensRawTransactionMobile(utxoEndpoint, scriptHash, wif, sendingAssetID, amount, remark, networkFeeAmountInGAS)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	log.Printf("txID =%v", tx.TXID)
	log.Printf("tx = %x", tx.Data)
}
