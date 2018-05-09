package neoutils

const (
	VERSION = "1.0.4"
)

//RELEASE NOTES
// V. 1.0.4
// - Updated to use UTXO from O3

// V. 1.0.3
// - mintTokens now triggers Verification

// V. 1.0.2
// - added txid in return
// - MintTokensRawTransaction(wallet Wallet, assetToSend smartcontract.NativeAsset, amount float64, unspent smartcontract.Unspent, remark string) ([]byte, string, error)
// - TransferNEP5RawTransaction(wallet Wallet, toAddress smartcontract.NEOAddress, amount float64, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte) ([]byte, string, error)
