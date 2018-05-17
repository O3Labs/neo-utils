package neoutils_test

import (
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
	//TEST WIF on testnet
	wif := "L4Ns4Uh4WegsHxgDG49hohAYxuhj41hhxG6owjjTWg95GSrRRbLL"
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
	amount := float64(20)
	toAddress := "Adm9ER3UwdJfimFtFhHq1L5MQ5gxLLTUes"
	to := smartcontract.ParseNEOAddress(toAddress)
	remark := "O3TX"
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	attributes[smartcontract.Remark1] = []byte(remark)

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
