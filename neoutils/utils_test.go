package neoutils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"testing"
)

func TestScriptHashToNEOAddress(t *testing.T) {
	hashLittleEndian := "2b41aea9d405fef2e809e3c8085221ce944527a7"
	expectedAddress := "AKibPRzkoZpHnPkF6qvuW2Q4hG9gKBwGpR"
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
	bigEndian := hex2bytes("74f2dc36a68fdc4682034178eb2220729231db76")
	address := ScriptHashToNEOAddress(fmt.Sprintf("%x", bigEndian))
	log.Printf("%v", address)
}

func TestNEOAddressToScriptHash(t *testing.T) {
	hash := NEOAddressToScriptHashWithEndian("Ab5atNiFFWzFTq55HAniJu4tMKN6hzdGEQ", binary.LittleEndian)
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

func TestConvertingByteArray(t *testing.T) {
	hexByteArray := "fec99a3b00000000"
	value := ConvertByteArrayToBigInt(hexByteArray)

	log.Printf("%v", value)
}

func TestParseNEP9(t *testing.T) {
	uri := "neo:AafQxV6wQhtGYGYFboEyBjw3eMYNtBFW8J?asset=GAS&amount=1"
	nep9, err := ParseNEP9URI(uri)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	log.Printf("%+v", nep9)
}

func TestReverse(t *testing.T) {
	b := HexTobytes("321253665742813601c8b5414c96d575990806f1")
	log.Printf("%x", ReverseBytes(b))
}

func TestConvertLittleEndianScripthashToAddress(t *testing.T) {
	b := ReverseBytes(hex2bytes("7933cdd780a209f9779a9745e81f566048d4288d"))
	address := ScriptHashToNEOAddress(fmt.Sprintf("%x", b))
	log.Printf("%v", address)
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
	publicKey := "032d47663ca1bb94f6f251df31b33615d43e1ca417c1b40322b6acd33f8fafd314"
	b, _ := hex.DecodeString(publicKey)
	address := PublicKeyToNEOAddress(b)
	log.Printf("%v", address)
}

// return list of primes less than N
func sieveOfEratosthenes(N int) (primes []int) {
	b := make([]bool, N)
	for i := 2; i < N; i++ {
		if b[i] == true {
			continue
		}
		primes = append(primes, i)
		for k := i * i; k < N; k += i {
			b[k] = true
		}
	}
	return
}

func TestPublicKeyToCustomAddress(t *testing.T) {
	publicKey := "032d47663ca1bb94f6f251df31b33615d43e1ca417c1b40322b6acd33f8fafd314"
	b, _ := hex.DecodeString(publicKey)
	primes := sieveOfEratosthenes(300)
	for _, p := range primes {
		// fmt.Println(p)

		h := fmt.Sprintf("%02x", p)

		bt, _ := hex.DecodeString(h)
		// log.Printf("%v %x", p, bt[0])
		prefix := uint8(bt[0])
		address := PublicKeyToCustomAddress(prefix, b)
		log.Printf("%v %x %v", prefix, bt[0], address)
	}

}
