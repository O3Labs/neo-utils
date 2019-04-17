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
	wif := "L23uNfE2CQxMhBgEma3QvKFbPbg9Wu7goWf71QwC1A9FaU7kPCga"
	wallet, _ := neoutils.GenerateFromWIF(wif)

	refundValue := float64(1)
	//c074a05e9dcf0141cbe6b4b3475dd67baf4dcb60 cneo
	//74f2dc36a68fdc4682034178eb2220729231db76 cgas
	cgas, _ := smartcontract.NewScriptHash("c074a05e9dcf0141cbe6b4b3475dd67baf4dcb60")

	// unspent, _ := utxo("test", wallet.Address)

	unspent := smartcontract.Unspent{}
	unspent.Assets = map[smartcontract.NativeAsset]*smartcontract.Balance{}

	//any utxo from SGAS address that is not marked as refund.
	txid := "0xf9f6df80e4085b262c90aae474d4b168fc701b76c0f6bed04f531b52236a385e"
	balance := smartcontract.Balance{
		Amount: float64(0) / float64(100000000),
		UTXOs:  []smartcontract.UTXO{},
	}

	utxo := smartcontract.UTXO{
		Index: 0,
		TXID:  txid,
		Value: 5,
	}
	balance.UTXOs = append(balance.UTXOs, utxo)
	unspent.Assets[smartcontract.NEO] = &balance

	from := smartcontract.ParseNEOAddress(wallet.Address)
	args := []interface{}{from}

	//New invocation transaction struct and fill with all necessary data
	tx := smartcontract.NewInvocationTransaction()
	// tx.GAS = uint64(0)
	txData := smartcontract.NewScriptBuilder().GenerateContractInvocationData(cgas, "refund", args)
	log.Printf("script %x", txData)
	tx.Data = txData

	//basically sending GAS to myself
	amountToSend := refundValue
	assetToSend := smartcontract.NEO

	networkFee := smartcontract.NetworkFeeAmount(0)

	txInputs, err := smartcontract.NewScriptBuilder().GenerateTransactionInput(unspent, assetToSend, amountToSend, networkFee)
	if err != nil {
		return
	}
	tx.Inputs = txInputs

	//this is a MUST
	sender := smartcontract.ParseNEOAddress("AQbg4gk1Q6FaGCtfEKu2ETSMP6U25YDVR3")
	receiver := smartcontract.ParseNEOAddress("AQbg4gk1Q6FaGCtfEKu2ETSMP6U25YDVR3")

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

	log.Printf("before sign %x", tx.ToBytes())

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
		Invocation:   []byte{0x00, 0x00},
		Verification: nil,
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

	//put CGAS contract in validation script
	// b := neoutils.HexTobytes(`012fc56b6c766b00527ac46c766b51527ac4616168164e656f2e52756e74696d652e47657454726967676572009c6c766b52527ac46c766b52c36418046161682953797374656d2e457865637574696f6e456e67696e652e476574536372697074436f6e7461696e65726c766b53527ac46c766b53c36168194e656f2e5472616e73616374696f6e2e476574496e707574736c766b54527ac46c766b53c361681a4e656f2e5472616e73616374696f6e2e4765744f7574707574736c766b55527ac4616c766b54c36c766b59527ac4006c766b5a527ac4622b016c766b59c36c766b5ac3c36c766b5b527ac4616c766b5bc36168124e656f2e496e7075742e476574496e646578009c6c766b5c527ac46c766b5cc364de00616168164e656f2e53746f726167652e476574436f6e7465787406726566756e64617c6545176c766b5d527ac46c766b5dc36c766b5bc36168114e656f2e496e7075742e47657448617368617c6555176c766b5e527ac46c766b5ec3c000a06c766b5f527ac46c766b5fc3646f00616c766b54c3c051907c907c9e6310006c766b55c3c0519c009c620400516c766b60527ac46c766b60c3640f00006c766b0111527ac4620f076c766b55c300c36168184e656f2e4f75747075742e476574536372697074486173686c766b5ec39c6c766b0111527ac462dc0661616c766b5ac351936c766b5a527ac46c766b5ac36c766b59c3c09f63ccfe61682d53797374656d2e457865637574696f6e456e67696e652e476574457865637574696e67536372697074486173686c766b56527ac4006c766b57527ac4616c766b53c361681d4e656f2e5472616e73616374696f6e2e4765745265666572656e6365736c766b0112527ac4006c766b0113527ac462e9006c766b0112c36c766b0113c3c36c766b0114527ac4616c766b0114c36168154e656f2e4f75747075742e476574417373657449646120e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c609e6c766b0115527ac46c766b0115c3640f00006c766b0111527ac462d0056c766b0114c36168184e656f2e4f75747075742e476574536372697074486173686c766b56c39c6c766b0116527ac46c766b0116c3642c006c766b57c36c766b0114c36168134e656f2e4f75747075742e47657456616c7565936c766b57527ac4616c766b0113c351936c766b0113527ac46c766b0113c36c766b0112c3c09f630cff006c766b58527ac4616c766b55c36c766b0117527ac4006c766b0118527ac4628b006c766b0117c36c766b0118c3c36c766b0119527ac4616c766b0119c36168184e656f2e4f75747075742e476574536372697074486173686c766b56c39c6c766b011a527ac46c766b011ac3642c006c766b58c36c766b0119c36168134e656f2e4f75747075742e47657456616c7565936c766b58527ac4616c766b0118c351936c766b0118527ac46c766b0118c36c766b0117c3c09f636aff6c766b58c36c766b57c39c6c766b0111527ac4627c046168164e656f2e52756e74696d652e47657454726967676572609c6c766b011b527ac46c766b011bc3649e026161682b53797374656d2e457865637574696f6e456e67696e652e47657443616c6c696e67536372697074486173686c766b011c527ac46c766b00c30962616c616e63654f66876c766b011d527ac46c766b011dc36419006c766b51c300c36165f7036c766b0111527ac462e2036c766b00c308646563696d616c73876c766b011e527ac46c766b011ec3641200616570046c766b0111527ac462b3036c766b00c30f676574526566756e64546172676574876c766b011f527ac46c766b011fc36419006c766b51c300c361653b046c766b0111527ac46276036c766b00c3096765745478496e666f876c766b0120527ac46c766b0120c36419006c766b51c300c36165b1046c766b0111527ac4623f036c766b00c30a6d696e74546f6b656e73876c766b0121527ac46c766b0121c36412006165e6056c766b0111527ac4620e036c766b00c3046e616d65876c766b0122527ac46c766b0122c36412006165200a6c766b0111527ac462e3026c766b00c306726566756e64876c766b0123527ac46c766b0123c36419006c766b51c300c36165fc096c766b0111527ac462af026c766b00c30673796d626f6c876c766b0124527ac46c766b0124c36412006165f70e6c766b0111527ac46282026c766b00c312737570706f727465645374616e6461726473876c766b0125527ac46c766b0125c36412006165ca0e6c766b0111527ac46249026c766b00c30b746f74616c537570706c79876c766b0126527ac46c766b0126c36412006165bc0e6c766b0111527ac46217026c766b00c3087472616e73666572876c766b0127527ac46c766b0127c36441006c766b51c300c36c766b51c351c36c766b51c352c36c766b011cc361537951795572755172755279527954727552727565d60e6c766b0111527ac462b9016162a9016168164e656f2e52756e74696d652e47657454726967676572519c6c766b0128527ac46c766b0128c3647d016161682d53797374656d2e457865637574696f6e456e67696e652e476574457865637574696e67536372697074486173686c766b0129527ac461682953797374656d2e457865637574696f6e456e67696e652e476574536372697074436f6e7461696e65726c766b012a527ac4616c766b012ac361681a4e656f2e5472616e73616374696f6e2e4765744f7574707574736c766b012b527ac4006c766b012c527ac462bb006c766b012bc36c766b012cc3c36c766b012d527ac4616c766b012dc36168184e656f2e4f75747075742e476574536372697074486173686c766b0129c3907c907c9e6347006c766b012dc36168154e656f2e4f75747075742e476574417373657449646120e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c609e620400006c766b012e527ac46c766b012ec3640f00006c766b0111527ac4623d00616c766b012cc351936c766b012c527ac46c766b012cc36c766b012bc3c09f633aff516c766b0111527ac4620f00006c766b0111527ac46203006c766b0111c3616c756654c56b6c766b00527ac4616c766b00c3c001149c009c6c766b52527ac46c766b52c36439003254686520706172616d65746572206163636f756e742053484f554c442062652032302d62797465206164647265737365732e6175f06168164e656f2e53746f726167652e476574436f6e74657874056173736574617c652f0f6c766b51527ac46c766b51c36c766b00c3617c65530f6c766b53527ac46203006c766b53c3616c756600c56b58616c756654c56b6c766b00527ac4616c766b00c3c001209c009c6c766b52527ac46c766b52c3643d003654686520706172616d6574657220747849642053484f554c442062652033322d62797465207472616e73616374696f6e20686173682e6175f06168164e656f2e53746f726167652e476574436f6e7465787406726566756e64617c657a0e6c766b51527ac46c766b51c36c766b00c3617c659e0e6c766b53527ac46203006c766b53c3616c756656c56b6c766b00527ac4616c766b00c3c001209c009c6c766b53527ac46c766b53c3643d003654686520706172616d6574657220747849642053484f554c442062652033322d62797465207472616e73616374696f6e20686173682e6175f06168164e656f2e53746f726167652e476574436f6e74657874067478496e666f617c65cd0d6c766b51527ac46c766b51c36c766b00c3617c65f10d6c766b52527ac46c766b52c3c0009c6c766b54527ac46c766b54c3640e00006c766b55527ac4622c006c766b52c36168174e656f2e52756e74696d652e446573657269616c697a656c766b55527ac46203006c766b55c3616c756653c56b6c766b00527ac4616c766b00c361681a4e656f2e426c6f636b636861696e2e476574436f6e74726163746c766b51527ac46c766b51c36424006c766b51c36168164e656f2e436f6e74726163742e497350617961626c65620400516c766b52527ac46203006c766b52c3616c75660114c56b6161682953797374656d2e457865637574696f6e456e67696e652e476574536372697074436f6e7461696e65726c766b00527ac4006c766b51527ac46c766b00c361681d4e656f2e5472616e73616374696f6e2e4765745265666572656e6365736c766b52527ac4616c766b52c36c766b59527ac4006c766b5a527ac46210016c766b59c36c766b5ac3c36c766b5b527ac4616c766b5bc36168154e656f2e4f75747075742e476574417373657449646120e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c609c6c766b5c527ac46c766b5cc36434006c766b51c376632400756c766b5bc36168184e656f2e4f75747075742e476574536372697074486173686c766b51527ac46c766b5bc36168184e656f2e4f75747075742e4765745363726970744861736861682d53797374656d2e457865637574696f6e456e67696e652e476574457865637574696e67536372697074486173689c6c766b5d527ac46c766b5dc3640e00006c766b5e527ac462dd02616c766b5ac351936c766b5a527ac46c766b5ac36c766b59c3c09f63e7fe6c766b00c36168174e656f2e5472616e73616374696f6e2e476574486173686165dafc00a06c766b5f527ac46c766b5fc3640e00006c766b5e527ac46280026c766b00c361681a4e656f2e5472616e73616374696f6e2e4765744f7574707574736c766b53527ac4006c766b54527ac4616c766b53c36c766b60527ac4006c766b0111527ac46203016c766b60c36c766b0111c3c36c766b0112527ac4616c766b0112c36168184e656f2e4f75747075742e4765745363726970744861736861682d53797374656d2e457865637574696f6e456e67696e652e476574457865637574696e6753637269707448617368907c907c9e6347006c766b0112c36168154e656f2e4f75747075742e476574417373657449646120e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c609c620400006c766b0113527ac46c766b0113c3642e00616c766b54c36c766b0112c36168134e656f2e4f75747075742e47657456616c7565936c766b54527ac461616c766b0111c351936c766b0111527ac46c766b0111c36c766b60c3c09f63f3fe6168164e656f2e53746f726167652e476574436f6e7465787408636f6e7472616374617c658b096c766b55527ac46c766b55c30b746f74616c537570706c79617c65030a6c766b56527ac46c766b56c36c766b54c3936c766b56527ac46c766b55c30b746f74616c537570706c796c766b56c361527265290a616168164e656f2e53746f726167652e476574436f6e74657874056173736574617c6514096c766b57527ac46c766b57c36c766b51c3617c6538096c766b58527ac46c766b57c36c766b51c36c766b58c36c766b54c39361527265260a61006c766b51c36c766b54c36152726596046161006c766b51c36c766b54c3615272087472616e7366657254c168124e656f2e52756e74696d652e4e6f7469667961516c766b5e527ac46203006c766b5ec3616c756600c56b084e45503520474153616c75660111c56b6c766b00527ac4616c766b00c3c001149c009c6c766b59527ac46c766b59c36436002f54686520706172616d657465722066726f6d2053484f554c442062652032302d62797465206164647265737365732e6175f061682953797374656d2e457865637574696f6e456e67696e652e476574536372697074436f6e7461696e65726c766b51527ac46c766b51c361681a4e656f2e5472616e73616374696f6e2e4765744f75747075747300c36c766b52527ac46c766b52c36168154e656f2e4f75747075742e476574417373657449646120e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c609e6c766b5a527ac46c766b5ac3640e00006c766b5b527ac46228036c766b52c36168184e656f2e4f75747075742e4765745363726970744861736861682d53797374656d2e457865637574696f6e456e67696e652e476574457865637574696e67536372697074486173689e6c766b5c527ac46c766b5cc3640e00006c766b5b527ac462bd026168164e656f2e53746f726167652e476574436f6e7465787406726566756e64617c65d5066c766b53527ac46c766b53c36c766b51c36168174e656f2e5472616e73616374696f6e2e47657448617368617c65df06c000a06c766b5d527ac46c766b5dc3640e00006c766b5b527ac4624b026c766b00c36168184e656f2e52756e74696d652e436865636b5769746e657373009c6c766b5e527ac46c766b5ec3640e00006c766b5b527ac4620f026168164e656f2e53746f726167652e476574436f6e74657874056173736574617c6528066c766b54527ac46c766b54c36c766b00c3617c654c066c766b55527ac46c766b52c36168134e656f2e4f75747075742e47657456616c75656c766b56527ac46c766b55c36c766b56c39f6c766b5f527ac46c766b5fc3640e00006c766b5b527ac46287016c766b55c36c766b56c39c6c766b60527ac46c766b60c36416006c766b54c36c766b00c3617c653f0761621f006c766b54c36c766b00c36c766b55c36c766b56c39461527265c606616c766b53c36c766b51c36168174e656f2e5472616e73616374696f6e2e476574486173686c766b00c3615272654007616168164e656f2e53746f726167652e476574436f6e7465787408636f6e7472616374617c6524056c766b57527ac46c766b57c30b746f74616c537570706c79617c659c056c766b58527ac46c766b58c36c766b56c3946c766b58527ac46c766b57c30b746f74616c537570706c796c766b58c361527265c205616c766b00c3006c766b56c3615272658c0061616c766b00c3006c766b56c3615272087472616e7366657254c168124e656f2e52756e74696d652e4e6f7469667961616c766b51c36168174e656f2e5472616e73616374696f6e2e476574486173686c766b00c3617c06726566756e6453c168124e656f2e52756e74696d652e4e6f7469667961516c766b5b527ac46203006c766b5bc3616c756656c56b6c766b00527ac46c766b51527ac46c766b52527ac46161682953797374656d2e457865637574696f6e456e67696e652e476574536372697074436f6e7461696e65726168174e656f2e5472616e73616374696f6e2e476574486173686c766b53527ac46153c5766c766b00c3007cc4766c766b51c3517cc4766c766b52c3527cc46c766b54527ac46168164e656f2e53746f726167652e476574436f6e74657874067478496e666f617c6587036c766b55527ac46c766b55c36c766b53c36c766b54c36168154e656f2e52756e74696d652e53657269616c697a6561527265470561616c756600c56b0443474153616c756600c56b1c7b224e45502d35222c20224e45502d37222c20224e45502d3130227d616c756652c56b616168164e656f2e53746f726167652e476574436f6e7465787408636f6e7472616374617c65f3026c766b00527ac46c766b00c30b746f74616c537570706c79617c656b036c766b51527ac46203006c766b51c3616c756653c56b6c766b00527ac46c766b51527ac46c766b52527ac451616c75665fc56b6c766b00527ac46c766b51527ac46c766b52527ac46c766b53527ac4616c766b00c3c00114907c907c9e6311006c766b51c3c001149c009c620400516c766b57527ac46c766b57c3643e003754686520706172616d65746572732066726f6d20616e6420746f2053484f554c442062652032302d62797465206164647265737365732e6175f06c766b52c300a16c766b58527ac46c766b58c36433002c54686520706172616d6574657220616d6f756e74204d5553542062652067726561746572207468616e20302e6175f06c766b51c3616575f4009c6c766b59527ac46c766b59c3640e00006c766b5a527ac462a9016c766b00c36168184e656f2e52756e74696d652e436865636b5769746e6573736311006c766b00c36c766b53c39e620400006c766b5b527ac46c766b5bc3640e00006c766b5a527ac4625d016168164e656f2e53746f726167652e476574436f6e74657874056173736574617c6542016c766b54527ac46c766b54c36c766b00c3617c6566016c766b55527ac46c766b55c36c766b52c39f6c766b5c527ac46c766b5cc3640e00006c766b5a527ac462f7006c766b00c36c766b51c39c6c766b5d527ac46c766b5dc3640e00516c766b5a527ac462d2006c766b55c36c766b52c39c6c766b5e527ac46c766b5ec36416006c766b54c36c766b00c3617c65560261621f006c766b54c36c766b00c36c766b55c36c766b52c39461527265dd01616c766b54c36c766b51c3617c65bd006c766b56527ac46c766b54c36c766b51c36c766b56c36c766b52c39361527265ab01616c766b00c36c766b51c36c766b52c36152726517fc61616c766b00c36c766b51c36c766b52c3615272087472616e7366657254c168124e656f2e52756e74696d652e4e6f7469667961516c766b5a527ac46203006c766b5ac3616c756653c56b6c766b00527ac46c766b51527ac4616152c5766c766b00c3007cc4766c766b51c3517cc46c766b52527ac46203006c766b52c3616c756654c56b6c766b00527ac46c766b51527ac4616c766b00c351c301007e6c766b51c37e6c766b52527ac46c766b00c300c36c766b52c3617c680f4e656f2e53746f726167652e4765746c766b53527ac46203006c766b53c3616c756654c56b6c766b00527ac46c766b51527ac4616c766b00c351c301007e6c766b51c37e6c766b52527ac46c766b00c300c36c766b52c3617c680f4e656f2e53746f726167652e4765746c766b53527ac46203006c766b53c3616c756654c56b6c766b00527ac46c766b51527ac46c766b52527ac4616c766b00c351c301007e6c766b51c37e6c766b53527ac46c766b00c300c36c766b53c36c766b52c3615272680f4e656f2e53746f726167652e50757461616c756654c56b6c766b00527ac46c766b51527ac46c766b52527ac4616c766b00c351c301007e6c766b51c37e6c766b53527ac46c766b00c300c36c766b53c36c766b52c3615272680f4e656f2e53746f726167652e50757461616c756653c56b6c766b00527ac46c766b51527ac4616c766b00c351c301007e6c766b51c37e6c766b52527ac46c766b00c300c36c766b52c3617c68124e656f2e53746f726167652e44656c65746561616c756654c56b6c766b00527ac46c766b51527ac46c766b52527ac4616c766b00c351c301007e6c766b51c37e6c766b53527ac46c766b00c300c36c766b53c36c766b52c3615272680f4e656f2e53746f726167652e50757461616c7566`)

	// cgasVerificationScript := smartcontract.TransactionValidationScript{
	// 	Invocation:   []byte{0x00, 0x00},
	// 	Verification: nil,
	// }

	// scripts = append([]interface{}{cgasVerificationScript}, scripts...)
	// scripts = append(scripts, cgasVerificationScript)

	txScripts := smartcontract.NewScriptBuilder().GenerateVerificationScripts(scripts)
	//assign scripts to the tx
	tx.Script = txScripts

	log.Printf("txid = %v", tx.ToTXID())
	log.Printf("endPayload = %x", tx.ToBytes())
}

func TestInvokeGetRefundTarget(t *testing.T) {
	cgas, _ := smartcontract.NewScriptHash("74f2dc36a68fdc4682034178eb2220729231db76")
	args := []interface{}{"cb2b434b73dfc3abf4bbb779ce21d0b58ff847c49c285f4ce3641fa6299f10fe"}
	s := smartcontract.NewScriptBuilder()

	s.GenerateContractInvocationScript(cgas, "getTxInfo", args)
	log.Printf("%x", s.ToBytes())
}

func TestRefund2ndStep(t *testing.T) {
	wif := "L23uNfE2CQxMhBgEma3QvKFbPbg9Wu7goWf71QwC1A9FaU7kPCga"
	wallet, _ := neoutils.GenerateFromWIF(wif)
	refundValue := float64(1)

	cgas, _ := smartcontract.NewScriptHash("c074a05e9dcf0141cbe6b4b3475dd67baf4dcb60")

	unspent := smartcontract.Unspent{}
	unspent.Assets = map[smartcontract.NativeAsset]*smartcontract.Balance{}

	//any utxo from SGAS address that is not marked as refund.
	txid := "442f6e349a8a5289aa77cc8cc6b422c67d9563c13770238083504fc0c550bfac"
	gasBalance := smartcontract.Balance{
		Amount: float64(0) / float64(100000000),
		UTXOs:  []smartcontract.UTXO{},
	}

	utxo := smartcontract.UTXO{
		Index: 0,
		TXID:  txid,
		Value: 1,
	}
	gasBalance.UTXOs = append(gasBalance.UTXOs, utxo)
	unspent.Assets[smartcontract.NEO] = &gasBalance

	from := smartcontract.ParseNEOAddress(wallet.Address)
	args := []interface{}{from}

	//New invocation transaction struct and fill with all necessary data
	tx := smartcontract.NewInvocationTransactionPayable()
	tx.GAS = uint64(0)
	txData := smartcontract.NewScriptBuilder().GenerateContractInvocationData(cgas, "refund", args)
	tx.Data = txData

	//basically sending GAS to myself
	amountToSend := refundValue
	assetToSend := smartcontract.NEO

	networkFee := smartcontract.NetworkFeeAmount(0)

	txInputs, err := smartcontract.NewScriptBuilder().GenerateTransactionInput(unspent, assetToSend, amountToSend, networkFee)
	if err != nil {
		return
	}
	tx.Inputs = txInputs
	log.Printf("input %x", txInputs)

	sender := smartcontract.ParseNEOAddress("AQbg4gk1Q6FaGCtfEKu2ETSMP6U25YDVR3")
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
	log.Printf("before sign %x", tx.ToBytes())
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

	// this empty verification script is needed in order to make it triggers Verification part
	emptyVerificationScript := smartcontract.TransactionValidationScript{
		Invocation:   []byte{0x00, 0x00},
		Verification: nil,
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
