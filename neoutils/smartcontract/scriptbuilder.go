package smartcontract

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/apisit/btckeygenie/btckey"
)

type ScriptHash []byte
type NEOAddress []byte

func (s ScriptHash) ToString() string {
	return hex.EncodeToString(s)
}
func ParseNEOAddress(address string) NEOAddress {
	v, b, _ := btckey.B58checkdecode(address)
	if v != 0x17 {
		return nil
	}
	return NEOAddress(b)
}

type ScriptBuilderInterface interface {
	pushContractInvoke(scriptHash ScriptHash, operation string, args []interface{})
	ToBytes() []byte
	FullHexString() string
	pushInt(value int) error
	pushData(data interface{}) error
	Clear()
}

func NewScriptBuilder() ScriptBuilderInterface {
	return &ScriptBuilder{RawBytes: []byte{}}
}

type ScriptBuilder struct {
	RawBytes []byte
}

func (s ScriptBuilder) ToBytes() []byte {
	return s.RawBytes
}

func (s *ScriptBuilder) Clear() {
	s.RawBytes = []byte{}
}

func (s ScriptBuilder) FullHexString() string {
	b := s.ToBytes()
	return hex.EncodeToString(b)
}

func (s *ScriptBuilder) pushOpCode(opcode OpCode) {
	log.Printf("pushing opcode %x\n", opcode)
	s.RawBytes = append(s.RawBytes, byte(opcode))
}

func (s *ScriptBuilder) pushInt(value int) error {
	switch {
	case value == -1:
		s.pushOpCode(PUSHM1)
		return nil
	case value == 0:
		s.pushOpCode(PUSH0)
		return nil
	case value >= 1 && value < 16:
		rawValue := byte(PUSH1) + byte(value) - 1
		s.RawBytes = append(s.RawBytes, rawValue)
		return nil
	case value >= 16:
		num := make([]byte, 8)
		binary.LittleEndian.PutUint64(num, uint64(value))
		s.RawBytes = append(s.RawBytes, bytes.TrimRight(num, "\x00")...)
		return nil
	}
	return nil
}
func (s *ScriptBuilder) pushHexString(hexString string) error {
	b, err := hex.DecodeString(hexString)
	if err != nil {
		return err
	}
	count := len(b)

	countBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(countBytes, uint64(count))
	trimmedCountByte := bytes.TrimRight(countBytes, "\x00")
	log.Printf("trimmed = %v", trimmedCountByte)

	if count < int(PUSHBYTES75) {
		s.RawBytes = append(s.RawBytes, trimmedCountByte...)
		s.RawBytes = append(s.RawBytes, b...)
	} else if count < 0x100 {
		s.pushOpCode(PUSHDATA1)
		s.RawBytes = append(s.RawBytes, trimmedCountByte...)
		s.RawBytes = append(s.RawBytes, b...)
	} else if count < 0x10000 {
		s.pushOpCode(PUSHDATA2)
		s.RawBytes = append(s.RawBytes, trimmedCountByte...)
		s.RawBytes = append(s.RawBytes, b...)
	} else {
		s.pushOpCode(PUSHDATA4)
		s.RawBytes = append(s.RawBytes, trimmedCountByte...)
		s.RawBytes = append(s.RawBytes, b...)
	}
	log.Println("pushing hex string", b)
	return nil
}

func (s *ScriptBuilder) pushData(data interface{}) error {
	switch e := data.(type) {
	case NEOAddress:
		//when pushing neo address as an arg. we need length so we need to push a hex string
		return s.pushHexString(fmt.Sprintf("%x", e))
	case ScriptHash:
		s.RawBytes = append(s.RawBytes, e...)
		return nil
	case string:
		return s.pushHexString(e)
	case []byte:
		return s.pushHexString(hex.EncodeToString(e))
	case bool:
		if e == true {
			s.pushOpCode(PUSH1)
		} else {
			s.pushOpCode(PUSH0)
		}
		return nil
	case []interface{}:

		count := len(e)
		//reverse the array first
		for i := len(e) - 1; i >= 0; i-- {
			s.pushData(e[i])
		}
		s.pushData(count)
		s.pushOpCode(PACK)
		return nil
	case int:
		log.Println("pushing int", e)
		s.pushInt(e)
		return nil
	}
	return nil
}

func NewScriptHash(hexString string) (ScriptHash, error) {
	b, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, err
	}
	//we need to reverse the script hash to little endian
	reversed := reverseBytes(b)
	return ScriptHash(reversed), nil
}

func (s *ScriptBuilder) pushContractInvoke(scriptHash ScriptHash, operation string, args []interface{}) {
	//args needs to be reversed.
	if args != nil {
		s.pushData(args)
	}
	s.pushData([]byte(operation)) //operation is in string we need to convert it to hex first

	s.pushOpCode(APPCALL) //use APPCALL only
	s.pushData(scriptHash)
	s.RawBytes = append([]byte{byte(len(s.RawBytes))}, s.RawBytes...)
}
