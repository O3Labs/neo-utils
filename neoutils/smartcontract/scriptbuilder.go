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
	generateContractInvocationData(scriptHash ScriptHash, operation string, args []interface{}) []byte
	generateTransactionAttributes(attributes map[TransactionAttribute][]byte) ([]byte, error)
	generateTransactionInput(unspent Unspent, assetToSend NativeAsset, amountToSend float64) ([]byte, error)

	generateTransactionOutput(sender NEOAddress, receiver NEOAddress, unspent Unspent, assetToSend NativeAsset, amountToSend float64) ([]byte, error)
	generateInvocationScriptWithSignatures(signatures []TransactionSignature) []byte

	emptyTransactionAttributes() []byte

	ToBytes() []byte
	FullHexString() string
	pushInt(value int) error
	pushData(data interface{}) error
	Clear()
	pushLength(count int)
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
		log.Printf("raw pushInt %x %v", rawValue, rawValue)
		s.RawBytes = append(s.RawBytes, rawValue)
		return nil
	case value >= 16:
		num := make([]byte, 8)

		binary.LittleEndian.PutUint64(num, uint64(value))
		// s.pushData(bytes.Trim(num, "\x00"))
		s.pushData(num)
		return nil
	}
	return nil
}

func (s *ScriptBuilder) pushLength(count int) {
	if count == 0 {
		s.RawBytes = append(s.RawBytes, 0x00)
		return
	}
	countBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(countBytes, uint64(count))
	trimmedCountByte := bytes.Trim(countBytes, "\x00")
	s.RawBytes = append(s.RawBytes, trimmedCountByte...)
}

func uintToBytes(value uint) []byte {
	countBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(countBytes, uint64(value))
	return bytes.TrimRight(countBytes, "\x00")
}

func uint16ToFixBytes(value uint16) []byte {
	countBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(countBytes, value)
	return countBytes //bytes.TrimRight(countBytes, "\x00")
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
	return nil
}

func (s *ScriptBuilder) pushData(data interface{}) error {
	switch e := data.(type) {
	case TransactionSignature:
		signatureLength := len(e.SignedData)
		b := []byte{}
		b = append(b, uintToBytes(uint(signatureLength))...)
		b = append(b, e.SignedData...)

		s.pushLength(len(b)) //this should be 0x41
		s.RawBytes = append(s.RawBytes, b...)
		log.Printf("b %x (%v)", b, len(b))
		s.RawBytes = append(s.RawBytes, 0x23) //0x23 = 35 this is the length of the next [publickey.length(2)]+[publickey(33)]]

		log.Printf("PublicKey %x (%v)", e.PublicKey, len(e.PublicKey))
		s.pushData(e.PublicKey)
		return nil
	case TransactionOutput:

		s.RawBytes = append(s.RawBytes, e.Asset.ToLittleEndianBytes()...) //32 bytes
		amountToSendBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(amountToSendBytes, uint64(e.Value))
		s.RawBytes = append(s.RawBytes, amountToSendBytes...) //8 bytes
		s.RawBytes = append(s.RawBytes, e.Address...)         //20 bytes
		return nil
	case UTXO:
		//reverse txID to little endian
		log.Printf("pusing utxo %v %v\n", e.TXID, e.Index)
		b, err := hex.DecodeString(e.TXID)
		if err != nil {
			return err
		}
		littleEndianTXID := reverseBytes(b)
		index := e.Index
		s.RawBytes = append(s.RawBytes, littleEndianTXID...)
		intBytes := uint16ToFixBytes(uint16(index))
		s.RawBytes = append(s.RawBytes, intBytes...)
		return nil
	case TradingVersion:
		s.RawBytes = append(s.RawBytes, byte(e))
		return nil
	case TransactionAttribute:
		s.RawBytes = append(s.RawBytes, byte(e))
		return nil
	case TransactionType:
		s.RawBytes = append(s.RawBytes, byte(e))
		return nil
	case NEOAddress:
		//when pushing neo address as an arg. we need length so we need to push a hex string
		return s.pushHexString(fmt.Sprintf("%x", e))
	case ScriptHash:
		s.RawBytes = append(s.RawBytes, e...)
		return nil
	case string:
		return s.pushHexString(e)
	case []byte:
		// length + data
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
		s.pushInt(count)
		s.pushOpCode(PACK)
		return nil
	case int:
		log.Printf("push int %v", int(e))
		s.pushInt(e)
		return nil
	case int64:
		log.Printf("push int64 %v", int(e))
		s.pushInt(int(e))
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

func (s ScriptHash) ToBigEndian() []byte {
	return reverseBytes([]byte(s))
}

// This is in a format of main(string operation, []object args) in c#
func (s *ScriptBuilder) generateContractInvocationData(scriptHash ScriptHash, operation string, args []interface{}) []byte {
	if args != nil {
		s.pushData(args)
	}
	s.pushData([]byte(operation))                                     //operation is in string we need to convert it to hex first
	s.pushOpCode(APPCALL)                                             //use APPCALL only
	s.pushData(scriptHash)                                            // script hash of the smart contract that we want to invoke
	s.RawBytes = append([]byte{byte(len(s.RawBytes))}, s.RawBytes...) //the length of the entire raw bytes
	return s.ToBytes()
}

func (s *ScriptBuilder) emptyTransactionAttributes() []byte {
	s.pushData(0x00)
	return s.ToBytes()
}

func (s *ScriptBuilder) generateTransactionAttributes(attributes map[TransactionAttribute][]byte) ([]byte, error) {

	count := len(attributes)
	s.pushLength(count) //number of transaction attributes
	// N x transaction attribute
	//transaction attribute =  TransactionAttribute + data.length + data
	for k, v := range attributes {
		s.pushData(k) //transaction attribute usage
		s.pushData(v) //push byte data in already includes the length of the data
	}

	return s.ToBytes(), nil
}

func (s *ScriptBuilder) generateTransactionInput(unspent Unspent, assetToSend NativeAsset, amountToSend float64) ([]byte, error) {
	//inputs = [input_count] + [[txID(32)] + [txIndex(2)]] = 34 x input_count bytes

	sendingAsset := unspent.Assets[assetToSend]
	if sendingAsset == nil {
		return nil, fmt.Errorf("Asset %v not found in UTXO", assetToSend)
	}

	if amountToSend > sendingAsset.TotalAmount() {
		return nil, fmt.Errorf("input Don't have enough balance. Sending %v but only have %v", amountToSend, sendingAsset.TotalAmount())
	}

	//sort min first
	sendingAsset.SortMinFirst()

	runningAmount := float64(0)
	index := 0
	count := 0
	inputs := []UTXO{}
	for runningAmount < amountToSend {
		addingUTXO := sendingAsset.UTXOs[index]
		inputs = append(inputs, addingUTXO)
		runningAmount += addingUTXO.Value
		index += 1
		count += 1
	}

	s.pushLength(count)
	for _, v := range inputs {
		//push utxo data
		s.pushData(v)
	}

	return s.ToBytes(), nil
}

func (s *ScriptBuilder) generateTransactionOutput(sender NEOAddress, receiver NEOAddress, unspent Unspent, assetToSend NativeAsset, amountToSend float64) ([]byte, error) {

	//output = [output_count] + [assetID(32)] + [amount(8)] + [sender_scripthash(20)] = 60 x output_count bytes

	sendingAsset := unspent.Assets[assetToSend]
	if sendingAsset == nil {
		return nil, fmt.Errorf("Asset %v not found in UTXO", assetToSend)
	}

	if amountToSend > sendingAsset.TotalAmount() {
		return nil, fmt.Errorf("output Don't have enough balance. Sending %v but only have %v", amountToSend, sendingAsset.TotalAmount())
	}
	//sort min first
	sendingAsset.SortMinFirst()

	runningAmount := float64(0)
	index := 0
	count := 0
	inputs := []UTXO{}
	for runningAmount < amountToSend {
		addingUTXO := sendingAsset.UTXOs[index]
		inputs = append(inputs, addingUTXO)
		runningAmount += addingUTXO.Value
		index += 1
		count += 1
	}

	//if the total amount of inputs is over amountToSend
	//we need to send the rest back to the sending address
	totalAmountInInputs := runningAmount

	needTwoOutputTransaction := totalAmountInInputs != amountToSend

	list := []TransactionOutput{}
	log.Printf("needTwoOutputTransaction = %v", needTwoOutputTransaction)
	if needTwoOutputTransaction {
		sendingOutput := TransactionOutput{
			Asset:   assetToSend,
			Value:   int64(amountToSend * float64(100000000)),
			Address: receiver,
		}
		list = append(list, sendingOutput)

		returningAmount := totalAmountInInputs - amountToSend

		//return the left over to sender
		returningOutput := TransactionOutput{
			Asset:   assetToSend,
			Value:   int64(returningAmount * float64(100000000)),
			Address: sender,
		}
		list = append(list, returningOutput)
	} else {
		out := TransactionOutput{
			Asset:   assetToSend,
			Value:   int64(amountToSend * float64(100000000)),
			Address: receiver,
		}
		list = append(list, out)
	}

	//number of outputs
	s.pushLength(len(list))
	for _, v := range list {
		s.pushData(v)
	}

	return s.ToBytes(), nil
}

func (s *ScriptBuilder) generateInvocationScriptWithSignatures(signatures []TransactionSignature) []byte {

	numberOfSignatures := len(signatures)
	if numberOfSignatures == 0 {
		return nil
	}

	s.pushLength(numberOfSignatures)

	for _, signature := range signatures {
		s.pushData(signature)
	}

	if numberOfSignatures >= 1 {
		s.pushOpCode(CHECKSIG)
	} else {
		s.pushOpCode(CHECKMULTISIG)
	}
	log.Printf("signature data %x", s.ToBytes())
	return s.ToBytes()
}
