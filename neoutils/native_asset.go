package neoutils

import (
	"log"

	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

type NativeAssetInterface interface {
	SendNativeAssetRawTransaction(wallet Wallet, asset smartcontract.NativeAsset, amount float64, to smartcontract.NEOAddress, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte) ([]byte, string, error)
	GenerateRawTx(fromAddress string, asset smartcontract.NativeAsset, amount float64, to smartcontract.NEOAddress, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte) ([]byte, string, error)
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

func (n *NativeAsset) SendNativeAssetRawTransaction(wallet Wallet, asset smartcontract.NativeAsset, amount float64, to smartcontract.NEOAddress, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte) ([]byte, string, error) {
	tx, txID, err := n.GenerateRawTx(wallet.Address, asset, amount, to, unspent, attributes)
	if err != nil {
		return nil, "", err
	}

	//begin signing
	privateKeyInHex := bytesToHex(wallet.PrivateKey)
	signedData, err := Sign(tx, privateKeyInHex)
	if err != nil {
		log.Printf("err signing %v", err)
		return nil, "", err
	}

	signature := smartcontract.TransactionSignature{
		SignedData: signedData,
		PublicKey:  wallet.PublicKey,
	}

	scripts := []interface{}{signature}
	txScripts := smartcontract.NewScriptBuilder().GenerateVerificationScripts(scripts)

	//concat data
	endPayload := []byte{}
	endPayload = append(endPayload, tx...)
	endPayload = append(endPayload, txScripts...)

	return endPayload, txID, nil
}

func (n *NativeAsset) GenerateRawTx(fromAddress string, asset smartcontract.NativeAsset, amount float64, to smartcontract.NEOAddress, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte) ([]byte, string, error) {
	//New invocation transaction struct and fill with all necessary data
	tx := smartcontract.NewContractTransaction()

	//generate transaction inputs
	txInputs, err := smartcontract.NewScriptBuilder().GenerateTransactionInput(unspent, asset, amount, n.NetworkFeeAmount)
	if err != nil {
		return nil, "", err
	}
	tx.Inputs = txInputs

	txAttributes, err := smartcontract.NewScriptBuilder().GenerateTransactionAttributes(attributes)
	if err != nil {
		return nil, "", err
	}

	tx.Attributes = txAttributes

	sender := smartcontract.ParseNEOAddress(fromAddress)

	txOutputs, err := smartcontract.NewScriptBuilder().GenerateTransactionOutput(sender, to, unspent, asset, amount, n.NetworkFeeAmount)
	if err != nil {
		log.Printf("%v", err)
		return nil, "", err
	}

	tx.Outputs = txOutputs

	return tx.ToBytes(), tx.ToTXID(), nil
}
