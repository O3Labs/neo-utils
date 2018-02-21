package smartcontract

import (
	"log"
	"testing"
)

func TestParseNEOAddress(t *testing.T) {
	to := ParseNEOAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
	log.Printf("%x", to)
}
func TestScriptBuilder(t *testing.T) {
	s := NewScriptBuilder()
	scriptHash, err := NewScriptHash("84a392ce6cddcc5892b9368aed43227e09b26341")
	if err != nil {
		log.Printf("err = %v", err)
		t.Fail()
		return
	}
	to := ParseNEOAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
	args := []interface{}{to, 1}
	s.pushContractInvoke(scriptHash, "mintTokensTo", args)
	log.Printf("%x", s.ToBytes())
	//from swift
	//3a51 14 23ba2703c53263e8d6e522dc32203339dcd8eee952c10c6d696e74546f6b656e73546f674163b2097e2243ed8a36b99258ccdd6cce92a384
	//3a51 14 23ba2703c53263e8d6e522dc32203339dcd8eee952c10c6d696e74546f6b656e73546f674163b2097e2243ed8a36b99258ccdd6cce92a384
}

func TestPushInt(t *testing.T) {
	s := NewScriptBuilder()
	s.pushInt(100000000)
	log.Printf("%x", s.ToBytes())
}

func TestPushDataWithInt(t *testing.T) {
	s := NewScriptBuilder()
	s.pushData(100000000)
	log.Printf("%x", s.ToBytes())
}

func TestPushArray(t *testing.T) {
	args := []interface{}{"e9eed8dc39332032dc22e5d6e86332c50327ba23", "e9eed8dc39332032dc22e5d6e86332c50327ba23", 1}
	s := NewScriptBuilder()
	s.pushData(args)
	log.Printf("%x", s.ToBytes())
}

func TestClear(t *testing.T) {
	s := NewScriptBuilder()
	s.pushData([]byte("test"))
	s.Clear()
	if len(s.ToBytes()) > 0 {
		t.Fail()
		return
	}
}
