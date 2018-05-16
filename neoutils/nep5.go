package neoutils

import (
	"fmt"
	"log"
	"strings"

	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

type NEP5Interface interface {
	MintTokensRawTransaction(wallet Wallet, assetToSend smartcontract.NativeAsset, amount float64, unspent smartcontract.Unspent, remark string) ([]byte, string, error)
	TransferNEP5RawTransaction(wallet Wallet, toAddress smartcontract.NEOAddress, amount float64, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte) ([]byte, string, error)
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

func (n *NEP5) TransferNEP5RawTransaction(wallet Wallet, toAddress smartcontract.NEOAddress, amount float64, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte) ([]byte, string, error) {

	from := smartcontract.ParseNEOAddress(wallet.Address)
	if from == nil {
		return nil, "", fmt.Errorf("Invalid from address")
	}

	to := smartcontract.ParseNEOAddress(toAddress.ToString())
	if to == nil {
		return nil, "", fmt.Errorf("Invalid from address")
	}
	numberOfTokens := amount
	args := []interface{}{from, to, numberOfTokens}

	//New invocation transaction struct and fill with all necessary data
	tx := smartcontract.NewInvocationTransaction()
	txData := smartcontract.NewScriptBuilder().GenerateContractInvocationData(n.ScriptHash, "transfer", args)
	log.Printf("txData = %x", txData)
	tx.Data = txData

	//for smart contract invocation we send the minimum amount of gas to it
	//0.00000001 gas
	amountToSend := float64(0.00000001)
	assetToSend := smartcontract.GAS

	//generate transaction inputs
	txInputs, err := smartcontract.NewScriptBuilder().GenerateTransactionInput(unspent, assetToSend, amountToSend, n.NetworkFeeAmount)
	if err != nil {
		return nil, "", err
	}
	log.Printf("txInputs = %x", txInputs)
	//transaction inputs
	tx.Inputs = txInputs

	//generate transaction outputs
	txAttributes, err := smartcontract.NewScriptBuilder().GenerateTransactionAttributes(attributes)
	if err != nil {
		return nil, "", err
	}
	log.Printf("txAttributes = %x", txAttributes)
	//transaction attributes
	tx.Attributes = txAttributes

	//send GAS to the same account
	sender := smartcontract.ParseNEOAddress(wallet.Address)
	receiver := smartcontract.ParseNEOAddress(wallet.Address)
	txOutputs, err := smartcontract.NewScriptBuilder().GenerateTransactionOutput(sender, receiver, unspent, assetToSend, amountToSend, n.NetworkFeeAmount)
	if err != nil {
		return nil, "", err
	}
	log.Printf("output = %x", txOutputs)
	tx.Outputs = txOutputs

	//begin signing process and invocation script
	privateKeyInHex := bytesToHex(wallet.PrivateKey)

	signedData, err := Sign(tx.ToBytes(), privateKeyInHex)
	if err != nil {
		return nil, "", err
	}

	needVerification := false
	if amountToSend == 0 {
		needVerification = true
	}

	signature := smartcontract.TransactionSignature{
		SignedData: signedData,
		PublicKey:  wallet.PublicKey,
	}

	scripts := []interface{}{signature}
	if needVerification == true {
		//this empty verification script is needed in order to make it triggers Verification part
		emptyVerificationScript := smartcontract.TransactionValidationScript{
			StackScript:  []byte{0x00, 0x00},
			RedeemScript: nil,
		}

		//this logic is still unknown to me
		//I need to check with the one who figured it out
		//https://github.com/CityOfZion/neon-js/blob/a9dfaefec870bfd05f3a8a0e5bc90a635fb6c5b9/src/api/core.js#L308

		//basically we need to sort in descending order for address and script hash
		scriptHashInt := ConvertByteArrayToBigInt(fmt.Sprintf("%x", n.ScriptHash))
		addressInt := ConvertByteArrayToBigInt(fmt.Sprintf("%x", wallet.HashedSignature))
		//https://godoc.org/math/big#Int.Cmp
		//if scripthash int is grether than address int
		if scriptHashInt.Cmp(addressInt) == 1 {
			scripts = append(scripts, emptyVerificationScript)
		} else {
			scripts = append([]interface{}{emptyVerificationScript}, scripts...)
		}
	}
	txScripts := smartcontract.NewScriptBuilder().GenerateVerificationScripts(scripts)
	//assign scripts to the tx
	log.Printf("txScripts = %x", txScripts)
	tx.Script = txScripts
	//end signing process

	//concat data
	endPayload := []byte{}
	endPayload = append(endPayload, tx.ToBytes()...)
	endPayload = append(endPayload, n.ScriptHash.ToBigEndian()...)

	//get tx id
	txID := tx.ToTXID()
	return endPayload, txID, nil
}

func (n *NEP5) MintTokensRawTransaction(wallet Wallet, assetToSend smartcontract.NativeAsset, amount float64, unspent smartcontract.Unspent, remark string) ([]byte, string, error) {

	needVerification := true
	operation := "mintTokens"
	args := []interface{}{}
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	attributes[smartcontract.Remark1] = []byte(remark)
	if needVerification == true {
		//add this to make it run VerificationTrigger
		attributes[smartcontract.Script] = n.ScriptHash
	}

	//New invocation transaction struct and fill with all necessary data
	tx := smartcontract.NewInvocationTransaction()
	txData := smartcontract.NewScriptBuilder().GenerateContractInvocationData(n.ScriptHash, operation, args)
	tx.Data = txData

	amountToSend := amount

	//generate transaction inputs
	txInputs, err := smartcontract.NewScriptBuilder().GenerateTransactionInput(unspent, assetToSend, amountToSend, n.NetworkFeeAmount)
	if err != nil {
		return nil, "", err
	}
	//transaction inputs
	tx.Inputs = txInputs

	//generate transaction outputs
	txAttributes, err := smartcontract.NewScriptBuilder().GenerateTransactionAttributes(attributes)
	if err != nil {
		return nil, "", err
	}
	//transaction attributes
	tx.Attributes = txAttributes

	//sender is a wallet address
	sender := smartcontract.ParseNEOAddress(wallet.Address)
	//when invoke the smart contract with amount of asset to send
	//we simply set the receiver to be the smart contract address
	receiver := smartcontract.NEOAddressFromScriptHash(n.ScriptHash.ToBigEndian())
	txOutputs, err := smartcontract.NewScriptBuilder().GenerateTransactionOutput(sender, receiver, unspent, assetToSend, amountToSend, n.NetworkFeeAmount)
	if err != nil {
		return nil, "", err
	}

	tx.Outputs = txOutputs

	//begin signing process and invocation script
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
	if needVerification == true {
		//this empty verification script is needed in order to make it triggers Verification part
		emptyVerificationScript := smartcontract.TransactionValidationScript{
			StackScript:  []byte{0x00, 0x00},
			RedeemScript: nil,
		}

		//this logic is still unknown to me
		//I need to check with the one who figured it out
		//https://github.com/CityOfZion/neon-js/blob/a9dfaefec870bfd05f3a8a0e5bc90a635fb6c5b9/src/api/core.js#L308

		//basically we need to sort in descending order for address and script hash
		scriptHashInt := ConvertByteArrayToBigInt(fmt.Sprintf("%x", n.ScriptHash))
		addressInt := ConvertByteArrayToBigInt(fmt.Sprintf("%x", wallet.HashedSignature))
		//https://godoc.org/math/big#Int.Cmp
		//if scripthash int is grether than address int
		if scriptHashInt.Cmp(addressInt) == 1 {
			scripts = append(scripts, emptyVerificationScript)
		} else {
			scripts = append([]interface{}{emptyVerificationScript}, scripts...)
		}
	}
	txScripts := smartcontract.NewScriptBuilder().GenerateVerificationScripts(scripts)
	//assign scripts to the tx
	tx.Script = txScripts

	//concat data
	endPayload := []byte{}
	endPayload = append(endPayload, tx.ToBytes()...)
	endPayload = append(endPayload, n.ScriptHash.ToBigEndian()...)

	//get tx id
	txID := tx.ToTXID()

	return endPayload, txID, nil
}
