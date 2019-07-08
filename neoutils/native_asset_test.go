package neoutils_test

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/neorpc"
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
	asset := smartcontract.NEO
	amount := float64(3)
	toAddress := "ANoW2zD8HmhbWJAjL4yKJWCZcF2WFb1ire" //this is multi signature adddress 3/2
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

func TestSendingNativeAsset(t *testing.T) {

	key := ""
	passphrase := ""
	privateNetwallet, err := neoutils.NEP2DecryptToWallet(key, passphrase)
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
	amount := float64(10)
	toAddress := privateNetwallet.Address
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
	log.Printf("%v\n%x %v", txid, rawtx, len(rawtx))
}

func TestSendingNEOFromMultiSig(t *testing.T) {
	fromAddress := "AXeKhuHRUXMJFAXLwyHxvyCNQb8X7mtnQU" //this is multi signature adddress 3/2
	neoclient := neorpc.NewClient("http://localhost:30333")
	unspentResponse := neoclient.GetUnspents(fromAddress)

	unspent := smartcontract.Unspent{
		Assets: map[smartcontract.NativeAsset]*smartcontract.Balance{},
	}

	gasBalance := smartcontract.Balance{
		Amount: float64(0),
		UTXOs:  []smartcontract.UTXO{},
	}

	neoBalance := smartcontract.Balance{
		Amount: float64(10000000),
		UTXOs:  []smartcontract.UTXO{},
	}

	for _, v := range unspentResponse.Result.Balance {
		for _, unspent := range v.Unspent {

			utxo := smartcontract.UTXO{
				TXID:  fmt.Sprintf("0x%v", unspent.Txid),
				Index: unspent.N,
				Value: float64(unspent.Value),
			}
			log.Printf("asset %+v", utxo)
			neoBalance.UTXOs = append(neoBalance.UTXOs, utxo)

		}
	}
	unspent.Assets[smartcontract.GAS] = &gasBalance
	unspent.Assets[smartcontract.NEO] = &neoBalance

	asset := smartcontract.NEO
	amount := float64(100 * 1000000)

	toAddress := "AVFobKv2y7i66gbGPAGDT67zv1RMQQj9GB"
	to := smartcontract.ParseNEOAddress(toAddress)

	attributes := map[smartcontract.TransactionAttribute][]byte{}

	fee := smartcontract.NetworkFeeAmount(0.0)
	nativeAsset := neoutils.UseNativeAsset(fee)
	rawtx, txid, err := nativeAsset.GenerateRawTx(fromAddress, asset, amount, to, unspent, attributes)
	if err != nil {
		log.Printf("error sending native %v", err)
		t.Fail()
		return
	}
	log.Printf("txid %v\n", txid)
	log.Printf("raw %x\n", rawtx)

	wallet1, _ := neoutils.GenerateFromWIF("")

	wallets := []*neoutils.Wallet{wallet1}

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

	redeemScript := "5121020ef8767aeb514780a8fb0a21f2568c521eb1e633a161dcdc39e78745762cb84351ae"
	b := neoutils.HexTobytes(redeemScript)
	length := len(b)
	log.Printf("%x%x%v", endPayload, length, redeemScript)
}
