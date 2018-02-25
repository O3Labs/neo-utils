package neoutils

import (
	"encoding/hex"
	"log"
	"testing"
)

func TestScriptHashToNEOAddress(t *testing.T) {
	hash := "5e4095c8c90a7b83f509d8e08a48ba0c11fbce3a"
	address := ScriptHashToNEOAddress(hash)

	scripthash := NEOAddressToScriptHash(address)
	log.Printf("address = %v result = %s", address, scripthash)

	if scripthash != hash {
		t.Fail()
	}
}

func TestNEOAddressToScriptHash(t *testing.T) {
	hash := NEOAddressToScriptHash("AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn")
	log.Printf("%v", hash)
	b, _ := hex.DecodeString(hash)
	log.Printf("%x", ReverseBytes(b))
}

func TestValidateNEOAddress(t *testing.T) {
	valid := ValidateNEOAddress("APYB8TgR8K3rAMfYt2cCfQj3zV2Rt1oTPe")
	if valid == false {
		t.Fail()
	}
}

func TestValidateNEOAddressInvalidAddress(t *testing.T) {
	valid := ValidateNEOAddress("APYB8TgR8K3rAMfYt2cCfQj3zV2Rt1oTPe1")
	if valid == true {
		t.Fail()
	}
}

func TestConverting(t *testing.T) {
	hex := "e1f505"

	value := ConvertByteArrayToBigInt(hex)

	log.Printf("%v", value)
}

func TestParseNEP9(t *testing.T) {
	uri := "neo:AeNkbJdiMx49kBStQdDih7BzfDwyTNVRfb?assetID=602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7&amount=0.11&description=for%20a%20coffee"
	nep9, err := ParseNEP9URI(uri)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	log.Printf("%+v", nep9)
}

// func TestParseScriptFromTX(t *testing.T) {
// 	// expectedScripthash := "b7c1f850a025e34455e7e98c588c784385077fb1"
// 	// expectedOperation := []byte("mintTokensTo") // 6d696e74546f6b656e73546f
// 	// expectedToAddress := "AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn"
// 	// expectedTokenAmount := 1

// 	// log.Printf("operation %x", expectedOperation)

// 	//0x67 = APPCALL
// 	// script := "51143acefb110cba488ae0d809f5837b0ac9c895405e52c10c6d696e74546f6b656e73546f67b17f078543788c588ce9e75544e325a050f8c1b7"

// }
