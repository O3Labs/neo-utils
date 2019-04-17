package neoutils

import "github.com/o3labs/neo-utils/neoutils/smartcontract"

var cgasScriptHash = "74f2dc36a68fdc4682034178eb2220729231db76"

func GASToCGAS(wallet Wallet, amountToSend float64, networkFee float64) ([]byte, string, error) {
	cgas, _ := smartcontract.NewScriptHash(cgasScriptHash)
	network := "test"
	operation := "mintTokens"
	remark := "O3XGAS2CGAS"
	networkFeeAmount := smartcontract.NetworkFeeAmount(networkFee)
	assetToSend := smartcontract.GAS

	args := []interface{}{}
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	attributes[smartcontract.Remark1] = []byte(remark)

	tx := smartcontract.NewInvocationTransaction()
	txData := smartcontract.NewScriptBuilder().GenerateContractInvocationData(cgas, operation, args)
	tx.Data = txData

	unspent, err := utxoFromO3Platform(network, wallet.Address)
	if err != nil {
		return nil, "", err
	}

	txInputs, err := smartcontract.NewScriptBuilder().GenerateTransactionInput(unspent, assetToSend, amountToSend, networkFeeAmount)
	if err != nil {
		return nil, "", err
	}

	tx.Inputs = txInputs

	txAttributes, err := smartcontract.NewScriptBuilder().GenerateTransactionAttributes(attributes)
	if err != nil {
		return nil, "", err
	}
	tx.Attributes = txAttributes

	sender := smartcontract.ParseNEOAddress(wallet.Address)
	receiver := smartcontract.NEOAddressFromScriptHash(cgas.ToBigEndian())
	txOutputs, err := smartcontract.NewScriptBuilder().GenerateTransactionOutput(sender, receiver, unspent, assetToSend, amountToSend, networkFeeAmount)
	if err != nil {
		return nil, "", err
	}

	tx.Outputs = txOutputs

	privateKeyInHex := bytesToHex(wallet.PrivateKey)

	signedData, err := Sign(tx.ToBytes(), privateKeyInHex)
	if err != nil {
		return nil, "", err
	}

	signature := smartcontract.TransactionSignature{
		SignedData: signedData,
		PublicKey:  wallet.PublicKey,
	}

	scripts := []interface{}{signature}

	txScripts := smartcontract.NewScriptBuilder().GenerateVerificationScripts(scripts)

	tx.Script = txScripts

	return tx.ToBytes(), tx.ToTXID(), nil
}

func CGASToGAS(wallet Wallet, amount float64, networkFee float64) ([]byte, string, error) {
	// cgas, _ := smartcontract.NewScriptHash(cgasScriptHash)
	// network := "test"
	// operation := "refund"
	// remark := "O3XCGAS2GAS"
	// networkFeeAmount := smartcontract.NetworkFeeAmount(networkFee)
	// assetToSend := smartcontract.GAS

	// //cgas contract utxo
	// unspent, err := utxoFromO3Platform(network, cgas.ToNEOAddress())

	return nil, "", nil
}
