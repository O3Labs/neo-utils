package smartcontract

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

type Transaction struct {
	Type       TransactionType
	Version    TradingVersion
	Data       []byte
	Attributes []byte
	Inputs     []byte
	Outputs    []byte
	//scripts contains two parts, Invocation script and Verification script
	Script []byte
	GAS    uint64 //only for version 1
}

type TransactionOutput struct {
	Asset   NativeAsset
	Value   int64
	Address NEOAddress
}

func (t *Transaction) ToBytes() []byte {
	payload := []byte{}
	payload = append(payload, byte(t.Type))
	payload = append(payload, byte(t.Version))
	payload = append(payload, t.Data...)

	if t.Version >= NEOTradingVersionPayableGAS {
		gasInBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(gasInBytes, uint64(t.GAS*uint64(100000000)))
		payload = append(payload, gasInBytes...)
	}

	payload = append(payload, t.Attributes...)
	payload = append(payload, t.Inputs...)
	payload = append(payload, t.Outputs...)
	payload = append(payload, t.Script...)

	return payload
}

//this ToHash256 returns little endian bytes.
//TXID is big endian bytes, so when calling json-rpc api we need to reverse it
func (t *Transaction) ToHash256() []byte {
	payload := []byte{}
	payload = append(payload, byte(t.Type))
	payload = append(payload, byte(t.Version))
	payload = append(payload, t.Data...)

	if t.Version >= NEOTradingVersionPayableGAS {
		gasInBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(gasInBytes, uint64(t.GAS*uint64(100000000)))
		payload = append(payload, gasInBytes...)
	}
	payload = append(payload, t.Attributes...)
	payload = append(payload, t.Inputs...)
	payload = append(payload, t.Outputs...)

	hash := sha256.Sum256(payload)
	hash = sha256.Sum256(hash[:])

	return hash[:]
}

func (t *Transaction) ToTXID() string {
	return fmt.Sprintf("%x", reverseBytes(t.ToHash256()))
}

func NewInvocationTransaction() Transaction {
	return Transaction{
		Type:    InvocationTransaction,
		Version: NEOTradingVersion,
	}
}

func NewInvocationTransactionPayable() Transaction {
	return Transaction{
		Type:    InvocationTransaction,
		Version: NEOTradingVersionPayableGAS, //version 1 this will allow paying GAS
	}
}

func NewContractTransaction() Transaction {
	return Transaction{
		Type:    ContractTransaction,
		Version: NEOTradingVersion,
	}
}

func NewTransactionWithType(txType byte, version int) Transaction {
	if txType == byte(InvocationTransaction) {
		return Transaction{
			Type:    InvocationTransaction,
			Version: TradingVersion(version),
		}
	}
	return Transaction{
		Type:    InvocationTransaction,
		Version: NEOTradingVersionPayableGAS, //version 1 this will allow paying GAS
	}
}
