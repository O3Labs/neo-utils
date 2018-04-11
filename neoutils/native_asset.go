package neoutils

import (
	"log"

	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func SendNativeAssetRawTransaction(wallet Wallet, asset smartcontract.NativeAsset, amount float64, to smartcontract.NEOAddress, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte) ([]byte, error) {
	//New invocation transaction struct and fill with all necessary data
	tx := smartcontract.NewContractTransaction()

	amountToSend := amount
	assetToSend := asset

	// fee := float64(0.00000001)
	//generate transaction inputs
	txInputs, err := smartcontract.NewScriptBuilder().GenerateTransactionInput(unspent, assetToSend, amountToSend)
	if err != nil {
		return nil, err
	}
	//transaction inputs
	tx.Inputs = txInputs

	//generate transaction outputs
	txAttributes, err := smartcontract.NewScriptBuilder().GenerateTransactionAttributes(attributes)
	if err != nil {
		return nil, err
	}
	//transaction attributes
	tx.Attributes = txAttributes

	//send GAS to the same account
	sender := smartcontract.ParseNEOAddress(wallet.Address)
	receiver := to
	txOutputs, err := smartcontract.NewScriptBuilder().GenerateTransactionOutput(sender, receiver, unspent, assetToSend, amountToSend)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	tx.Outputs = txOutputs

	//begin signing process and invocation script
	privateKeyInHex := bytesToHex(wallet.PrivateKey)

	signedData, err := Sign(tx.ToBytes(), privateKeyInHex)
	if err != nil {
		log.Printf("err signing %v", err)
		return nil, err
	}

	signature := smartcontract.TransactionSignature{
		SignedData: signedData,
		PublicKey:  wallet.PublicKey,
	}

	signatures := []smartcontract.TransactionSignature{signature}
	txScripts := smartcontract.NewScriptBuilder().GenerateInvocationAndVerificationScriptWithSignatures(signatures)
	//assign scripts to the tx
	tx.Script = txScripts
	//end signing process

	//concat data
	endPayload := []byte{}
	endPayload = append(endPayload, tx.ToBytes()...)

	return endPayload, nil
}