package neoutils

import (
	"log"
	"testing"
)

func TestScriptHashToNEOAddress(t *testing.T) {
	hash := "84a392ce6cddcc5892b9368aed43227e09b26341"
	address := ScriptHashToNEOAddress(hash)

	scripthash := NEOAddressToScriptHash(address)
	log.Printf("address = %v result = %s", address, scripthash)

	if scripthash != hash {
		t.Fail()
	}
}

func TestNEOAddressToScriptHash(t *testing.T) {
	hash := NEOAddressToScriptHash("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
	log.Printf("%v", hash)
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
	hex := "00e1f505"
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
