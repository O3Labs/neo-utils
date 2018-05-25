package smartcontract_test

import (
	"encoding/hex"
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func TestParseNEOAddress(t *testing.T) {
	to := smartcontract.ParseNEOAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
	log.Printf("%x", to)
}

func TestNewScriptHash(t *testing.T) {
	scriptHash, err := smartcontract.NewScriptHash("b7c1f850a025e34455e7e98c588c784385077fb1")
	if err != nil {
		log.Printf("err = %v", err)
		t.Fail()
		return
	}
	s := hex.EncodeToString(scriptHash.ToBigEndian())
	log.Printf("%v", s)
}

func TestGenerateInvokeScript(t *testing.T) {
	scriptHash, err := smartcontract.NewScriptHash("0x7cd338644833db2fd8824c410e364890d179e6f8")
	if err != nil {
		log.Printf("err = %v", err)
		t.Fail()
		return
	}
	args := []interface{}{}
	s := smartcontract.NewScriptBuilder()
	s.GenerateContractInvocationScript(scriptHash, "name", args)
	s.GenerateContractInvocationScript(scriptHash, "symbol", args)
	s.GenerateContractInvocationScript(scriptHash, "totalSupply", args)
	log.Printf("%x", s.ToBytes())
}

// func TestPushContractInvocationScript(t *testing.T) {
// 	s := NewScriptBuilder()
// 	scriptHash, err := NewScriptHash("b7c1f850a025e34455e7e98c588c784385077fb1")
// 	if err != nil {
// 		log.Printf("err = %v", err)
// 		t.Fail()
// 		return
// 	}
// 	to := ParseNEOAddress("AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn")
// 	if to == nil {
// 		//invalid neo address
// 		t.Fail()
// 		return
// 	}
// 	args := []interface{}{to, 715799899998}
// 	b := s.GenerateContractInvocationData(scriptHash, "mintTokensTo", args)
// 	log.Printf("%x", b)
// 	//from swift
// 	//3a511423ba273c53263e8d6e522dc32203339dcd8eee952 c1 c6d696e74546f6b656e73546f67b798b0251a6a85d2699928911afbdaefaf8470
// 	//from go
// 	// 8e830000001423ba273c53263e8d6e522dc32203339dcd8eee952c1c6d696e74546f6b656e73546f671b245557dc34b4ac60c520d335361bbe15a57ce
// 	//3be8031423ba2703c53263e8d6e522dc32203339dcd8eee952c10c6d696e74546f6b656e73546f671b245557dc34b4ac60c5200d335361bbe15a57ce
// }

// func TestPushInt(t *testing.T) {
// 	s := NewScriptBuilder()
// 	v := int(1234567890)
// 	s.pushInt(v)

// 	log.Printf("%+v %x %x", s.ToBytes(), s.ToBytes(), uintToBytes(uint(v)))

// 	//from go    715799899998 = 5eafffa8a6
// 	//from swift 715799899998 = 5eafffa8a6000
// }

// func TestPushDataWithInt(t *testing.T) {
// 	s := NewScriptBuilder()
// 	s.pushData(100000000)
// 	log.Printf("%x", s.ToBytes())
// }

// func TestPushArray(t *testing.T) {
// 	to := ParseNEOAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
// 	args := []interface{}{to, 1000}
// 	s := NewScriptBuilder()
// 	s.pushData(args)
// 	log.Printf("%x", s.ToBytes())
// }

// func TestToScriptHash(t *testing.T) {
// 	to := ParseNEOAddress("AJShjraX4iMJjwVt8WYYzZyGvDMxw6Xfbe")
// 	s := NewScriptBuilder()
// 	s.pushData(to)
// 	log.Printf("%x", s.ToScriptHash())
// }

// func TestClear(t *testing.T) {
// 	s := NewScriptBuilder()
// 	s.pushData([]byte("test"))
// 	s.Clear()
// 	if len(s.ToBytes()) > 0 {
// 		t.Fail()
// 		return
// 	}
// }

// func TestPushTransactionType(t *testing.T) {
// 	s := NewScriptBuilder()
// 	s.pushData(InvocationTransaction)
// 	log.Printf("%x", s.ToBytes())
// }

// func TestPushTransactionAttibute(t *testing.T) {
// 	s := NewScriptBuilder()
// 	s.pushData(Remark1)
// 	log.Printf("%x", s.ToBytes())
// }

// func TestPushLength(t *testing.T) {
// 	s := NewScriptBuilder()
// 	s.pushLength(33)
// 	log.Printf("%x", s.ToBytes())
// }

// func TestGenerateTransactionAttributes(t *testing.T) {
// 	s := NewScriptBuilder()
// 	attributes := map[TransactionAttribute][]byte{}
// 	attributes[Remark] = []byte("test")
// 	attributes[Remark2] = []byte("test2")
// 	attributes[Remark3] = []byte("test3")
// 	b, err := s.GenerateTransactionAttributes(attributes)
// 	if err != nil {
// 		t.Fail()
// 		return
// 	}
// 	log.Printf("%v", b)
// }

// func TestGenerateTransactionInput(t *testing.T) {
// 	s := NewScriptBuilder()
// 	assetToSend := GAS
// 	amount := float64(5000)
// 	unspent := UTXODataForSmartContract()
// 	b, err := s.GenerateTransactionInput(unspent, assetToSend, amount)
// 	if err != nil {
// 		log.Printf("err = %v", err)
// 		t.Fail()
// 		return
// 	}

// 	log.Printf("%x %v", b, len(b))
// 	//swift
// 	//2c0848942be7b95beeda620ed484c26c763459a987a5836ea3d87e12dc2658dad00fe65fcc69b6d8bea4c7ff2e3b158ae089f055e1af8567ab747a12ec7f641b00
// 	//go
// 	//2c0848942be7b95beeda620ed484c26c763459a987a5836ea3d87e12dc2658dad00 fe65fc0c69b6d8bea4c7ff2e3b158ae089f055e1af8567ab747a120ec70f641b 00
// }

// func TestGenerateTransactionOutput(t *testing.T) {
// 	s := NewScriptBuilder()
// 	assetToSend := GAS
// 	amountToSend := float64(0.00000001)
// 	unspent := UTXODataForSmartContract()
// 	sender := ParseNEOAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
// 	receiver := ParseNEOAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")

// 	b, err := s.GenerateTransactionOutput(sender, receiver, unspent, assetToSend, amountToSend)
// 	if err != nil {
// 		log.Printf("%v", err)
// 		t.Fail()
// 	}
// 	log.Printf("%x", b)
// 	//52e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c600088526a741423ba2703c53263e8d6e522dc32203339dcd8eee9e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c6000584fa73d1423ba2703c53263e8d6e522dc32203339dcd8eee9
// }
