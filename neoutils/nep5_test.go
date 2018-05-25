package neoutils_test

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func TestMintTokens(t *testing.T) {
	scripthash := "55d8d97603701a34f1bda8c30777c8c04deefe55"
	fee := smartcontract.NetworkFeeAmount(0.001)
	nep5 := neoutils.UseNEP5WithNetworkFee(scripthash, fee)

	wif := ""
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	unspent, err := unspent(privateNetwallet.Address)
	if err != nil {
		t.Fail()
		return
	}

	remark := "O3TX"

	asset := smartcontract.NEO
	amount := float64(10)

	tx, txID, err := nep5.MintTokensRawTransaction(*privateNetwallet, asset, amount, unspent, remark)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("txID = %v", txID)
	log.Printf("%x", tx)
}

//TEST with fee. succeeded
func TestTransferNEP5PrivateNet(t *testing.T) {

	//this is APT token
	scripthash := "0x7cd338644833db2fd8824c410e364890d179e6f8"
	fee := smartcontract.NetworkFeeAmount(0)
	nep5 := neoutils.UseNEP5WithNetworkFee(scripthash, fee)

	wif := ""
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	log.Printf("wallet address %v", privateNetwallet.Address)
	hash := neoutils.NEOAddressToScriptHashWithEndian(privateNetwallet.Address, binary.LittleEndian)
	log.Printf("wallet hash %v", hash)

	amount := float64(1000)
	to := smartcontract.ParseNEOAddress("AQaZPqcv9Kg2x1eSrF8UBYXLK4WQoTSLH5")

	unspent := smartcontract.Unspent{}

	attributes := map[smartcontract.TransactionAttribute][]byte{}
	//address hash is a hex string
	addressScriptHash := neoutils.NEOAddressToScriptHashWithEndian(privateNetwallet.Address, binary.LittleEndian)
	b, _ := hex.DecodeString(addressScriptHash)
	attributes[smartcontract.Script] = []byte(b)
	attributes[smartcontract.Remark1] = []byte(fmt.Sprintf("O3TXAPT%v", time.Now().Unix()))

	tx, txID, err := nep5.TransferNEP5RawTransaction(*privateNetwallet, to, amount, unspent, attributes)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("txID %v ", txID)
	log.Printf("%x", tx)
}
