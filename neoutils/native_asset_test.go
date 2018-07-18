package neoutils_test

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/o3"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func utxoFromO3Platform(network string, address string) (smartcontract.Unspent, error) {

	unspent := smartcontract.Unspent{
		Assets: map[smartcontract.NativeAsset]*smartcontract.Balance{},
	}

	client := o3.DefaultO3APIClient()
	if network == "test" {
		client = o3.APIClientWithNEOTestnet()
	} else if network == "private" {
		client = o3.APIClientWithNEOPrivateNet()
	}

	response := client.GetNEOUTXO(address)
	if response.Code != 200 {
		return unspent, fmt.Errorf("Error cannot get utxo")
	}

	gasBalance := smartcontract.Balance{
		Amount: float64(0),
		UTXOs:  []smartcontract.UTXO{},
	}

	neoBalance := smartcontract.Balance{
		Amount: float64(0),
		UTXOs:  []smartcontract.UTXO{},
	}

	for _, v := range response.Result.Data {
		if strings.Contains(v.Asset, string(smartcontract.GAS)) {
			value, err := strconv.ParseFloat(v.Value, 64)
			if err != nil {
				continue
			}
			gasTX1 := smartcontract.UTXO{
				Index: v.Index,
				TXID:  v.Txid,
				Value: value,
			}
			gasBalance.UTXOs = append(gasBalance.UTXOs, gasTX1)
		}

		if strings.Contains(v.Asset, string(smartcontract.NEO)) {
			value, err := strconv.ParseFloat(v.Value, 64)
			if err != nil {
				continue
			}
			tx := smartcontract.UTXO{
				Index: v.Index,
				TXID:  v.Txid,
				Value: value,
			}
			neoBalance.UTXOs = append(neoBalance.UTXOs, tx)
		}
	}

	unspent.Assets[smartcontract.GAS] = &gasBalance
	unspent.Assets[smartcontract.NEO] = &neoBalance
	return unspent, nil
}

func TestSendingGAS(t *testing.T) {
	//TEST WIF on privatenet
	wif := ""
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
	}

	unspent, err := utxoFromO3Platform("test", privateNetwallet.Address)
	if err != nil {
		log.Printf("error %v", err)
		t.Fail()
		return
	}
	asset := smartcontract.GAS
	amount := float64(0.1)
	toAddress := "AKo8k27H5nCG8MwSirmnraH6uUG6fQQVC2" //this is multi signature adddress 3/2
	to := smartcontract.ParseNEOAddress(toAddress)
	// remark := "O3TX"
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	// attributes[smartcontract.Remark1] = []byte(remark)

	fee := smartcontract.NetworkFeeAmount(0.0)
	nativeAsset := neoutils.UseNativeAsset(fee)
	rawtx, txid, err := nativeAsset.SendNativeAssetRawTransaction(*privateNetwallet, asset, amount, to, unspent, attributes)
	if err != nil {
		log.Printf("error sending natie %v", err)
		t.Fail()
		return
	}
	log.Printf("%v\n%x", txid, rawtx)
}

func TestSendingNEO(t *testing.T) {
	//TEST WIF on privatenet
	wif := ""
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
	}

	unspent, err := utxoFromO3Platform("private", privateNetwallet.Address)
	if err != nil {
		log.Printf("error %v", err)
		t.Fail()
		return
	}
	asset := smartcontract.GAS
	amount := float64(1000)
	toAddress := "Adm9ER3UwdJfimFtFhHq1L5MQ5gxLLTUes"
	to := smartcontract.ParseNEOAddress(toAddress)
	// remark := "O3TX"
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	// attributes[smartcontract.Remark1] = []byte(remark)

	fee := smartcontract.NetworkFeeAmount(0.0)
	nativeAsset := neoutils.UseNativeAsset(fee)
	rawtx, txid, err := nativeAsset.SendNativeAssetRawTransaction(*privateNetwallet, asset, amount, to, unspent, attributes)
	if err != nil {
		log.Printf("error sending natie %v", err)
		t.Fail()
		return
	}
	log.Printf("%v\n%x", txid, rawtx)
}

func TestSendingGASFromMultiSig(t *testing.T) {
	fromAddress := "AKo8k27H5nCG8MwSirmnraH6uUG6fQQVC2" //this is multi signature adddress 3/2
	unspent, err := utxoFromO3Platform("test", fromAddress)
	if err != nil {
		log.Printf("error %v", err)
		t.Fail()
		return
	}
	asset := smartcontract.GAS
	amount := float64(0.1)

	toAddress := "AKo8k27H5nCG8MwSirmnraH6uUG6fQQVC2"
	to := smartcontract.ParseNEOAddress(toAddress)

	attributes := map[smartcontract.TransactionAttribute][]byte{}

	fee := smartcontract.NetworkFeeAmount(0.0)
	nativeAsset := neoutils.UseNativeAsset(fee)
	rawtx, txid, err := nativeAsset.GenerateRawTx(fromAddress, asset, amount, to, unspent, attributes)
	if err != nil {
		log.Printf("error sending natie %v", err)
		t.Fail()
		return
	}
	log.Printf("txid %v\n", txid)
	log.Printf("raw %x\n", rawtx)

	wallet1, _ := neoutils.GenerateFromWIF("")
	wallet2, _ := neoutils.GenerateFromWIF("")

	wallets := []*neoutils.Wallet{wallet1, wallet2}

	signatures := []smartcontract.TransactionSignature{}

	for _, w := range wallets {
		privateKeyInHex := neoutils.BytesToHex(w.PrivateKey)

		signedData, err := neoutils.Sign(rawtx, privateKeyInHex)
		if err != nil {
			log.Printf("err signing %v", err)
			return
		}

		signature := smartcontract.TransactionSignature{
			SignedData: signedData,
			PublicKey:  w.PublicKey,
		}
		signatures = append(signatures, signature)

		log.Printf("pub key = %x\n", w.PublicKey)
		log.Printf("signedData = %x\n", signedData)
		hash := sha256.Sum256(rawtx)
		valid := neoutils.Verify(w.PublicKey, signedData, hash[:])
		log.Printf("valid %v", valid)
	}

	verificationScripts := smartcontract.NewScriptBuilder().GenerateVerificationScriptsMultiSig(signatures)

	//concat data
	endPayload := []byte{}
	endPayload = append(endPayload, rawtx...)
	endPayload = append(endPayload, verificationScripts...)

	redeemScript := "5221024da93f9a66981e499b36ce763e57fd89a47a052e86d40b42f81708c40fe9eff02102e77ff280db51ef3638009f11947c544ed094d4e5f2d96a9e654dc817bc3a898652ae"
	b := neoutils.HexTobytes(redeemScript)
	length := len(b)
	log.Printf("%x%x%v", endPayload, length, redeemScript)
}
