package smartcontract

type Transaction struct {
	Type       TransactionType
	Version    TradingVersion
	Data       []byte
	Attributes []byte
	Inputs     []byte
	Outputs    []byte
	Script     []byte
}

type TransactionOutput struct {
	Asset   NativeAsset
	Value   int64
	Address NEOAddress
}

func (t *Transaction) ToBytes() []byte {

	return nil
}

func (t *Transaction) NewInvocationTransaction() Transaction {
	return Transaction{
		Type:    InvocationTransaction,
		Version: NEOTradingVersion,
	}
}
