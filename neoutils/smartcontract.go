package neoutils

import (
	"log"
	"strings"

	"github.com/o3labs/neo-utils/neoutils/smartcontract"
	// "github.com/o3labs/neo-utils/neoutils/smartcontract"
)

type SmartContractInterface interface {
	GenerateInvokeFunctionRawTransaction(wallet Wallet, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte, operation string, args []interface{}) ([]byte, error)
}

type SmartContract struct {
	ScriptHash smartcontract.ScriptHash
}

func UseSmartContract(scriptHashHex string) SmartContractInterface {
	if len(strings.TrimSpace(scriptHashHex)) == 0 {
		return nil
	}
	scripthash, err := smartcontract.NewScriptHash(scriptHashHex)
	if err != nil {
		return nil
	}

	return &SmartContract{ScriptHash: scripthash}
}

func (s *SmartContract) GenerateInvokeFunctionRawTransaction(wallet Wallet, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte, operation string, args []interface{}) ([]byte, error) {

	//New invocation transaction struct and fill with all necessary data
	tx := smartcontract.NewInvocationTransaction()
	txData := smartcontract.NewScriptBuilder().GenerateContractInvocationData(s.ScriptHash, operation, args)
	tx.Data = txData

	//for smart contract invocation we send the minimum amount of gas to it
	//0.00000001 gas
	amountToSend := float64(0.00000001)
	assetToSend := smartcontract.GAS

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
	receiver := smartcontract.ParseNEOAddress(wallet.Address)
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
	endPayload = append(endPayload, s.ScriptHash.ToBigEndian()...)

	return endPayload, nil
}
