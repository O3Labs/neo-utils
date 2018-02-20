package neoutils

import (
	"log"
	"testing"
)

func TestScriptHashToNEOAddress(t *testing.T) {
	hash := "e9eed8dc39332032dc22e5d6e86332c50327ba23"
	address := ScriptHashToNEOAddress(hash)

	scripthash := NEOAddressToScriptHash(address)
	log.Printf("address = %v result = %s", address, scripthash)

	if scripthash != hash {
		t.Fail()
	}
}

func TestNEOAddressToScriptHash(t *testing.T) {
	hash := NEOAddressToScriptHash("APYB8TgR8K3rAMfYt2cCfQj3zV2Rt1oTPe")
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
