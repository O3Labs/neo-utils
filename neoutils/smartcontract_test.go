package neoutils_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

var validSmartContract = neoutils.NewSmartContract("ce575ae1bb6153330d20c560acb434dc5755241b")

func TestInvalidSmartContractStruct(t *testing.T) {
	sc := neoutils.NewSmartContract("ce575ae1bb6153330d2")
	if sc != nil {
		t.Fail()
		return
	}
	log.Printf("%v", sc)
}

func TestSmartContractStruct(t *testing.T) {

	sc := neoutils.NewSmartContract("ce575ae1bb6153330d20c560acb434dc5755241b")
	if sc == nil {
		t.Fail()
		return
	}
	log.Printf("%v", sc)
}

func UTXODataForSmartContract() smartcontract.Unspent {

	gasTX1 := smartcontract.UTXO{
		Index: 0,
		TXID:  "969ab37e0560aa2717fa8a7878c84ef37946e40bf34a06933f779a1f1b1816e2",
		Value: float64(1) / float64(100000000),
	}

	gasBalance := smartcontract.Balance{
		Amount: float64(715800000000) / float64(100000000),
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

	wif := "KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr"
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
	}

	to := smartcontract.ParseNEOAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
	if to == nil {
		//invalid neo address
		t.Fail()
		return
	}
	args := []interface{}{to, 1000}
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
