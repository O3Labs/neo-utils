package smartcontract

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/o3labs/neo-utils/neoutils/btckey"
	"golang.org/x/crypto/ripemd160"
)

type ScriptHash []byte
type NEOAddress []byte

type TokenAmount uint

const (
	Uint160Length = 20
)

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
func NEOAddressFromScriptHash(scriptHashBytes []byte) NEOAddress {
	address := btckey.B58checkencodeNEO(0x17, reverseBytes(scriptHashBytes))
	return ParseNEOAddress(address)
}

func (n NEOAddress) ToString() string {
	return btckey.B58checkencodeNEO(0x17, n)
}

type ScriptBuilderInterface interface {
	GenerateContractInvocationData(scriptHash ScriptHash, operation string, args []interface{}) []byte
	GenerateTransactionAttributes(attributes map[TransactionAttribute][]byte) ([]byte, error)

	//this is to send the UTXO of asset that will be used in TransactionOutput
	GenerateTransactionInput(unspent Unspent, assetToSend NativeAsset, amountToSend float64, networkFeeAmount NetworkFeeAmount) ([]byte, error)
	GenerateTransactionOutput(sender NEOAddress, receiver NEOAddress, unspent Unspent, assetToSend NativeAsset, amountToSend float64, networkFeeAmount NetworkFeeAmount) ([]byte, error)

	GenerateVerificationScripts(signatures []interface{}) []byte
	EmptyTransactionAttributes() []byte
	ToBytes() []byte
	FullHexString() string
	Clear()

	EmitPush(data interface{}) error
	ToScriptHash() []byte //UInt160

	pushInt(value int) error
	pushData(data interface{}) error
	pushLength(count int)
}

func NewScriptBuilder() ScriptBuilderInterface {
	return &ScriptBuilder{RawBytes: []byte{}}
}

type ScriptBuilder struct {
	RawBytes []byte
}

func (s *ScriptBuilder) ToScriptHash() []byte {
	sha := sha256.New()
	sha.Write(s.ToBytes())
	b := sha.Sum(nil)
	ripemd := ripemd160.New()
	ripemd.Write(b)
	b = ripemd.Sum(nil)
	return b[0:Uint160Length]
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
func (s *ScriptBuilder) pushInt8bytes(value int) error {
	num := make([]byte, 8)
	binary.LittleEndian.PutUint64(num, uint64(value))
	return s.pushData(num)
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
		//we push as []byte so then it prefixes with length
		s.pushData(bytes.Trim(num, "\x00"))
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

func (s *ScriptBuilder) EmitPush(data interface{}) error {
	return s.pushData(data)
}

func (s *ScriptBuilder) pushData(data interface{}) error {
	switch e := data.(type) {
	case TransactionValidationScript:
		s.pushData(e.StackScript)
		if e.RedeemScript == nil {
			s.RawBytes = append(s.RawBytes, 0x00)
		} else {
			s.pushData(e.RedeemScript)
		}
		return nil
	case TransactionSignature:
		signatureLength := len(e.SignedData)
		b := []byte{}
		b = append(b, uintToBytes(uint(signatureLength))...)
		b = append(b, e.SignedData...)
		s.pushLength(len(b)) //this should be 0x41
		s.RawBytes = append(s.RawBytes, b...)
		s.RawBytes = append(s.RawBytes, 0x23) //0x23 = 35 this is the length of the next [publickey.length(2)]+[publickey(33)]]
		//this part is for verification script
		//push public key in there and call CHECKSIG or CHECKMULTISIG
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
		//remove prefix 0x here
		//check if the scripthash is prefixed with 0x. if so, trim it out.
		trimmed0x := e.TXID
		if has0xPrefix(e.TXID) == true {
			trimmed0x = e.TXID[2:]
		}
		//reverse txID to little endian
		b, err := hex.DecodeString(trimmed0x)
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
			log.Printf("case []interface{}: %x", e[i])
			s.pushData(e[i])
		}
		s.pushInt(count)
		s.pushOpCode(PACK)
		return nil
	case int:
		log.Printf("push int %v", e)
		s.pushInt(e)
		return nil
	case int64:
		s.pushInt(int(e))
		return nil
	case TokenAmount:
		log.Printf("token amount =%v", e)
		s.pushInt8bytes(int(e))
		return nil
	}
	log.Printf("unknown type %v", data)
	return nil
}

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func NewScriptHash(hexString string) (ScriptHash, error) {
	//check if the scripthash is prefixed with 0x. if so, trim it out.
	trimmed0x := hexString
	if has0xPrefix(hexString) == true {
		trimmed0x = hexString[2:]
	}
	b, err := hex.DecodeString(trimmed0x)
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
func (s *ScriptBuilder) GenerateContractInvocationData(scriptHash ScriptHash, operation string, args []interface{}) []byte {
	if args != nil {
		s.pushData(args)
	}
	s.pushData([]byte(operation))                                     //operation is in string we need to convert it to hex first
	s.pushOpCode(APPCALL)                                             //use APPCALL only
	s.pushData(scriptHash)                                            //script hash of the smart contract that we want to invoke
	s.RawBytes = append([]byte{byte(len(s.RawBytes))}, s.RawBytes...) //the length of the entire raw bytes
	return s.ToBytes()
}

func (s *ScriptBuilder) EmptyTransactionAttributes() []byte {
	s.pushData(0x00)
	return s.ToBytes()
}

func (s *ScriptBuilder) GenerateTransactionAttributes(attributes map[TransactionAttribute][]byte) ([]byte, error) {

	count := len(attributes)
	s.pushLength(count) //number of transaction attributes
	// N x transaction attribute
	//transaction attribute =  TransactionAttribute + data.length + data
	for k, v := range attributes {
		s.pushData(k) //transaction attribute usage
		if k == Script {
			//if it's a Script field, we just need to put it as is in a little endian bytes
			s.pushData(ScriptHash(v))
		} else {
			s.pushData(v) //push byte data in already includes the length of the data
		}
	}

	return s.ToBytes(), nil
}

func (s *ScriptBuilder) GenerateTransactionInput(unspent Unspent, assetToSend NativeAsset, amountToSend float64, networkFeeAmount NetworkFeeAmount) ([]byte, error) {
	//inputs = [input_count] + [[txID(32)] + [txIndex(2)]] = 34 x input_count bytes

	//empty unspent
	if len(unspent.Assets) == 0 || amountToSend == 0 {
		s.pushLength(0)
		return s.ToBytes(), nil
	}
	sendingAsset := unspent.Assets[assetToSend]
	if sendingAsset == nil {
		return nil, fmt.Errorf("Asset %v not found in UTXO", assetToSend)
	}
	//network fee
	feeAmount := networkFeeAmount

	//if assetToSend is NEO and fee amount is more than zero
	needAnotherAssetForFee := false
	if assetToSend == NEO && feeAmount > 0 {
		//we need another input because fee is in GAS
		needAnotherAssetForFee = true
	}

	if amountToSend > sendingAsset.TotalAmount() {
		return nil, fmt.Errorf("input Don't have enough balance. Sending %v but only have %v", amountToSend, sendingAsset.TotalAmount())
	}

	//sort min first
	sendingAsset.SortMinFirst()

	utxoSumAmount := float64(0)
	index := 0
	count := 0
	inputs := []UTXO{}
	//loop until we get enough sum amount
	for utxoSumAmount < amountToSend {
		addingUTXO := sendingAsset.UTXOs[index]
		log.Printf("%+v", addingUTXO)
		inputs = append(inputs, addingUTXO)
		utxoSumAmount += addingUTXO.Value
		index += 1
		count += 1
	}

	//fee input part
	if needAnotherAssetForFee == true {
		gasBalanceForFee := unspent.Assets[GAS]
		gasBalanceForFee.SortMinFirst()
		if float64(feeAmount) > gasBalanceForFee.TotalAmount() {
			return nil, fmt.Errorf("you don't have enough balance for network fee.")
		}
		utxoSumFeeAmount := float64(0)
		feeIndex := 0
		for utxoSumFeeAmount < float64(feeAmount) {
			addingUTXO := gasBalanceForFee.UTXOs[feeIndex]
			inputs = append(inputs, addingUTXO)
			utxoSumFeeAmount += addingUTXO.Value
			feeIndex += 1
			count += 1
		}
		//end fee input part
	}

	s.pushLength(count)
	for _, v := range inputs {
		//push utxo data
		s.pushData(v)
	}

	return s.ToBytes(), nil
}

func (s *ScriptBuilder) GenerateTransactionOutput(sender NEOAddress, receiver NEOAddress, unspent Unspent, assetToSend NativeAsset, amountToSend float64, networkFeeAmount NetworkFeeAmount) ([]byte, error) {

	//output = [output_count] + [assetID(32)] + [amount(8)] + [sender_scripthash(20)] = 60 x output_count bytes
	//empty unspent
	if len(unspent.Assets) == 0 || amountToSend == 0 {
		log.Printf("unspent is empty")
		s.pushLength(0)
		return s.ToBytes(), nil
	}

	sendingAsset := unspent.Assets[assetToSend]
	if sendingAsset == nil {
		return nil, fmt.Errorf("Asset %v not found in UTXO", assetToSend)
	}

	//network fee
	feeAmount := networkFeeAmount

	//if assetToSend is NEO and fee amount is more than zero
	needAnotherAssetForFee := false
	if assetToSend == NEO && feeAmount > 0 {
		needAnotherAssetForFee = true
	}

	if amountToSend > sendingAsset.TotalAmount() {
		return nil, fmt.Errorf("you don't have enough balance. Sending %v but only have %v", amountToSend, sendingAsset.TotalAmount())
	}
	//sort min first
	sendingAsset.SortMinFirst()

	utxoSumAmount := float64(0)
	index := 0
	count := 0
	inputs := []UTXO{}
	for utxoSumAmount < amountToSend {
		addingUTXO := sendingAsset.UTXOs[index]
		inputs = append(inputs, addingUTXO)
		utxoSumAmount += addingUTXO.Value
		index += 1
		count += 1
	}

	//if the total amount of inputs is over amountToSend
	//we need to send the rest back to the sending address
	totalAmountInInputs := utxoSumAmount
	needTwoOutputTransaction := totalAmountInInputs != amountToSend
	list := []TransactionOutput{}

	if needTwoOutputTransaction {
		//first output is the amount to send to the receiver
		sendingOutput := TransactionOutput{
			Asset:   assetToSend,
			Value:   int64(amountToSend * float64(100000000)),
			Address: receiver,
		}
		list = append(list, sendingOutput)

		//second output is the returning amount you will be sending back to yourself.
		returningAmount := totalAmountInInputs - amountToSend

		//so if we don't need another asset input and fee is more than 0
		//we then make returningAmount = returningAmount - fee
		if needAnotherAssetForFee == false && float64(feeAmount) > 0 {
			returningAmount -= float64(feeAmount)
		}
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

	//if set network fee is more than 0
	//add more output for fee
	if needAnotherAssetForFee == true {

		gasBalanceForFee := unspent.Assets[GAS]
		gasBalanceForFee.SortMinFirst()
		if float64(feeAmount) > gasBalanceForFee.TotalAmount() {
			return nil, fmt.Errorf("you don't have enough balance for network fee.")
		}
		runningFeeAmount := float64(0)
		feeIndex := 0
		for runningFeeAmount < float64(feeAmount) {
			addingUTXO := gasBalanceForFee.UTXOs[feeIndex]
			inputs = append(inputs, addingUTXO)
			runningFeeAmount += addingUTXO.Value
			feeIndex += 1
			count += 1
		}

		// To allow user to set network fee is to make send GAS back to yourself
		// minus the amount of gas that you want it to be network fee
		// for example
		// GAS balance = 10
		// sending back amount = 9
		// this will make network fee = 1
		returningAmount := runningFeeAmount - float64(feeAmount)
		returningOutput := TransactionOutput{
			Asset:   GAS,
			Value:   int64(returningAmount * float64(100000000)),
			Address: sender,
		}
		list = append(list, returningOutput)
	}

	//number of outputs
	s.pushLength(len(list))
	for _, v := range list {
		s.pushData(v)
	}

	return s.ToBytes(), nil
}

func (s *ScriptBuilder) GenerateVerificationScripts(scripts []interface{}) []byte {

	numberOfScripts := len(scripts)
	if numberOfScripts == 0 {
		return nil
	}

	s.pushLength(numberOfScripts)
	for _, script := range scripts {
		switch e := script.(type) {
		case TransactionSignature:
			s.pushData(e)
			s.pushOpCode(CHECKSIG)
			continue
		case TransactionValidationScript:
			s.pushData(e)
		}
	}
	return s.ToBytes()
}
