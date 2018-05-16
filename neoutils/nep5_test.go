package neoutils_test

import (
	"encoding/binary"
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

//TEST with fee. succeeded
func TestTransferNEP5(t *testing.T) {

	scripthash := "7c1de0a1fba67cbddbfea27aed370ff2ff35e8b2"
	fee := smartcontract.NetworkFeeAmount(0)
	nep5 := neoutils.UseNEP5WithNetworkFee(scripthash, fee)

	//AQmq2yU7DupE4VddmEoweKiJFyGhAAEZeH
	wif := ""
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	amount := float64(1)
	to := smartcontract.ParseNEOAddress("AQaZPqcv9Kg2x1eSrF8UBYXLK4WQoTSLH5")

	// unspent, err := utxoFromO3Platform("test", privateNetwallet.Address)
	// if err != nil {
	// 	t.Fail()
	// 	return
	// }

	unspent := smartcontract.Unspent{}
	addressScriptHash := neoutils.NEOAddressToScriptHashWithEndian(privateNetwallet.Address, binary.LittleEndian)
	// log.Printf("address hash %v", addressScriptHash)
	// remark := "O3TX111"
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	attributes[smartcontract.Script] = []byte(addressScriptHash)
	// attributes[smartcontract.Remark1] = []byte(remark)

	tx, txID, err := nep5.TransferNEP5RawTransaction(*privateNetwallet, to, amount, unspent, attributes)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("txID %v ", txID)
	log.Printf("%x", tx)
}

// func TestMintTokens(t *testing.T) {
// 	scripthash := "cc1bf80ceb9db91792c84feb8353921d9df3b4e8"
// 	fee := smartcontract.NetworkFeeAmount(0.001)
// 	nep5 := neoutils.UseNEP5WithNetworkFee(scripthash, fee)

// 	wif := "L5h6cTh45egotcxFZ2rkF1gv7rLxx9rScfuja9kEVEE9mEj9Uwtv"
// 	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
// 	if err != nil {
// 		log.Printf("%v", err)
// 		t.Fail()
// 		return
// 	}
// 	unspent, err := unspent(privateNetwallet.Address)
// 	if err != nil {
// 		t.Fail()
// 		return
// 	}

// 	remark := "O3TX"

// 	asset := smartcontract.NEO
// 	amount := float64(10)

// 	tx, txID, err := nep5.MintTokensRawTransaction(*privateNetwallet, asset, amount, unspent, remark)
// 	if err != nil {
// 		t.Fail()
// 		return
// 	}
// 	log.Printf("txID = %v", txID)
// 	log.Printf("%x", tx)
// }
