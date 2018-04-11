package smartcontract

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
