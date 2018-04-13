package neoutils_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/coz"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

//TESTED with fee. succeeded
func TestSendNativeAsset(t *testing.T) {
	wif := "KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr"
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
	}
	cozClient := coz.NewClient("http://localhost:5000/")

	unspentCoz, err := cozClient.GetUnspentByAddress(privateNetwallet.Address)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	gasBalance := smartcontract.Balance{
		Amount: float64(0) / float64(100000000),
		UTXOs:  []smartcontract.UTXO{},
	}

	neoBalance := smartcontract.Balance{
		Amount: float64(0) / float64(100000000),
		UTXOs:  []smartcontract.UTXO{},
	}

	for _, v := range unspentCoz.GAS.Unspent {
		gasTX1 := smartcontract.UTXO{
			Index: v.Index,
			TXID:  v.Txid,
			Value: v.Value,
		}
		log.Printf("utxo value = %.8f", v.Value)
		gasBalance.UTXOs = append(gasBalance.UTXOs, gasTX1)
	}

	for _, v := range unspentCoz.NEO.Unspent {
		tx := smartcontract.UTXO{
			Index: v.Index,
			TXID:  v.Txid,
			Value: v.Value,
		}
		log.Printf("utxo value = %.8f", v.Value)
		neoBalance.UTXOs = append(neoBalance.UTXOs, tx)
	}

	unspent := smartcontract.Unspent{
		Assets: map[smartcontract.NativeAsset]*smartcontract.Balance{},
	}

	unspent.Assets[smartcontract.GAS] = &gasBalance
	unspent.Assets[smartcontract.NEO] = &neoBalance

	asset := smartcontract.NEO
	amount := float64(10)
	toAddress := "Adm9ER3UwdJfimFtFhHq1L5MQ5gxLLTUes"
	to := smartcontract.ParseNEOAddress(toAddress)
	remark := "O3TX"
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	attributes[smartcontract.Remark1] = []byte(remark)

	fee := smartcontract.NetworkFeeAmount(0.001)
	nativeAsset := neoutils.UseNativeAsset(fee)
	tx, err := nativeAsset.SendNativeAssetRawTransaction(*privateNetwallet, asset, amount, to, unspent, attributes)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("%x", tx)
}
