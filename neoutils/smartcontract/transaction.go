package smartcontract

import (
	"crypto/sha256"
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

//version is 0 currently
//it needs to change to 1 eventually to support pay gas to run smart contract
//https://github.com/neo-project/neo/blob/11d8db11568d9eadeeb86c5b8c21a1d3937e0912/neo/Core/InvocationTransaction.cs#L23
func NewInvocationTransaction() Transaction {
	return Transaction{
		Type:    InvocationTransaction,
		Version: NEOTradingVersion,
	}
}

func NewContractTransaction() Transaction {
	return Transaction{
		Type:    ContractTransaction,
		Version: NEOTradingVersion,
	}
}
