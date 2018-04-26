package neoutils

import (
	"log"

	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

type NativeAssetInterface interface {
	SendNativeAssetRawTransaction(wallet Wallet, asset smartcontract.NativeAsset, amount float64, to smartcontract.NEOAddress, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte) ([]byte, error)
}

type NativeAsset struct {
	NetworkFeeAmount smartcontract.NetworkFeeAmount //allow users to override the network fee here
}

func UseNativeAsset(networkFeeAmount smartcontract.NetworkFeeAmount) NativeAsset {
	return NativeAsset{
		NetworkFeeAmount: networkFeeAmount,
	}
}

var _ NativeAssetInterface = (*NativeAsset)(nil)

func (n *NativeAsset) SendNativeAssetRawTransaction(wallet Wallet, asset smartcontract.NativeAsset, amount float64, to smartcontract.NEOAddress, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte) ([]byte, error) {
	//New invocation transaction struct and fill with all necessary data
	tx := smartcontract.NewContractTransaction()

	amountToSend := amount
	assetToSend := asset

	//generate transaction inputs
	txInputs, err := smartcontract.NewScriptBuilder().GenerateTransactionInput(unspent, assetToSend, amountToSend, n.NetworkFeeAmount)
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

	sender := smartcontract.ParseNEOAddress(wallet.Address)
	receiver := to
	txOutputs, err := smartcontract.NewScriptBuilder().GenerateTransactionOutput(sender, receiver, unspent, assetToSend, amountToSend, n.NetworkFeeAmount)
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

	scripts := []interface{}{signature}
	txScripts := smartcontract.NewScriptBuilder().GenerateVerificationScripts(scripts)
	//assign scripts to the tx
	tx.Script = txScripts
	//end signing process

	//concat data
	endPayload := []byte{}
	endPayload = append(endPayload, tx.ToBytes()...)

	return endPayload, nil
}
