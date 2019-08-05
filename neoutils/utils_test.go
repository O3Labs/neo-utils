package neoutils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"testing"
)

func TestScriptHashToNEOAddress(t *testing.T) {
	hashLittleEndian := "d3c3e2d374c2bb9471e7e66010caf176cb407a88"
	expectedAddress := "Ab5atNiFFWzFTq55HAniJu4tMKN6hzdGEQ"
	bigEndian := ReverseBytes(hex2bytes(hashLittleEndian))

	//ScriptHashToNEOAddress always takes big endian hash
	address := ScriptHashToNEOAddress(fmt.Sprintf("%x", bigEndian))
	scripthash := NEOAddressToScriptHashWithEndian(address, binary.LittleEndian)
	log.Printf("address = %v result = %s", address, scripthash)
	if address != expectedAddress {
		t.Fail()
	}
}

func TestSmartContractScripthashToAddress(t *testing.T) {
	address := ScriptHashToNEOAddress("9121e89e8a0849857262d67c8408601b5e8e0524")
	log.Printf("%v", address)
}

func TestNEOAddressToScriptHash(t *testing.T) {
	hash := NEOAddressToScriptHashWithEndian("AQV8FNNi2o7EtMNn4etWBYx1cqBREAifgE", binary.LittleEndian)
	b, _ := hex.DecodeString(hash)
	log.Printf("\nlittle endian %v \nbig endian %x", hash, ReverseBytes(b))
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
	hexByteArray := "80778e06" //500000000000000000
	//hex := "005c7c875e" = 405991873536
	value := ConvertByteArrayToBigInt(hexByteArray)
	vvv := float64(value.Int64()) / float64(math.Pow10(8))
	log.Printf("%v %.8f", value, vvv)
}

func TestParseNEP9(t *testing.T) {
	uri := "neo:AeNkbJdiMx49kBStQdDih7BzfDwyTNVRfb?asset=602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7&amount=0.11&description=for%20a%20coffee"
	nep9, err := ParseNEP9URI(uri)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	log.Printf("%+v", nep9)
}

func TestReverse(t *testing.T) {
	b := HexTobytes("73ef176d9f12809e64363b2b5f4553abecca7aae157327f190323cfa0e42c815")
	log.Printf("%x", ReverseBytes(b))
}

func TestHash160(t *testing.T) {
	address := "AJShjraX4iMJjwVt8WYYzZyGvDMxw6Xfbe"
	b := Hash160([]byte(address))
	log.Printf("%v", b)
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

func TestPublicKeyToNEOAddress(t *testing.T) {
	publicKey := "020ef8767aeb514780a8fb0a21f2568c521eb1e633a161dcdc39e78745762cb843"
	b, _ := hex.DecodeString(publicKey)
	address := PublicKeyToNEOAddress(b)
	log.Printf("%v", address)
}

func TestVMCodeToNEOAdress(t *testing.T) {
	code := "5121020ef8767aeb514780a8fb0a21f2568c521eb1e633a161dcdc39e78745762cb843ae"
	b, _ := hex.DecodeString(code)
	address := VMCodeToNEOAddress(b)
	log.Printf("%v", address)
}
