package neoutils_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
)

func TestNEP2ToWallet(t *testing.T) {
	encrypted := "6PYVPVe1fQznphjbUxXP9KZJqPMVnVwCx5s5pr5axRJ8uHkMtZg97eT5kL"
	passphrase := "TestingOneTwoThree"

	w, err := neoutils.NEP2DecryptToWallet(encrypted, passphrase)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("%v", w)
}
