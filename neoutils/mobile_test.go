package neoutils_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
)

func TestMintTokensFromMobile(t *testing.T) {
	scriptHash := "b2eb148d3783f60e678e35f2c496de1a2a7ead93"

	//both of these are whitelisted
	// wif := "Kz3dZoCXU8SsmE67GLoGZKaghD1bG1kbePY72LVKpuchMqmRwKer"
	// wif := "KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr"
	//this one is not whitelisted
	wif := "L5h6cTh45egotcxFZ2rkF1gv7rLxx9rScfuja9kEVEE9mEj9Uwtv"
	wallet, _ := neoutils.GenerateFromWIF(wif)
	log.Printf("address = %v\n address hash = %v", wallet.Address, neoutils.NEOAddressToScriptHash(wallet.Address))

	sendingAssetID := "c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b" //NEO
	amount := float64(1)
	remark := "o3x"
	utxoEndpoint := "http://localhost:5000/"
	networkFeeAmountInGAS := float64(0.0012)

	tx, err := neoutils.MintTokensRawTransactionMobile(utxoEndpoint, scriptHash, wif, sendingAssetID, amount, remark, networkFeeAmountInGAS)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	log.Printf("txID =%v", tx.TXID)
	log.Printf("tx = %x", tx.Data)
}
