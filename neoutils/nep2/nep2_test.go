package nep2_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils/nep2"
)

func TestNEP2DecryptToWallet(t *testing.T) {
	encrypted := "6PYVPVe1fQznphjbUxXP9KZJqPMVnVwCx5s5pr5axRJ8uHkMtZg97eT5kL"
	passphrase := "TestingOneTwoThree"
	p, err := nep2.NEP2DecryptToPrivateKey(encrypted, passphrase)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("%+v", p)
}

func TestNEP2(t *testing.T) {
	passphrase := "TestingOneTwoThree"
	WIF := "L44B5gGEpqEDRS9vVPz7QT35jcBG2r3CZwSwQ4fCewXAhAhqGVpP" //AStZHy8E6StCqYQbzMqi4poH7YNDHQKxvt
	expectedEncrypted := "6PYVPVe1fQznphjbUxXP9KZJqPMVnVwCx5s5pr5axRJ8uHkMtZg97eT5kL"

	// expctedUnencryptedHex := "CBF4B9F70470856BB4F40F80B87EDB90865997FFEE6DF315AB166D713AF433A5"

	encrypted, address, err := nep2.NEP2Encrypt(WIF, passphrase)
	if err != nil {
		log.Printf("err %v", err)
		return
	}
	log.Printf("encrypted = %v address = %v", encrypted, address)

	if encrypted != expectedEncrypted {
		log.Printf("expected %v : got %v", expectedEncrypted, encrypted)
		t.Fail()
		return
	}

	decrypted, err := nep2.NEP2Decrypt(encrypted, passphrase)
	if err != nil {
		log.Printf("err %v", err)
		return
	}

	if decrypted != WIF {
		t.Fail()
		return
	}
	log.Printf("decrypted = %v", decrypted)
}
