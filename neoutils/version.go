package neoutils

const (
	VERSION = "1.3.0"
)

//RELEASE NOTES

// V.1.3.0
// - Sample Deploy Smart Contract without syncing
// - Fixed add empty string as a argument when invoke script
// - Sample NNS invoke script
// - Sample Multisignature address

// V.1.2.2
// - Allow sending NEP5 with fee

// V.1.2.1
// - Added ClaimONG

// V.1.2.0
// - Added generate multisignature address redeem script
// - Added gasPrice and gasLimit to ontology transfer method

// V. 1.1.0
// - Added Ontology asset transfer

// V. 1.0.8
// - Added Generate NEP6 wallet for mobile

// V. 1.0.7
// - Make sure to round to fixed 8 decimals in output

// V. 1.0.6
// - Added Verify method to verify signed data

// V. 1.0.5
// - Added NEP6 Wallet format
// - Make https handshake timeout lower to make get best node faster
// - Added generate invocation script

// V. 1.0.4
// - Updated to use UTXO from O3

// V. 1.0.3
// - mintTokens now triggers Verification

// V. 1.0.2
// - added txid in return
// - MintTokensRawTransaction(wallet Wallet, assetToSend smartcontract.NativeAsset, amount float64, unspent smartcontract.Unspent, remark string) ([]byte, string, error)
// - TransferNEP5RawTransaction(wallet Wallet, toAddress smartcontract.NEOAddress, amount float64, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte) ([]byte, string, error)
