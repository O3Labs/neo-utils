package smartcontract

import (
	"encoding/hex"
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
)

func TestInvocationTransactionToBytes(t *testing.T) {
	tx := NewInvocationTransaction()
	log.Printf("%x", tx.ToBytes())
}

func mintTokensToData() []byte {

	s := NewScriptBuilder()
	scriptHash, err := NewScriptHash("ce575ae1bb6153330d20c560acb434dc5755241b")
	if err != nil {
		log.Printf("err = %v", err)
		return nil
	}
	to := ParseNEOAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
	if to == nil {
		//invalid neo address
		log.Println("invalid neo address")
		return nil
	}
	args := []interface{}{to, 1000}
	b := s.GenerateContractInvocationData(scriptHash, "mintTokensTo", args)
	return b
}

func UTXODataForSmartContract() Unspent {

	gasTX1 := UTXO{
		Index: 0,
		TXID:  "880081a69debf8f94187f83e91e67af5d53615bdd2383d3611b7a8eead049ea4",
		Value: float64(1) / float64(100000000),
	}

	gasBalance := Balance{
		Amount: float64(715799899999) / float64(100000000),
		UTXOs:  []UTXO{gasTX1},
	}

	neoTX1 := UTXO{
		Index: 0,
		TXID:  "e8b8bf4f98490368fc1caa86f8646e7383bb52751ffc3a1a7e296d715c4382ed",
		Value: float64(10000000000000000) / float64(100000000),
	}

	neoBalance := Balance{
		Amount: float64(10000000000000000) / float64(100000000),
		UTXOs:  []UTXO{neoTX1},
	}

	unspent := Unspent{
		Assets: map[NativeAsset]*Balance{},
	}
	unspent.Assets[neo] = &neoBalance
	unspent.Assets[gas] = &gasBalance
	return unspent
}

func inputs() []byte {
	s := NewScriptBuilder()
	assetToSend := gas
	amount := float64(0.00000001)
	unspent := UTXODataForSmartContract()
	b, err := s.GenerateTransactionInput(unspent, assetToSend, amount)
	if err != nil {
		log.Printf("err = %v", err)

		return nil
	}
	return b
}

func outputs() []byte {
	s := NewScriptBuilder()
	assetToSend := gas
	amountToSend := float64(0.00000001)
	unspent := UTXODataForSmartContract()
	sender := ParseNEOAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
	receiver := ParseNEOAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
	b, err := s.GenerateTransactionOutput(sender, receiver, unspent, assetToSend, amountToSend)
	if err != nil {
		log.Printf("%v", err)
		return nil
	}
	return b
}

func attributes() []byte {
	s := NewScriptBuilder()
	// return s.EmptyTransactionAttributes()
	attributes := map[TransactionAttribute][]byte{}
	attributes[Remark] = []byte("test")
	b, err := s.GenerateTransactionAttributes(attributes)
	if err != nil {
		return nil
	}
	return b
}

func TestMintTokensToInvocation(t *testing.T) {
	scriptHash, _ := NewScriptHash("ce575ae1bb6153330d20c560acb434dc5755241b")

	tx := NewInvocationTransaction()
	tx.Data = mintTokensToData()
	tx.Inputs = inputs()
	tx.Outputs = outputs()
	tx.Attributes = attributes()

	wif := "KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr"
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
	}
	privateKeyInHex := hex.EncodeToString(privateNetwallet.PrivateKey)

	signedData, err := neoutils.Sign(tx.ToBytes(), privateKeyInHex)
	if err != nil {
		log.Printf("err signing %v", err)
		t.Fail()
	}
	s := NewScriptBuilder()
	signature := TransactionSignature{
		SignedData: signedData,
		PublicKey:  privateNetwallet.PublicKey,
	}
	signatures := []TransactionSignature{signature}
	scripts := s.GenerateInvocationAndVerificationScriptWithSignatures(signatures)

	tx.Script = scripts

	endPayload := []byte{}
	endPayload = append(endPayload, tx.ToBytes()...)
	endPayload = append(endPayload, scriptHash.ToBigEndian()...)

	log.Printf("%x", endPayload)
}
