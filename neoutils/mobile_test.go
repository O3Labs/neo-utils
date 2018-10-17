package neoutils_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/nep2"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func TestMintTokensFromMobile(t *testing.T) {
	scriptHash := "9121e89e8a0849857262d67c8408601b5e8e0524"

	// encryptedKey := ""
	// passphrase := ""
	// wif, _ := neoutils.NEP2Decrypt(encryptedKey, passphrase)
	wif := ""
	wallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}

	log.Printf("address = %v\n address hash = %x", wallet.Address, wallet.HashedSignature)
	// neo := string(smartcontract.NEO)
	gas := string(smartcontract.GAS)
	amount := float64(1)
	remark := "FIRST! APISIT FROM O3 :D"
	network := "main"
	networkFeeAmountInGAS := float64(0)
	tx, err := neoutils.MintTokensRawTransactionMobile(network, scriptHash, wif, gas, amount, remark, networkFeeAmountInGAS)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	log.Printf("txID =%v", tx.TXID)
	log.Printf("tx = %x", tx.Data)
}

func TestNEP6MobileMethod(t *testing.T) {
	passphase := "TestingOneTwoThree"
	WIF := "L44B5gGEpqEDRS9vVPz7QT35jcBG2r3CZwSwQ4fCewXAhAhqGVpP" //AStZHy8E6StCqYQbzMqi4poH7YNDHQKxvt
	encryptedKey, address, err := nep2.NEP2Encrypt(WIF, passphase)
	if err != nil {
		log.Printf("err %v", err)
		return
	}
	log.Printf("encrypted = %v", encryptedKey)
	walletName := "o3wallet"
	addressLabel := "spending"

	nep6Wallet := neoutils.GenerateNEP6FromEncryptedKey(walletName, addressLabel, address, encryptedKey)
	log.Printf("%+v", nep6Wallet)

}

func TestSerializeTX(t *testing.T) {
	data := `
{
	"sha256": "ab1ad3efa1bf2fca51219b73c676dadaf9f446b81acd72f2557fecb7a8e7d243",
	"type": 209,
	"attributes": [{
		"usage": 32,
		"data": "4d17abe11020df91ce627af28c03c9c0cfb2a6c4"
	}],
	"scripts": [],
	"gas": 0,
	"version": 1,
	"hash": "a67d3f9314383c4f7234cc9c8b7cf50602f91bf908bf6496a0c81bdc37fac7da",
	"inputs": [{
		"prevIndex": 0,
		"prevHash": "d21043bb53d70a4762ebad4fcd55fb00528f4898d97cbe4182aef5b91139ec60"
	}, {
		"prevIndex": 6,
		"prevHash": "6005967b1f6697d03cf241995fd4b2e71e56945ce0e4f815033700b993150c15"
	}],
	"outputs": [{
		"assetId": "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7",
		"scriptHash": "e707714512577b42f9a011f8b870625429f93573",
		"value": 1e-08
	}],
	"script": "0800e1f505000000001432e125258b7db0a0dffde5bd03b2b859253538ab144d17abe11020df91ce627af28c03c9c0cfb2a6c453c1076465706f73697467823b63e7c70a795a7615a38d1ba67d9e54c195a1"
}
`
	final := neoutils.SerializeTX(data)
	log.Printf("%x", final)
}

func TestNEOAddresstoScriptHashBigEndian(t *testing.T) {
	log.Printf("%v", neoutils.NEOAddresstoScriptHashBigEndian("AQV8FNNi2o7EtMNn4etWBYx1cqBREAifgE"))
}

func TestGetVarUInt(t *testing.T) {
	log.Printf("%x", neoutils.GetVarUInt(286))
}
