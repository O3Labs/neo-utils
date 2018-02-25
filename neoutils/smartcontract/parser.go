package smartcontract

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"reflect"

	"github.com/o3labs/neo-utils/neoutils/btckey"
)

type ParserInterface interface {
	Parse(methodSignature interface{}) error
}

type Parser struct {
	Script string
}

func NewParserWithScript(script string) Parser {
	return Parser{Script: script}
}

var _ ParserInterface = (*Parser)(nil)

func ReadHexString(reader *bufio.Reader) string {
	first := make([]byte, 1)
	reader.Read(first)
	length := int(first[0])
	if length == 0 {
		return ""
	}
	buf := make([]byte, length)
	reader.Read(buf)
	return string(buf)
}

func ReadNEOAddress(reader *bufio.Reader) (*NEOAddress, error) {
	first := make([]byte, 1)
	reader.Read(first)
	length := int(first[0])
	if length == 0 {
		return nil, fmt.Errorf("Length is zero")
	}

	buf := make([]byte, length)
	reader.Read(buf)

	address := btckey.B58checkencodeNEO(0x17, buf)
	log.Printf("reading neo address %v", address)
	neoAddress := ParseNEOAddress(address)
	return &neoAddress, nil
}
func ReadBigInt(reader *bufio.Reader) (value big.Int, e error) {

	//first byte is the number of bytes
	//however, we have to check against the PUSH[X] opcode too
	first := make([]byte, 1)
	reader.Read(first)

	//very simple switch statement for each case for code readability
	switch first[0] {
	case byte(PUSHM1):
		v := new(big.Int).SetInt64(-1)
		value = *v
	case byte(PUSH1):
		v := new(big.Int).SetInt64(1)
		value = *v
	case byte(PUSH2):
		v := new(big.Int).SetInt64(2)
		value = *v
	case byte(PUSH3):
		v := new(big.Int).SetInt64(3)
		value = *v
	case byte(PUSH4):
		v := new(big.Int).SetInt64(4)
		value = *v
	case byte(PUSH5):
		v := new(big.Int).SetInt64(5)
		value = *v
	case byte(PUSH6):
		v := new(big.Int).SetInt64(6)
		value = *v
	case byte(PUSH7):
		v := new(big.Int).SetInt64(7)
		value = *v
	case byte(PUSH8):
		v := new(big.Int).SetInt64(8)
		value = *v
	case byte(PUSH9):
		v := new(big.Int).SetInt64(9)
		value = *v
	case byte(PUSH10):
		v := new(big.Int).SetInt64(10)
		value = *v
	case byte(PUSH11):
		v := new(big.Int).SetInt64(11)
		value = *v
	case byte(PUSH12):
		v := new(big.Int).SetInt64(12)
		value = *v
	case byte(PUSH13):
		v := new(big.Int).SetInt64(13)
		value = *v
	case byte(PUSH14):
		v := new(big.Int).SetInt64(14)
		value = *v
	case byte(PUSH15):
		v := new(big.Int).SetInt64(15)
		value = *v
	case byte(PUSH16):
		v := new(big.Int).SetInt64(16)
		value = *v
	default: //with a length
		length := first[0]
		buf := make([]byte, length)
		reader.Read(buf)
		reversed := reverseBytes(buf)
		reversedHex := hex.EncodeToString(reversed)
		v, _ := new(big.Int).SetString(reversedHex, 16)
		value = *v
	}

	e = nil
	return
}

func bufioReaderFromBytes(b []byte) *bufio.Reader {
	return bufio.NewReaderSize(bytes.NewReader(b), len(b))
}

func (p *Parser) splitScriptWithAPPCALL() (operationAndArgs []byte, scriptHash []byte, e error) {
	b, err := hex.DecodeString(p.Script)
	if err != nil {
		return nil, nil, err
	}
	//split the script by the APPCALL opcode(0x67) then we will get
	//script hash on the right and args + operation on the left
	splitted := bytes.Split(b, []byte{byte(APPCALL)})
	if len(splitted) < 2 {
		return nil, nil, fmt.Errorf("invalid script: Script doesn't have APPCALL opcode")
	}
	operationAndArgs = splitted[0]
	scriptHash = splitted[1]
	return
}

func (p *Parser) GetScriptHash() (string, error) {
	_, scripthashBytes, err := p.splitScriptWithAPPCALL()
	if err != nil {
		return "", err
	}
	//this script hash from rawtransaction's script is reversed
	//in order to return the proper one that you get when call "getcontractstate"
	//we have to reverse the bytes
	return fmt.Sprintf("%x", reverseBytes(scripthashBytes)), nil
}

func (p *Parser) GetOperationName() (string, error) {
	operationAndArgs, _, err := p.splitScriptWithAPPCALL()
	if err != nil {
		return "", err
	}
	splittedPack := bytes.Split(operationAndArgs, []byte{byte(PACK)})
	if len(splittedPack) < 2 {
		return "", fmt.Errorf("invalid script: Script doesn't conform main(operation, args)")
	}
	operationString := ReadHexString(bufioReaderFromBytes(splittedPack[1]))
	return operationString, nil
}

func (p *Parser) Parse(methodSignature interface{}) error {

	operationAndArgs, bigEndianScripthash, err := p.splitScriptWithAPPCALL()
	if err != nil {
		return err
	}

	//We can split the operation and args by PACK opcode(0xc1)
	//it's packing the array of arg and the next is the number of args
	//e.g. 02c1 = pack 2 arguments
	//this can be done by reading the last byte if it's a PACK opcode
	//second last byte is the number of array that are packed in
	splittedPack := bytes.Split(operationAndArgs, []byte{byte(PACK)})
	if len(splittedPack) < 2 {
		return fmt.Errorf("invalid script: Script doesn't conform main(operation, args)")
	}

	argsWithNumbers := splittedPack[0]
	operation := splittedPack[1]
	//read from hex
	operationString := ReadHexString(bufioReaderFromBytes(operation))

	//after split, last byte is the number of args in an array
	numberOfArgsBytes := argsWithNumbers[len(argsWithNumbers)-1:]
	//This seems overkill because of the number of args should never be this large.
	//maybe we can rewrite readInt again
	numberOfArgs, err := ReadBigInt(bufioReaderFromBytes(numberOfArgsBytes))
	if err != nil {
		return err
	}
	log.Printf("numberOfArgs = %v", numberOfArgs)

	s := reflect.ValueOf(methodSignature).Elem()
	typeOfT := s.Type()

	//new bytes reader
	bytesReader := bytes.NewReader(argsWithNumbers)
	reader := bufio.NewReaderSize(bytesReader, len(argsWithNumbers))

	//because the bytes of the script is reversed
	//we will have to reverse the order of the fields in the struct too
	for i := s.NumField() - 1; i >= 0; i-- {
		field := s.Field(i)
		t := typeOfT.Field(i)
		switch t.Type {
		case reflect.TypeOf(Operation("")):
			field.SetString(operationString)
		case reflect.TypeOf(NEOAddress{}):
			neoAddress, err := ReadNEOAddress(reader)
			if err != nil {
				return err
			}
			field.SetBytes(*neoAddress)
		case reflect.TypeOf(int(0)):
			v, err := ReadBigInt(reader)
			if err != nil {
				return err
			}
			field.SetInt(v.Int64())
		}
	}

	log.Printf("%x %x", operationAndArgs, reverseBytes(bigEndianScripthash))
	return nil
}
