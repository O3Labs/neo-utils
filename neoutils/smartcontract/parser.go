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
	Parse(methodSignature interface{}) ([]interface{}, error)
	GetListOfOperations() ([]string, error)
	GetListOfScriptHashes() ([]string, error)
	ContainsOperation(operation string) bool
	ContainsScriptHash(scripthash string) bool
	ContainsScriptHashAndOperation(scripthash string, operation string) bool
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
		// log.Printf("big int = %v", v.Int64())
		value = *v
	}

	e = nil
	return
}

func bufioReaderFromBytes(b []byte) *bufio.Reader {
	return bufio.NewReaderSize(bytes.NewReader(b), len(b))
}

type appcall struct {
	operationAndArgs []byte
	scriptHash       []byte
}

func (p *Parser) splitScriptWithAPPCALL() ([]appcall, error) {

	list := []appcall{}
	b, err := hex.DecodeString(p.Script)
	if err != nil {
		return list, err
	}

	//perhaps a better way to do this is to check the 21st byte from the right side
	//if it's the APPCALL
	twentyFirstByte := b[len(b)-21]
	if twentyFirstByte == byte(APPCALL) {
		//if the 21st byte is the APPCALL. it's safe to grab the last 20 bytes. it's a scripthash
		rightSide := b[len(b)-20:] //last 20 bytes is a script hash
		leftSide := b[:len(b)-21]  //from start until the len(b) - 21

		list = append(list, appcall{
			operationAndArgs: leftSide,
			scriptHash:       reverseBytes(rightSide),
		})
		return list, nil
	}

	//if the 21st byte is not the APPCALL. we can try it again to see if the 10th byte is THROWIFNOT

	if twentyFirstByte != byte(APPCALL) {
		tenthByte := b[len(b)-10]
		//this could happen with single APPCALL or multiple APPCALL
		//so we should try to split the left side with the APPCALL and see the number of APPCALL as well
		testSplit := bytes.Split(b[:len(b)-(21+10)], []byte{byte(APPCALL)})
		if tenthByte == byte(THROWIFNOT) && len(testSplit) >= 1 {
			rightSide := b[len(b)-(20+10):] //last 20 bytes is a script hash
			leftSide := b[:len(b)-(21+10)]  //from start until the len(b) - (21 + 10)
			list = append(list, appcall{
				operationAndArgs: leftSide,
				scriptHash:       reverseBytes(rightSide),
			})
			return list, nil
		}
	}

	//split the script by the APPCALL opcode(0x67) then we will get
	//script hash on the right and args + operation on the left
	splitted := bytes.Split(b, []byte{byte(APPCALL)})
	numberOfAPPCALLs := len(splitted) - 1
	if numberOfAPPCALLs == 0 {
		return list, fmt.Errorf("invalid script: Script doesn't have APPCALL opcode")
	}

	//if we found more than 1 APPCALL. we need to split them properly

	//the script hash from rawtransaction's script is reversed
	//in order to return the proper one that you get when call "getcontractstate"
	//we have to reverse the bytes
	if numberOfAPPCALLs > 1 {
		for index := 0; index < len(splitted)-1; index++ {
			tempOperationAndArgs := []byte{}
			tempScriptHash := []byte{}
			if index == 0 {
				tempOperationAndArgs = splitted[index]
				tempScriptHash = reverseBytes(splitted[index+1][:20])
			} else {
				//check the length here first
				if len(splitted[index]) < 21 {
					log.Printf("%x", splitted[index])
					continue
				}
				//Multiple APPCALL contains THROWIFNOT to make sure that every APPCALL runs otherwise reject all
				//THROWIFNOT is 0xf1.
				//scripthash(20 bytes) + 0xf1 + [actual operation and args]
				//check the length first
				if len(splitted[index+1]) < 20 {
					continue
				}
				tempOperationAndArgs = splitted[index][21:]
				tempScriptHash = reverseBytes(splitted[index+1][:20])
			}
			a := appcall{
				operationAndArgs: tempOperationAndArgs, //left side is the operation and args
				scriptHash:       tempScriptHash,       //first 20 bytes of right side is the script hash
			}
			list = append(list, a)
		}
		return list, nil
	}

	return list, nil
}

func (p *Parser) GetListOfScriptHashes() ([]string, error) {
	list, err := p.splitScriptWithAPPCALL()
	if err != nil {
		return []string{}, err
	}
	listOfScriptHash := []string{}
	temp := map[string]bool{}
	for _, v := range list {
		scriptHash := fmt.Sprintf("%x", v.scriptHash[len(v.scriptHash)-20:])
		exist := temp[scriptHash]
		if exist == true {
			continue
		}
		listOfScriptHash = append(listOfScriptHash, scriptHash)
		temp[scriptHash] = true
	}

	return listOfScriptHash, nil
}

func (p *Parser) GetListOfOperations() ([]string, error) {
	list, err := p.splitScriptWithAPPCALL()
	if err != nil {
		return []string{}, err
	}
	listOfOperations := []string{}
	for _, v := range list {
		splittedPack := bytes.Split(v.operationAndArgs, []byte{byte(PACK)})
		if len(splittedPack) < 2 {
			return []string{}, fmt.Errorf("invalid script: Script doesn't conform main(operation, args)")
		}
		operationString := ReadHexString(bufioReaderFromBytes(splittedPack[1]))
		listOfOperations = append(listOfOperations, operationString)
	}
	return listOfOperations, nil
}

// This method return an array of given method signature
// because sometimes script can contains multiple app calls
func (p *Parser) Parse(methodSignature interface{}) ([]interface{}, error) {

	list, err := p.splitScriptWithAPPCALL()

	if err != nil {
		return nil, err
	}

	results := []interface{}{}
	for _, v := range list {
		operationAndArgs := v.operationAndArgs
		// scripthash := v.scriptHash

		//We can split the operation and args by PACK opcode(0xc1)
		//it's packing the array of arg and the next is the number of args
		//e.g. 02c1 = pack 2 arguments, 53c1 (PUSH3 + PACK)
		//this can be done by reading the last byte if it's a PACK opcode
		//second last byte is the number of array that are packed in
		splittedPack := bytes.Split(operationAndArgs, []byte{byte(PACK)})
		if len(splittedPack) < 2 {
			return nil, fmt.Errorf("invalid script: Script doesn't conform main(operation, args)")
		}

		//the operation should be on the farthest right
		//so when we split with PACK. the last index should be operation
		//this is here to make sure in case if there is 0x21 somewhere before the actual PACK
		operation := splittedPack[len(splittedPack)-1]

		//then the rest on the left should be args with the number of args
		//so we get operationsAndArgs until operations
		//operationAndArgs = [args] + [number of args]+ 0xc1 +[N bytes of operation] + [operation]
		//args(N bytes) + number of args (1 byte)
		//so we get from the first byte until (0xc1 +[N bytes of operation] + [operation]) bytes
		argsWithNumbers := operationAndArgs[:len(operationAndArgs)-len(operation)-1]

		//read from hex. first byte is the number of bytes to be read
		operationString := ReadHexString(bufioReaderFromBytes(operation))

		//after split, last byte is the number of args in an array
		numberOfArgsBytes := argsWithNumbers[len(argsWithNumbers)-1:]

		//This seems overkill because of the number of args should never be this large.
		//maybe we can rewrite readInt again
		numberOfArgs, err := ReadBigInt(bufioReaderFromBytes(numberOfArgsBytes))
		if err != nil {
			return nil, err
		}
		s := reflect.ValueOf(methodSignature).Elem()
		typeOfT := s.Type()

		//we check given method signature and parsed number of args
		if int(numberOfArgs.Int64()) != s.NumField()-2 {
			return nil, fmt.Errorf("The number of args is %v but given method signature has %v", numberOfArgs, s.NumField())
		}

		//new bytes reader
		bytesReader := bytes.NewReader(argsWithNumbers)
		reader := bufio.NewReaderSize(bytesReader, len(argsWithNumbers))

		//because the bytes of the script is reversed
		//we will have to reverse the order of the fields in the struct too
		for i := s.NumField() - 1; i >= 0; i-- {
			field := s.Field(i)
			t := typeOfT.Field(i)
			switch t.Type {
			case reflect.TypeOf(ScriptHash{}):
				field.SetBytes(v.scriptHash[:20])
			case reflect.TypeOf(Operation("")):
				field.SetString(operationString)
			case reflect.TypeOf(NEOAddress{}):
				neoAddress, err := ReadNEOAddress(reader)
				if err != nil {
					return nil, err
				}
				field.SetBytes(*neoAddress)
			case reflect.TypeOf(int(0)):
				v, err := ReadBigInt(reader)
				if err != nil {
					return nil, err
				}
				field.SetInt(v.Int64())
			}
		}

		results = append(results, methodSignature)
	}

	return results, nil
}

func (p *Parser) ContainsOperation(operation string) bool {
	operationBytes := []byte(operation)
	scriptBytes, err := hex.DecodeString(p.Script)
	if err != nil {
		return false
	}
	return bytes.Contains(scriptBytes, operationBytes)
}

func (p *Parser) ContainsScriptHash(scripthash string) bool {
	scripthashBytes, err := hex.DecodeString(scripthash)
	if err != nil {
		return false
	}
	target := reverseBytes(scripthashBytes)
	scriptBytes, err := hex.DecodeString(p.Script)
	if err != nil {
		return false
	}
	log.Printf("%x", target)
	return bytes.Contains(scriptBytes, target)
}
func (p *Parser) ContainsScriptHashAndOperation(scripthash string, operation string) bool {
	scripthashBytes, err := hex.DecodeString(scripthash)
	if err != nil {
		return false
	}
	target := reverseBytes(scripthashBytes)
	operationBytes := []byte(operation)
	scriptBytes, err := hex.DecodeString(p.Script)
	if err != nil {
		return false
	}
	return bytes.Contains(scriptBytes, target) && bytes.Contains(scriptBytes, operationBytes)
}
