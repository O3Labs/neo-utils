package smartcontract

//naming base on NEO network protocol
//http://docs.neo.org/en-us/network/network-protocol.html
type TransactionValidationScript struct {
	StackScript  []byte
	RedeemScript []byte
}
