package neoutils_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/nep2"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func TestMintTokensFromMobile(t *testing.T) {
	scriptHash := "0x3e390ae61acb6713389c8fbbd47d1d69c32655a3"

	wif := ""
	wallet, _ := neoutils.GenerateFromWIF(wif)
	log.Printf("address = %v\n address hash = %x", wallet.Address, wallet.HashedSignature)

	neo := string(smartcontract.NEO)
	// gas := string(smartcontract.GAS)
	amount := float64(2)
	remark := "O3XMOONLIGHT2"
	network := "test"
	networkFeeAmountInGAS := float64(0)
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
