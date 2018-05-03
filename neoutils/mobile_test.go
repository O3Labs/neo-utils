package neoutils_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func TestMintTokensFromMobile(t *testing.T) {
	scriptHash := "0xf820184470f1a5f38fa4ecc8db746336b371bda5"

	//both of these are whitelisted

	//ALL PRIVATE NET TEST ADDRESSES
	// wif := "KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr" // contract owner AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y
	// wif := "L2W3eBvPYMUaxDZGEb395HZf26tLPZgU5qN351HpyVSAG1DWgDtx" // AQmq2yU7DupE4VddmEoweKiJFyGhAAEZeH
	wif := "L5h6cTh45egotcxFZ2rkF1gv7rLxx9rScfuja9kEVEE9mEj9Uwtv" //AQaZPqcv9Kg2x1eSrF8UBYXLK4WQoTSLH5
	//this addresss is not whitelisted
	// wif := "L5gmcoaetU6YGSzg4wNqvKBEEAfwCAxWseuL3pxvLvEMZB9WyUYp"
	wallet, _ := neoutils.GenerateFromWIF(wif)
	log.Printf("address = %v\n address hash = %x", wallet.Address, wallet.HashedSignature)

	neo := string(smartcontract.NEO)
	// gas := string(smartcontract.GAS)
	amount := float64(2)
	remark := "o3x"
	utxoEndpoint := "http://localhost:5000/"
	networkFeeAmountInGAS := float64(0.0011)

	tx, err := neoutils.MintTokensRawTransactionMobile(utxoEndpoint, scriptHash, wif, neo, amount, remark, networkFeeAmountInGAS)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	log.Printf("txID =%v", tx.TXID)
	log.Printf("tx = %x", tx.Data)
}
