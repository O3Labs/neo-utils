package neoutils

import (
	"encoding/hex"
	"fmt"
	"log"
	"testing"
)

func TestConvertScripthashFromParamToNEOAddress(t *testing.T) {
	hash := "7bee835ff211327677c453d5f19b693e70a361ab"
	b := hex2bytes(hash)
	b = ReverseBytes(b)

	address := ScriptHashToNEOAddress(bytesToHex(b))

	scripthash := NEOAddressToScriptHash(address)
	log.Printf("address = %v result = %s", address, scripthash)

	if scripthash != hash {
		t.Fail()
	}
}

func TestScriptHashToNEOAddress(t *testing.T) {
	hash := "cc1bf80ceb9db91792c84feb8353921d9df3b4e8"

	address := ScriptHashToNEOAddress(hash)

	scripthash := NEOAddressToScriptHash(address)
	log.Printf("address = %v result = %s", address, scripthash)

	if scripthash != hash {
		t.Fail()
	}
}

func TestNEOAddressToScriptHash(t *testing.T) {
	hash := NEOAddressToScriptHash("ASH41gtWftHvhuYhZz1jj7ee7z9vp9D9wk")
	b, _ := hex.DecodeString(hash)
	log.Printf("%x %x", ReverseBytes(b), b)
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
	hex := "001175f11e"
	//hex := "005c7c875e" = 405991873536

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

func TestReverse(t *testing.T) {
	b := HexTobytes("f782294e0db7a64066f108e8c4400f1af2c20c28")
	log.Printf("%x", ReverseBytes(b))
}

func TestHash160(t *testing.T) {
	address := "AJShjraX4iMJjwVt8WYYzZyGvDMxw6Xfbe"
	b := Hash160([]byte(address))
	log.Printf("%x", b)
}

func TestHash256(t *testing.T) {
	raw := "d1002200c10a6d696e74546f6b656e736793ad7e2a1ade96c4f2358e670ef683378d14ebb201f1036f337802967e38191d9c0f2039e4890294689b7bf4a7153937fada20aa2425fc196ada7f0100967e38191d9c0f2039e4890294689b7bf4a7153937fada20aa2425fc196ada7f0200039b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc500e1f5050000000093ad7e2a1ade96c4f2358e670ef683378d14ebb29b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc5004ad4642a84230023ba2703c53263e8d6e522dc32203339dcd8eee9e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c605d3af17c0f00000023ba2703c53263e8d6e522dc32203339dcd8eee9"
	expectedResult := "389470367287e9f99e561a66d6ab5875f8375506ec1a16d54e9c628f34b8efe8"
	b, _ := hex.DecodeString(raw)

	txid := ReverseBytes(Hash256(b))
	result := fmt.Sprintf("%x", txid)
	if result != expectedResult {
		t.Fail()
	}
}
