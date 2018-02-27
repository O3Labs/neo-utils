package neoutils_test

import (
	"encoding/hex"
	"fmt"
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/sss"
)

func TestNewWallet(t *testing.T) {
	w, err := neoutils.NewWallet()
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("address %v", w.Address)
	log.Printf("WIF %v", w.WIF)
}

func TestGenKey(t *testing.T) {
	privateKey := "0C28FCA386C7A227600B2FE50B7CAE11EC86D3BF1FBE471BE89827E19D72AA1D"
	wallet, _ := neoutils.GenerateFromPrivateKey(privateKey)

	log.Printf("%+v", wallet)
}

func TestGenFromWIF(t *testing.T) {
	wif := "KzULqzStT2tseGnqogXnTLG5NCT1YXa3F9Wp1Kdv9xMxFhvV6H2A"
	wallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%+v", err)
		t.Fail()
	}

	log.Printf("private key %+v", hex.EncodeToString(wallet.PrivateKey))
	log.Printf("public key %+v (%d)", hex.EncodeToString(wallet.PublicKey), len(wallet.PublicKey))
	log.Printf("wallet%+v", wallet)
	log.Printf("wallet address %+v %d", wallet.Address, len(wallet.Address))
}

func TestSSS(t *testing.T) {
	secret := "well hello there!" // our secret
	n := byte(2)                  // create 30 shares
	k := byte(2)                  // require 2 of them to combine

	shares, err := sss.Split(n, k, []byte(secret)) // split into 30 shares
	if err != nil {
		fmt.Println(err)
		return
	}

	// select a random subset of the total shares
	subset := make(map[byte][]byte, k)
	for x, y := range shares { // just iterate since maps are randomized
		subset[x] = y
		if len(subset) == int(k) {
			break
		}
	}

	// combine two shares and recover the secret
	recovered := string(sss.Combine(subset))
	if secret != recovered {
		t.Fail()
		return
	}
	fmt.Println(recovered)
}

func TestGenerateSSS(t *testing.T) {
	sharedSecret, err := neoutils.GenerateShamirSharedSecret("0C28FCA386C7A227600B2FE50B7CAE11EC86D3BF1FBE471BE89827E19D72AA1D")
	if err != nil {
		t.Fail()
		return
	}
	recovered, err := neoutils.RecoverFromSharedSecret(sharedSecret.First, sharedSecret.Second)

	fmt.Printf("%v\n%v\n%v", neoutils.BytesToHex(sharedSecret.First), neoutils.BytesToHex(sharedSecret.Second), recovered)
}

func TestRecoverFromString(t *testing.T) {

	first := neoutils.HexTobytes("b636f0821f65399cd0a2334ab6b229bfb8dce9fe569c795d14ec6177c7eba62be80668fdd6fe743286e6aec5b856ca9a15d7c9b0c82d06fa80f51a920d32ff90")
	second := neoutils.HexTobytes("27a9ad57f40fb176f305a3cdb429043c31f39921fae93de578059b2b5602040504c998bb7bb22eae441d815e37f5dce9e5fdc233dd03c3bc503d6d69d9a7b6f7")

	recovered, err := neoutils.RecoverFromSharedSecret(first, second)
	if err != nil {
		t.Fail()
		return
	}
	fmt.Printf("%v", recovered)
}
