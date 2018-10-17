package neoutils_test

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func TestInvalidSmartContractStruct(t *testing.T) {
	sc := neoutils.UseSmartContract("ce575ae1bb6153330d2")
	if sc != nil {
		t.Fail()
		return
	}
	log.Printf("%v", sc)
}

func TestUseSmartContractWithEmptyScripthash(t *testing.T) {
	sc := neoutils.UseSmartContract("")
	if sc != nil {
		t.Fail()
		return
	}
	log.Printf("%v", sc)
}

func TestSmartContractStruct(t *testing.T) {

	sc := neoutils.UseSmartContract("ce575ae1bb6153330d20c560acb434dc5755241b")
	if sc == nil {
		t.Fail()
		return
	}
	log.Printf("%v", sc)
}

func UTXODataForSmartContract() smartcontract.Unspent {

	gasTX1 := smartcontract.UTXO{
		Index: 0,
		TXID:  "307d756074d9ee11220ccebf003bedb99f9b1a54e194a25e6ea5df1a7b2de84b",
		Value: float64(713399700000) / float64(100000000),
	}

	gasBalance := smartcontract.Balance{
		Amount: float64(713399700000) / float64(100000000),
		UTXOs:  []smartcontract.UTXO{gasTX1},
	}

	neoTX1 := smartcontract.UTXO{
		Index: 0,
		TXID:  "e8b8bf4f98490368fc1caa86f8646e7383bb52751ffc3a1a7e296d715c4382ed",
		Value: float64(10000000000000000) / float64(100000000),
	}

	neoBalance := smartcontract.Balance{
		Amount: float64(10000000000000000) / float64(100000000),
		UTXOs:  []smartcontract.UTXO{neoTX1},
	}

	unspent := smartcontract.Unspent{
		Assets: map[smartcontract.NativeAsset]*smartcontract.Balance{},
	}
	unspent.Assets[smartcontract.NEO] = &neoBalance
	unspent.Assets[smartcontract.GAS] = &gasBalance
	return unspent
}

func TestInvokeFunctionRawTransaction(t *testing.T) {
	var validSmartContract = neoutils.UseSmartContract("b7c1f850a025e34455e7e98c588c784385077fb1")

	wif := ""
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
	}

	to := smartcontract.ParseNEOAddress("AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn")
	if to == nil {
		//invalid neo address
		t.Fail()
		return
	}
	args := []interface{}{to, 1}
	unspent := UTXODataForSmartContract()

	transactionID := "thisisauniquetoken_from_stripe"
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	attributes[smartcontract.Remark1] = []byte(transactionID)
	tx, err := validSmartContract.GenerateInvokeFunctionRawTransaction(*privateNetwallet, unspent, attributes, "mintTokensTo", args)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("%x", tx)
}

func TestGenerateInvokeTransferNEP5Token(t *testing.T) {
	var validSmartContract = neoutils.UseSmartContract("b7c1f850a025e34455e7e98c588c784385077fb1")

	wif := ""
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
	}

	from := smartcontract.ParseNEOAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
	if from == nil {
		//invalid neo address
		t.Fail()
		return
	}

	to := smartcontract.ParseNEOAddress("AM8pnu1yK7ViMt7Sw2nPpbtPQXTwjjkykn")
	if to == nil {
		//invalid neo address
		t.Fail()
		return
	}
	numberOfTokens := 1
	args := []interface{}{from, to, numberOfTokens}
	unspent := UTXODataForSmartContract()

	remark := "this is a remark data in attribute"
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	attributes[smartcontract.Remark1] = []byte(remark)
	tx, err := validSmartContract.GenerateInvokeFunctionRawTransaction(*privateNetwallet, unspent, attributes, "transfer", args)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("%x", tx)
}

func TestCallDeployFunction(t *testing.T) {

	encryptedKey := ""
	passphrase := ""
	wif, err := neoutils.NEP2Decrypt(encryptedKey, passphrase)

	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
	}

	unspent := smartcontract.Unspent{}

	sc := neoutils.UseSmartContract("323571cfc42a40d48d64832a7da594039fbac76a")
	args := []interface{}{}
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	addressScriptHash := neoutils.NEOAddressToScriptHashWithEndian(privateNetwallet.Address, binary.LittleEndian)
	b, _ := hex.DecodeString(addressScriptHash)
	attributes[smartcontract.Script] = []byte(b)
	attributes[smartcontract.Remark1] = []byte(fmt.Sprintf("O3TXSCC%v", time.Now().Unix()))

	tx, err := sc.GenerateInvokeFunctionRawTransaction(*privateNetwallet, unspent, attributes, "deploy", args)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("%x", tx)
}

func TestRefund1stCGAS(t *testing.T) {
	wif := ""
	wallet, _ := neoutils.GenerateFromWIF(wif)

	refundValue := float64(1)

	cgas, _ := smartcontract.NewScriptHash("9121e89e8a0849857262d67c8408601b5e8e0524")

	// unspent, _ := utxo("test", wallet.Address)

	unspent := smartcontract.Unspent{}
	unspent.Assets = map[smartcontract.NativeAsset]*smartcontract.Balance{}

	//any utxo from SGAS address that is not marked as refund.
	txid := "0x472eadfd5fc2b726d07b78a83887a5f9ea00eafe1bc7dcc899dcbed21c9c99af"
	gasBalance := smartcontract.Balance{
		Amount: float64(0) / float64(100000000),
		UTXOs:  []smartcontract.UTXO{},
	}

	gasTX1 := smartcontract.UTXO{
		Index: 0,
		TXID:  txid,
		Value: 1,
	}
	gasBalance.UTXOs = append(gasBalance.UTXOs, gasTX1)
	unspent.Assets[smartcontract.GAS] = &gasBalance

	from := smartcontract.ParseNEOAddress(wallet.Address)
	args := []interface{}{from}

	//New invocation transaction struct and fill with all necessary data
	tx := smartcontract.NewInvocationTransaction()

	txData := smartcontract.NewScriptBuilder().GenerateContractInvocationData(cgas, "refund", args)
	tx.Data = txData

	//basically sending GAS to myself
	amountToSend := refundValue
	assetToSend := smartcontract.GAS

	networkFee := smartcontract.NetworkFeeAmount(0)

	txInputs, err := smartcontract.NewScriptBuilder().GenerateTransactionInput(unspent, assetToSend, amountToSend, networkFee)
	if err != nil {
		return
	}
	tx.Inputs = txInputs
	log.Printf("input %x", txInputs)
	//this is a MUST
	sender := smartcontract.ParseNEOAddress("AK4LdT5ZXR9DQZjfk5X6Xy79mE8ad8jKAW")
	receiver := smartcontract.ParseNEOAddress("AK4LdT5ZXR9DQZjfk5X6Xy79mE8ad8jKAW")
	log.Printf("receiver %v", receiver.ToString())
	txOutputs, err := smartcontract.NewScriptBuilder().GenerateTransactionOutput(sender, receiver, unspent, assetToSend, amountToSend, networkFee)
	if err != nil {
		return
	}

	tx.Outputs = txOutputs

	attributes := map[smartcontract.TransactionAttribute][]byte{}
	// attributes[smartcontract.Remark1] = []byte(remark)
	attributes[smartcontract.Script] = neoutils.HexTobytes(neoutils.NEOAddressToScriptHashWithEndian(wallet.Address, binary.LittleEndian))

	//generate transaction outputs
	txAttributes, err := smartcontract.NewScriptBuilder().GenerateTransactionAttributes(attributes)
	if err != nil {
		return
	}
	//transaction attributes
	tx.Attributes = txAttributes

	//begin signing process and invocation script
	privateKeyInHex := neoutils.BytesToHex(wallet.PrivateKey)
	signedData, err := neoutils.Sign(tx.ToBytes(), privateKeyInHex)
	if err != nil {
		return
	}

	signature := smartcontract.TransactionSignature{
		SignedData: signedData,
		PublicKey:  wallet.PublicKey,
	}

	scripts := []interface{}{signature}

	//this empty verification script is needed in order to make it triggers Verification part
	//and use Script field in Transaction attribute
	emptyVerificationScript := smartcontract.TransactionValidationScript{
		StackScript:  []byte{0x00, 0x00},
		RedeemScript: nil,
	}

	//basically we need to sort in descending order for address and script hash
	scriptHashInt := neoutils.ConvertByteArrayToBigInt(fmt.Sprintf("%x", cgas))
	addressInt := neoutils.ConvertByteArrayToBigInt(fmt.Sprintf("%x", wallet.HashedSignature))
	//https://godoc.org/math/big#Int.Cmp
	//if scripthash int is grether than address int
	if scriptHashInt.Cmp(addressInt) == 1 {
		scripts = append(scripts, emptyVerificationScript)
	} else {
		scripts = append([]interface{}{emptyVerificationScript}, scripts...)
	}

	txScripts := smartcontract.NewScriptBuilder().GenerateVerificationScripts(scripts)
	//assign scripts to the tx
	tx.Script = txScripts

	log.Printf("txid = %v", tx.ToTXID())
	log.Printf("endPayload = %x", tx.ToBytes())
}

func TestRefund2ndStep(t *testing.T) {
	wif := ""
	wallet, _ := neoutils.GenerateFromWIF(wif)

	refundValue := float64(1)

	cgas, _ := smartcontract.NewScriptHash("9121e89e8a0849857262d67c8408601b5e8e0524")

	unspent := smartcontract.Unspent{}
	unspent.Assets = map[smartcontract.NativeAsset]*smartcontract.Balance{}
	txid := "0x9296fef6b14f85eb29155639ab3cf46edd6fcc529177b7259bfeb2a932278238"
	gasBalance := smartcontract.Balance{
		Amount: float64(0) / float64(100000000),
		UTXOs:  []smartcontract.UTXO{},
	}

	gasTX1 := smartcontract.UTXO{
		Index: 0,
		TXID:  txid,
		Value: 1,
	}
	gasBalance.UTXOs = append(gasBalance.UTXOs, gasTX1)
	unspent.Assets[smartcontract.GAS] = &gasBalance

	from := smartcontract.ParseNEOAddress(wallet.Address)
	args := []interface{}{from}

	//New invocation transaction struct and fill with all necessary data
	tx := smartcontract.NewInvocationTransaction()

	txData := smartcontract.NewScriptBuilder().GenerateContractInvocationData(cgas, "refund", args)
	tx.Data = txData

	//basically sending GAS to myself
	amountToSend := refundValue
	assetToSend := smartcontract.GAS

	networkFee := smartcontract.NetworkFeeAmount(0)

	txInputs, err := smartcontract.NewScriptBuilder().GenerateTransactionInput(unspent, assetToSend, amountToSend, networkFee)
	if err != nil {
		return
	}
	tx.Inputs = txInputs
	log.Printf("input %x", txInputs)

	sender := smartcontract.ParseNEOAddress("AK4LdT5ZXR9DQZjfk5X6Xy79mE8ad8jKAW")
	receiver := smartcontract.ParseNEOAddress(wallet.Address)
	log.Printf("receiver %v", receiver.ToString())
	txOutputs, err := smartcontract.NewScriptBuilder().GenerateTransactionOutput(sender, receiver, unspent, assetToSend, amountToSend, networkFee)
	if err != nil {
		return
	}

	tx.Outputs = txOutputs

	attributes := map[smartcontract.TransactionAttribute][]byte{}
	// attributes[smartcontract.Remark1] = []byte(remark)
	attributes[smartcontract.Script] = neoutils.HexTobytes(neoutils.NEOAddressToScriptHashWithEndian(wallet.Address, binary.LittleEndian))

	//generate transaction outputs
	txAttributes, err := smartcontract.NewScriptBuilder().GenerateTransactionAttributes(attributes)
	if err != nil {
		return
	}
	//transaction attributes
	tx.Attributes = txAttributes

	//begin signing process and invocation script
	privateKeyInHex := neoutils.BytesToHex(wallet.PrivateKey)
	signedData, err := neoutils.Sign(tx.ToBytes(), privateKeyInHex)
	if err != nil {
		return
	}

	signature := smartcontract.TransactionSignature{
		SignedData: signedData,
		PublicKey:  wallet.PublicKey,
	}

	scripts := []interface{}{signature}

	//this empty verification script is needed in order to make it triggers Verification part
	emptyVerificationScript := smartcontract.TransactionValidationScript{
		StackScript:  []byte{0x00, 0x00},
		RedeemScript: nil,
	}

	//basically we need to sort in descending order for address and script hash
	scriptHashInt := neoutils.ConvertByteArrayToBigInt(fmt.Sprintf("%x", cgas))
	addressInt := neoutils.ConvertByteArrayToBigInt(fmt.Sprintf("%x", wallet.HashedSignature))
	//https://godoc.org/math/big#Int.Cmp
	//if scripthash int is grether than address int
	if scriptHashInt.Cmp(addressInt) == 1 {
		scripts = append(scripts, emptyVerificationScript)
	} else {
		scripts = append([]interface{}{emptyVerificationScript}, scripts...)
	}

	txScripts := smartcontract.NewScriptBuilder().GenerateVerificationScripts(scripts)
	//assign scripts to the tx
	tx.Script = txScripts

	log.Printf("txid = %v", tx.ToTXID())
	log.Printf("endPayload = %x", tx.ToBytes())

}
