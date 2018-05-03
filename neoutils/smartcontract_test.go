package neoutils_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func TestInvalidSmartContractStruct(t *testing.T) {
	sc := neoutils.UseSmartContract("ce575ae1bb6153330d2")
	if sc != nil {
		t.Fail()
		return
	}
	log.Printf("%v", sc)
}

func TestUseSmartContractWithEmptyScripthash(t *testing.T) {
	sc := neoutils.UseSmartContract("")
	if sc != nil {
		t.Fail()
		return
	}
	log.Printf("%v", sc)
}

func TestSmartContractStruct(t *testing.T) {

	sc := neoutils.UseSmartContract("ce575ae1bb6153330d20c560acb434dc5755241b")
	if sc == nil {
		t.Fail()
		return
	}
	log.Printf("%v", sc)
}

func UTXODataForSmartContract() smartcontract.Unspent {

	gasTX1 := smartcontract.UTXO{
		Index: 0,
		TXID:  "307d756074d9ee11220ccebf003bedb99f9b1a54e194a25e6ea5df1a7b2de84b",
		Value: float64(713399700000) / float64(100000000),
	}

	gasBalance := smartcontract.Balance{
		Amount: float64(713399700000) / float64(100000000),
		UTXOs:  []smartcontract.UTXO{gasTX1},
	}

	neoTX1 := smartcontract.UTXO{
		Index: 0,
		TXID:  "e8b8bf4f98490368fc1caa86f8646e7383bb52751ffc3a1a7e296d715c4382ed",
		Value: float64(10000000000000000) / float64(100000000),
	}

	neoBalance := smartcontract.Balance{
		Amount: float64(10000000000000000) / float64(100000000),
		UTXOs:  []smartcontract.UTXO{neoTX1},
	}

	unspent := smartcontract.Unspent{
		Assets: map[smartcontract.NativeAsset]*smartcontract.Balance{},
	}
	unspent.Assets[smartcontract.NEO] = &neoBalance
	unspent.Assets[smartcontract.GAS] = &gasBalance
	return unspent
}

func TestInvokeFunctionRawTransaction(t *testing.T) {
	var validSmartContract = neoutils.UseSmartContract("b7c1f850a025e34455e7e98c588c784385077fb1")

	wif := "KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr"
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
	}

	to := smartcontract.ParseNEOAddress("AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn")
	if to == nil {
		//invalid neo address
		t.Fail()
		return
	}
	args := []interface{}{to, 1}
	unspent := UTXODataForSmartContract()

	transactionID := "thisisauniquetoken_from_stripe"
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	attributes[smartcontract.Remark1] = []byte(transactionID)
	tx, err := validSmartContract.GenerateInvokeFunctionRawTransaction(*privateNetwallet, unspent, attributes, "mintTokensTo", args)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("%x", tx)
}

func TestGenerateInvokeTransferNEP5Token(t *testing.T) {
	var validSmartContract = neoutils.UseSmartContract("b7c1f850a025e34455e7e98c588c784385077fb1")

	wif := "KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr"
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
	}

	from := smartcontract.ParseNEOAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
	if from == nil {
		//invalid neo address
		t.Fail()
		return
	}

	to := smartcontract.ParseNEOAddress("AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn")
	if to == nil {
		//invalid neo address
		t.Fail()
		return
	}
	numberOfTokens := 1
	args := []interface{}{from, to, numberOfTokens}
	unspent := UTXODataForSmartContract()

	remark := "this is a remark data in attribute"
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	attributes[smartcontract.Remark1] = []byte(remark)
	tx, err := validSmartContract.GenerateInvokeFunctionRawTransaction(*privateNetwallet, unspent, attributes, "transfer", args)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("%x", tx)
}
