package smartcontract_test

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func TestReadBigInt(t *testing.T) {
	expectedInt := 1234567890
	b, _ := hex.DecodeString("04d2029649")

	bytesReader := bytes.NewReader(b)

	reader := bufio.NewReaderSize(bytesReader, len(b))

	value, _ := smartcontract.ReadBigInt(reader)
	log.Printf("expected %v got %v (%v)", expectedInt, value.Int64())
	// log.Printf("%v value = %v %v", n, value, err)

}
func TestParserGetOperationName(t *testing.T) {
	expectedOperation := "mintTokensTo"
	p := smartcontract.NewParserWithScript("51143acefb110cba488ae0d809f5837b0ac9c895405e52c10c6d696e74546f6b656e73546f67b17f078543788c588ce9e75544e325a050f8c1b7")
	result, err := p.GetOperationName()
	if err != nil {
		log.Printf("Expected: %v but got error : %v", expectedOperation, err)
		t.Fail()
		return
	}
	if result != expectedOperation {
		log.Printf("Expected: %v but got: %v", expectedOperation, result)
		t.Fail()
		return
	}
	log.Printf("result = %v", result)
}

func TestGetScripthashFromScript(t *testing.T) {
	expectedResult := "b7c1f850a025e34455e7e98c588c784385077fb1"
	p := smartcontract.NewParserWithScript("51143acefb110cba488ae0d809f5837b0ac9c895405e52c10c6d696e74546f6b656e73546f67b17f078543788c588ce9e75544e325a050f8c1b7")
	result, err := p.GetScriptHash()
	if err != nil {
		log.Printf("Expected: %v but got error : %v", expectedResult, err)
		t.Fail()
		return
	}
	if result != expectedResult {
		log.Printf("Expected: %v but got: %v", expectedResult, result)
		t.Fail()
		return
	}
	log.Printf("result = %v", result)
}

func TestParser(t *testing.T) {
	// expectedToAddress := "AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn"

	p := smartcontract.NewParserWithScript("51143acefb110cba488ae0d809f5837b0ac9c895405e52c10c6d696e74546f6b656e73546f67b17f078543788c588ce9e75544e325a050f8c1b7")

	//the order of your method signature has the be exact to the one
	//in your deployed smart contract
	type methodSignature struct {
		Operation smartcontract.Operation  //operation
		To        smartcontract.NEOAddress //args[0]
		Amount    int                      //args[1]
	}
	m := methodSignature{}
	err := p.Parse(&m)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("%+v", m)
	log.Printf("%+v %v", m.To.ToString(), m.Amount)
}

func TestParserNEP5Transfer(t *testing.T) {

	p := smartcontract.NewParserWithScript("05006fe0d60114a20d664878bacc0114f8c594b5dc9065ce04f6eb14e484ee21fef450c92e9aed3968c6de1d58d8a9e853c1087472616e7366657267f91d6b7085db7c5aaf09f19eeec1ca3c0db2c6ec")

	//the order of your method signature has the be exact to the one
	//in your deployed smart contract
	type methodSignature struct {
		Operation smartcontract.Operation  //operation
		From      smartcontract.NEOAddress //args[0]
		To        smartcontract.NEOAddress //args[1]
		Amount    int                      //args[2]
	}
	m := methodSignature{}
	err := p.Parse(&m)
	if err != nil {
		t.Fail()
		return
	}
	// log.Printf("%+v", m)
	log.Printf("%v from %v to %v amount =%v", m.Operation, m.From.ToString(), m.To.ToString(), m.Amount)
}
