package smartcontract

//naming base on NEO network protocol
//http://docs.neo.org/en-us/network/network-protocol.html
type TransactionValidationScript struct {
	Invocation   interface{}
	Verification []byte
}
