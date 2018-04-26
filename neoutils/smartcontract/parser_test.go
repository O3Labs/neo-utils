package smartcontract_test

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"log"
	"math"
	"testing"

	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func TestReadBigInt(t *testing.T) {
	expectedInt := 9193970688
	b, _ := hex.DecodeString("0500dc5c2402")

	bytesReader := bytes.NewReader(b)

	reader := bufio.NewReaderSize(bytesReader, len(b))

	value, _ := smartcontract.ReadBigInt(reader)
	log.Printf("expected %v got %v (%v)", expectedInt, value.Int64(), value.BitLen())
	// log.Printf("%v value = %v %v", n, value, err)
}

func TestParserGetListOfOperations(t *testing.T) {
	expectedOperation := "mintTokensTo"
	p := smartcontract.NewParserWithScript("51143acefb110cba488ae0d809f5837b0ac9c895405e52c10c6d696e74546f6b656e73546f67b17f078543788c588ce9e75544e325a050f8c1b7")
	result, err := p.GetListOfOperations()
	if err != nil {
		log.Printf("Expected: %v but got error : %v", expectedOperation, err)
		t.Fail()
		return
	}

	log.Printf("result = %v", result)
}

func TestGetScripthashFromScript(t *testing.T) {
	expectedResult := "b7c1f850a025e34455e7e98c588c784385077fb1"
	p := smartcontract.NewParserWithScript("0830f4e2644b020000140c17e908b4014177e01d1a7fc3e6b5ed1ea83905141ffb723601fe7bf5e78b9ec6f6c79d69e317b9c753c1087472616e7366657267cf9472821400ceb06ca780c2a937fec5bbec51b9661a50b4e3696743e8")
	result, err := p.GetListOfScriptHashes()
	if err != nil {
		log.Printf("Expected: %v but got error : %v", expectedResult, err)
		t.Fail()
		return
	}

	log.Printf("result = %v", result)
}

func TestParserSingleAPPCALL(t *testing.T) {
	// expectedToAddress := "AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn"

	p := smartcontract.NewParserWithScript("51143acefb110cba488ae0d809f5837b0ac9c895405e52c10c6d696e74546f6b656e73546f67b17f078543788c588ce9e75544e325a050f8c1b7")

	//the order of your method signature has the be exact to the one
	//in your deployed smart contract
	type methodSignature struct {
		ScriptHash smartcontract.ScriptHash
		Operation  smartcontract.Operation  //operation
		To         smartcontract.NEOAddress //args[0]
		Amount     int                      //args[1]
	}
	m := methodSignature{}
	list, err := p.Parse(&m)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	for _, v := range list {
		log.Printf("%+v", v.(*methodSignature))
	}
}

func TestParseMultipleTransfers(t *testing.T) {
	script := `0500bca06501145a936d7abbaae28579dd36609f910f9b50de972f147bee835ff211327677c453d5f19b693e70a361ab53c1087472616e7366657267b6155db85e53298f01e0280cc2f21a0f40c4e808f10400e1f505147e548ecd2a87dd58731e6171752b1aa11494c62f147bee835ff211327677c453d5f19b693e70a361ab53c1087472616e7366657267b6155db85e53298f01e0280cc2f21a0f40c4e808f10500dc5c240214c10704464fade3197739536450ec9531a1f24a37147bee835ff211327677c453d5f19b693e70a361ab53c1087472616e7366657267b6155db85e53298f01e0280cc2f21a0f40c4e808f166b2263911344b5b15`
	p := smartcontract.NewParserWithScript(script)
	type methodSignature struct {
		ScriptHash smartcontract.ScriptHash
		Operation  smartcontract.Operation  //operation
		From       smartcontract.NEOAddress //args[0]
		To         smartcontract.NEOAddress //args[1]
		Amount     int                      //args[2]
	}
	m := methodSignature{}
	list, err := p.Parse(&m)
	if err != nil {
		t.Fail()
		return
	}

	for _, v := range list {
		sc, ok := v.(*methodSignature)
		if ok == false {
			continue
		}
		log.Printf("%x %v %v %v", sc.ScriptHash, sc.Operation, sc.From.ToString(), sc.To.ToString(), sc.Amount)
	}
}

func TestContainsOperation(t *testing.T) {
	p := smartcontract.NewParserWithScript("51143acefb110cba488ae0d809f5837b0ac9c895405e52c10c6d696e74546f6b656e73546f67b17f078543788c588ce9e75544e325a050f8c1b7")
	contains := p.ContainsOperation("mintTokensTo")
	log.Printf("%v", contains)
	if contains == false {
		t.Fail()
	}
}

func TestContainsTransfer(t *testing.T) {
	p := smartcontract.NewParserWithScript("05007f3e3602146b55668bb616336a5c6d2da6a035e4eb856f88c41445fc40a091bd0de5e5408e3dbf6b023919a6f7d953c1087472616e7366657267c5cc1cb5392019e2cc4e6d6b5ea54c8d4b6d11acf166605efb0156b867db")
	contains := p.ContainsOperation("transfer")
	log.Printf("%v", contains)
	if contains == false {
		t.Fail()
	}
}

func TestContainsOperationTransfer(t *testing.T) {
	script := `0500bca06501145a936d7abbaae28579dd36609f910f9b50de972f147bee835ff211327677c453d5f19b693e70a361ab53c1087472616e7366657267b6155db85e53298f01e0280cc2f21a0f40c4e808f10400e1f505147e548ecd2a87dd58731e6171752b1aa11494c62f147bee835ff211327677c453d5f19b693e70a361ab53c1087472616e7366657267b6155db85e53298f01e0280cc2f21a0f40c4e808f10500dc5c240214c10704464fade3197739536450ec9531a1f24a37147bee835ff211327677c453d5f19b693e70a361ab53c1087472616e7366657267b6155db85e53298f01e0280cc2f21a0f40c4e808f166b2263911344b5b15`
	p := smartcontract.NewParserWithScript(script)
	contains := p.ContainsOperation("transfer")
	log.Printf("%v", contains)
	if contains == false {
		t.Fail()
	}
}

func TestParserNEP5Transfer(t *testing.T) {

	p := smartcontract.NewParserWithScript("0480969800146063795d3b9b3cd55aef026eae992b91063db0db142c06a9124a089e43874d7a06a6532569df05d0ab53c1087472616e7366657267fb1c540417067c270dee32f21023aa8b9b71abcef1669ef64cb3c09aa3b3")

	//the order of your method signature has the be exact to the one
	//in your deployed smart contract
	type methodSignature struct {
		ScriptHash smartcontract.ScriptHash
		Operation  smartcontract.Operation  //operation
		From       smartcontract.NEOAddress //args[0]
		To         smartcontract.NEOAddress //args[1]
		Amount     int                      //args[2]
	}
	m := methodSignature{}
	list, err := p.Parse(&m)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	for _, v := range list {
		sc, ok := v.(*methodSignature)
		if ok == false {
			continue
		}
		log.Printf("%x %v %v %v %v", sc.ScriptHash, sc.Operation, sc.From.ToString(), sc.To.ToString(), sc.Amount)
	}
}

func TestParserNEP5TransferAnother(t *testing.T) {

	p := smartcontract.NewParserWithScript("0500b78b5b1b148e47621558c061b43d93ba407d81703edb8e36ec146063795d3b9b3cd55aef026eae992b91063db0db53c1087472616e7366657267b6155db85e53298f01e0280cc2f21a0f40c4e808f16640a6b70d4e2982f7")

	//the order of your method signature has the be exact to the one
	//in your deployed smart contract
	type methodSignature struct {
		ScriptHash smartcontract.ScriptHash
		Operation  smartcontract.Operation  //operation
		From       smartcontract.NEOAddress //args[0]
		To         smartcontract.NEOAddress //args[1]
		Amount     int                      //args[2]
	}
	m := methodSignature{}
	list, err := p.Parse(&m)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	for _, v := range list {
		sc, ok := v.(*methodSignature)
		if ok == false {
			continue
		}
		log.Printf("%x %v %v %v %.4f", sc.ScriptHash, sc.Operation, sc.From.ToString(), sc.To.ToString(), float64(sc.Amount)/math.Pow(10, 8))
	}
}

func TestGetByteWithIndex(t *testing.T) {
	b, err := hex.DecodeString("08007856a33904000014eeea87dfa23bb57b80a5001d5da2fada7effa3b614442f3f603a57542373ead52d81440f95da71245a53c1087472616e7366657267187fc13bec8ff0906c079e7f4cc8276709472913")
	if err != nil {
		t.Fail()
		return
	}

	log.Printf("%x", b[len(b)-21])
}

func TestContainsTargetScriptHash(t *testing.T) {
	//https://neotracker.io/tx/8f691ec7b9e9979964de9ce3f994588f31d6b6fea2588081d20010c14f32138d
	p := smartcontract.NewParserWithScript("0600027264cd0414aaef53a5153128fcb268b0337e8e7eae5724c78f146063795d3b9b3cd55aef026eae992b91063db0db53c1087472616e7366657267cf9472821400ceb06ca780c2a937fec5bbec51b9f166f73405540036bc1d")
	dbcScriptHash := "b951ecbbc5fe37a9c280a76cb0ce0014827294cf"
	contains := p.ContainsScriptHash(dbcScriptHash)
	log.Printf("%v", contains)
}

func TestContainsTargetScriptHashAndTransferOperation(t *testing.T) {
	//https://neotracker.io/tx/8f691ec7b9e9979964de9ce3f994588f31d6b6fea2588081d20010c14f32138d
	p := smartcontract.NewParserWithScript("0600027264cd0414aaef53a5153128fcb268b0337e8e7eae5724c78f146063795d3b9b3cd55aef026eae992b91063db0db53c1087472616e7366657267cf9472821400ceb06ca780c2a937fec5bbec51b9f166f73405540036bc1d")
	dbcScriptHash := "b951ecbbc5fe37a9c280a76cb0ce0014827294cf"
	targetOperation := "transfer"
	contains := p.ContainsScriptHashAndOperation(dbcScriptHash, targetOperation)
	log.Printf("%v", contains)
}
