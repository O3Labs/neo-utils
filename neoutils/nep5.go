package neoutils

import (
	"fmt"
	"strings"

	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

type NEP5Interface interface {
	MintTokensRawTransaction(wallet Wallet, assetToSend smartcontract.NativeAsset, amount float64, unspent smartcontract.Unspent, remark string) ([]byte, error)
	TransferNEP5RawTransaction(wallet Wallet, toAddress smartcontract.NEOAddress, amount float64, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte) ([]byte, error)
}

type NEP5 struct {
	ScriptHash       smartcontract.ScriptHash
	NetworkFeeAmount smartcontract.NetworkFeeAmount //allow users to override the network fee here
}

func UseNEP5WithNetworkFee(scriptHashHex string, networkFeeAmount smartcontract.NetworkFeeAmount) *NEP5 {
	if len(strings.TrimSpace(scriptHashHex)) == 0 {
		return nil
	}
	scripthash, err := smartcontract.NewScriptHash(scriptHashHex)
	if err != nil {
		return nil
	}

	return &NEP5{ScriptHash: scripthash, NetworkFeeAmount: networkFeeAmount}
}

var _ NEP5Interface = (*NEP5)(nil)

func (n *NEP5) TransferNEP5RawTransaction(wallet Wallet, toAddress smartcontract.NEOAddress, amount float64, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte) ([]byte, error) {

	from := smartcontract.ParseNEOAddress(wallet.Address)
	if from == nil {
		return nil, fmt.Errorf("Invalid from address")
	}

	to := smartcontract.ParseNEOAddress(toAddress.ToString())
	if to == nil {
		return nil, fmt.Errorf("Invalid from address")
	}
	numberOfTokens := amount
	args := []interface{}{from, to, numberOfTokens}

	//New invocation transaction struct and fill with all necessary data
	tx := smartcontract.NewInvocationTransaction()
	txData := smartcontract.NewScriptBuilder().GenerateContractInvocationData(n.ScriptHash, "transfer", args)
	tx.Data = txData

	//for smart contract invocation we send the minimum amount of gas to it
	//0.00000001 gas
	amountToSend := float64(0.00000001)
	assetToSend := smartcontract.GAS

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

	//send GAS to the same account
	sender := smartcontract.ParseNEOAddress(wallet.Address)
	receiver := smartcontract.ParseNEOAddress(wallet.Address)
	txOutputs, err := smartcontract.NewScriptBuilder().GenerateTransactionOutput(sender, receiver, unspent, assetToSend, amountToSend, n.NetworkFeeAmount)
	if err != nil {
		return nil, err
	}

	tx.Outputs = txOutputs

	//begin signing process and invocation script
	privateKeyInHex := bytesToHex(wallet.PrivateKey)

	signedData, err := Sign(tx.ToBytes(), privateKeyInHex)
	if err != nil {
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
	endPayload = append(endPayload, n.ScriptHash.ToBigEndian()...)

	return endPayload, nil
}

func (n *NEP5) MintTokensRawTransaction(wallet Wallet, assetToSend smartcontract.NativeAsset, amount float64, unspent smartcontract.Unspent, remark string) ([]byte, error) {

	scString := fmt.Sprintf("%x", n.ScriptHash.ToBigEndian())
	//this is because the script hash object is already in little endian
	sc := UseSmartContractWithNetworkFee(scString, n.NetworkFeeAmount)
	operation := "mintTokens"
	args := []interface{}{}
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	attributes[smartcontract.Remark1] = []byte(remark)

	tx, err := sc.GenerateInvokeFunctionRawTransactionWithAmountToSend(wallet, assetToSend, amount, unspent, attributes, operation, args)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
