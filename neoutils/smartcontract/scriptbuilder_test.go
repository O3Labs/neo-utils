package smartcontract

import (
	"log"
	"testing"
)

func TestParseNEOAddress(t *testing.T) {
	to := ParseNEOAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
	log.Printf("%x", to)
}

func TestPushContractInvocationScript(t *testing.T) {
	s := NewScriptBuilder()
	scriptHash, err := NewScriptHash("84a392ce6cddcc5892b9368aed43227e09b26341")
	if err != nil {
		log.Printf("err = %v", err)
		t.Fail()
		return
	}
	to := ParseNEOAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y1")
	if to == nil {
		//invalid neo address
		t.Fail()
		return
	}
	args := []interface{}{to, 1}
	b := s.generateContractInvocationScript(scriptHash, "mintTokensTo", args)
	log.Printf("%x", b)
	//from swift
	//3a51 14 23ba2703c53263e8d6e522dc32203339dcd8eee952c10c6d696e74546f6b656e73546f674163b2097e2243ed8a36b99258ccdd6cce92a384
	//from go
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

func TestPushTransactionType(t *testing.T) {
	s := NewScriptBuilder()
	s.pushData(InvocationTransaction)
	log.Printf("%x", s.ToBytes())
}

func TestPushTransactionAttibute(t *testing.T) {
	s := NewScriptBuilder()
	s.pushData(Remark1)
	log.Printf("%x", s.ToBytes())
}

func TestPushLength(t *testing.T) {
	s := NewScriptBuilder()
	s.pushLength(10)
	log.Printf("%x", s.ToBytes())
}

func TestGenerateTransactionAttributes(t *testing.T) {
	s := NewScriptBuilder()
	attributes := map[TransactionAttribute][]byte{}
	attributes[Remark] = []byte("test")
	attributes[Remark2] = []byte("test2")
	attributes[Remark3] = []byte("test3")
	b, err := s.generateTransactionAttributes(attributes)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("%v", b)
}

func UTXOData() Unspent {
	gasTX2 := UTXO{
		Index: 0,
		TXID:  "ad8d65c22de1873dea36587a989a4563c7264c48ed20a6edbe957bbe428984c0",
		Value: 40.0,
	}
	gasTX1 := UTXO{
		Index: 1,
		TXID:  "1b640fc70e127a74ab6785afe155f089e08a153b2effc7a4bed8b6690cfc65fe",
		Value: 7608.0,
	}

	gasBalance := Balance{
		Amount: 7648.0,
		UTXOs:  []UTXO{gasTX1, gasTX2},
	}

	neoTX1 := UTXO{
		Index: 0,
		TXID:  "e8b8bf4f98490368fc1caa86f8646e7383bb52751ffc3a1a7e296d715c4382ed",
		Value: 100000000,
	}

	neoBalance := Balance{
		Amount: 100000000,
		UTXOs:  []UTXO{neoTX1},
	}

	unspent := Unspent{
		Assets: map[NativeAsset]*Balance{},
	}
	unspent.Assets[neo] = &neoBalance
	unspent.Assets[gas] = &gasBalance
	return unspent
}

func TestGenerateTransactionInput(t *testing.T) {
	s := NewScriptBuilder()
	assetToSend := gas
	amount := float64(5000)
	unspent := UTXOData()
	b, err := s.generateTransactionInput(unspent, assetToSend, amount)
	if err != nil {
		log.Printf("err = %v", err)
		t.Fail()
		return
	}

	log.Printf("%x %v", b, len(b))
	//swift
	//2c0848942be7b95beeda620ed484c26c763459a987a5836ea3d87e12dc2658dad00fe65fcc69b6d8bea4c7ff2e3b158ae089f055e1af8567ab747a12ec7f641b00
	//go
	//2c0848942be7b95beeda620ed484c26c763459a987a5836ea3d87e12dc2658dad00 fe65fc0c69b6d8bea4c7ff2e3b158ae089f055e1af8567ab747a120ec70f641b 00
}
