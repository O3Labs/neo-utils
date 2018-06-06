package neoutils_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/nep2"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func TestMintTokensFromMobile(t *testing.T) {
	scriptHash := "5f03828cb45198eedd659d264b6d3a1c889978ce"

	wif := ""
	wallet, _ := neoutils.GenerateFromWIF(wif)
	log.Printf("address = %v\n address hash = %x", wallet.Address, wallet.HashedSignature)

	neo := string(smartcontract.NEO)
	// gas := string(smartcontract.GAS)
	amount := float64(2)
	remark := "o3x"
	network := "private"
	networkFeeAmountInGAS := float64(0.0011)
	tx, err := neoutils.MintTokensRawTransactionMobile(network, scriptHash, wif, neo, amount, remark, networkFeeAmountInGAS)
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
