package neoutils

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"

	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func WriteVarUint(w io.Writer, val uint64) error {
	// if val < 0xfd {
	// 	binary.Write(w, binary.LittleEndian, uint8(val))
	// 	return nil
	// }
	// if val < 0xFFFF {
	// 	binary.Write(w, binary.LittleEndian, byte(0xfd))
	// 	binary.Write(w, binary.LittleEndian, uint16(val))
	// 	return nil
	// }
	// if val < 0xFFFFFFFF {
	// 	binary.Write(w, binary.LittleEndian, byte(0xfe))
	// 	binary.Write(w, binary.LittleEndian, uint32(val))
	// 	return nil
	// }

	// binary.Write(w, binary.LittleEndian, byte(0xff))
	// binary.Write(w, binary.LittleEndian, val)

	// return nil
	if val < 0xfd {
		return binary.Write(w, binary.LittleEndian, uint8(val))
	}
	if val < 0xFFFF {
		if err := binary.Write(w, binary.LittleEndian, byte(0xfd)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, uint16(val))
	}
	if val < 0xFFFFFFFF {
		if err := binary.Write(w, binary.LittleEndian, byte(0xfe)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, uint32(val))
	}

	if err := binary.Write(w, binary.LittleEndian, byte(0xff)); err != nil {
		return err
	}

	return binary.Write(w, binary.LittleEndian, val)
}

type SmartContractInfo struct {
	AVMHEX      string
	Name        string
	Version     string
	Author      string
	Email       string
	Description string
	Properties  smartcontract.Properties
	InputTypes  []smartcontract.ParameterType
	ReturnType  smartcontract.ParameterType
}

func (s *SmartContractInfo) GetScriptHash() string {
	address := VMCodeToNEOAddress(hex2bytes(s.AVMHEX))
	scripthash := NEOAddressToScriptHashWithEndian(address, binary.BigEndian)
	return scripthash
}

func (s *SmartContractInfo) Serialize() []byte {

	params := []byte{}
	for _, p := range s.InputTypes {
		params = append(params, p.Byte())
	}

	scriptBuilder := smartcontract.NewScriptBuilder()
	scriptBuilder.Push([]byte(s.Description))
	scriptBuilder.Push([]byte(s.Email))
	scriptBuilder.Push([]byte(s.Author))
	scriptBuilder.Push([]byte(s.Version))
	scriptBuilder.Push([]byte(s.Name))
	scriptBuilder.Push(int(s.Properties))
	scriptBuilder.Push([]byte{s.ReturnType.Byte()})
	scriptBuilder.Push(params)
	scriptBuilder.Push(hex2bytes(s.AVMHEX))
	scriptBuilder.PushSysCall("Neo.Contract.Create")

	b := scriptBuilder.ToBytes()
	buff := new(bytes.Buffer)
	WriteVarUint(buff, uint64(len(b)))
	endPayload := []byte{}
	endPayload = append(endPayload, buff.Bytes()...)
	endPayload = append(endPayload, b...)
	return endPayload
}

func DeploySmartContractScript(contractInfo SmartContractInfo, wallet Wallet, asset smartcontract.NativeAsset, amount float64, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte) ([]byte, error) {

	tx := smartcontract.NewInvocationTransactionPayable()

	tx.Data = contractInfo.Serialize()
	tx.GAS = uint64(amount)

	amountToSend := amount
	assetToSend := asset

	networkFee := smartcontract.NetworkFeeAmount(0)

	txInputs, err := smartcontract.NewScriptBuilder().GenerateTransactionInput(unspent, assetToSend, amountToSend, networkFee)
	if err != nil {
		return nil, err
	}
	tx.Inputs = txInputs

	txAttributes, err := smartcontract.NewScriptBuilder().GenerateTransactionAttributes(attributes)
	if err != nil {
		return nil, err
	}
	tx.Attributes = txAttributes

	//when deploy smart contract, you don't actually send asset to another address
	//so the receiver is the same address
	sender := smartcontract.ParseNEOAddress(wallet.Address)
	receiver := smartcontract.ParseNEOAddress(wallet.Address)
	txOutputs, err := smartcontract.NewScriptBuilder().GenerateTransactionOutputPayableGAS(sender, receiver, unspent, assetToSend, amount, networkFee, float64(tx.GAS))
	if err != nil {
		return nil, err
	}

	tx.Outputs = txOutputs

	//begin signing process and invocation script
	privateKeyInHex := bytesToHex(wallet.PrivateKey)
	signedData, err := Sign(tx.ToBytes(), privateKeyInHex)
	if err != nil {
		return nil, err
	}

	signature := smartcontract.TransactionSignature{
		SignedData: signedData,
		PublicKey:  wallet.PublicKey,
	}

	scripts := []interface{}{signature}
	txScripts := smartcontract.NewScriptBuilder().GenerateVerificationScripts(scripts)
	tx.Script = txScripts
	//end signing process

	log.Printf("txid = %v", tx.ToTXID())

	return tx.ToBytes(), nil
}
