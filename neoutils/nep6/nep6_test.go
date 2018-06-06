package nep6_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/nep2"
)

func TestNEWNEP6Wallet(t *testing.T) {
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

	b, err := json.Marshal(nep6Wallet)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("%v", string(b))
}
