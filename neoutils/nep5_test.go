package neoutils_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/coz"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func unspent(address string) (smartcontract.Unspent, error) {
	cozClient := coz.NewClient("http://localhost:5000/")

	unspentCoz, err := cozClient.GetUnspentByAddress(address)
	if err != nil {
		log.Printf("%v", err)
		return smartcontract.Unspent{}, err
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
	return unspent, nil
}

//TEST with fee. succeeded
func TestTransferNEP5(t *testing.T) {

	scripthash := "b2eb148d3783f60e678e35f2c496de1a2a7ead93"
	fee := smartcontract.NetworkFeeAmount(1)
	nep5 := neoutils.UseNEP5WithNetworkFee(scripthash, fee)

	wif := "KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr"
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	amount := float64(1)
	to := smartcontract.ParseNEOAddress("Adm9ER3UwdJfimFtFhHq1L5MQ5gxLLTUes")

	unspent, err := unspent(privateNetwallet.Address)
	if err != nil {
		t.Fail()
		return
	}

	remark := "O3TX"
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	attributes[smartcontract.Remark1] = []byte(remark)

	tx, err := nep5.TransferNEP5RawTransaction(*privateNetwallet, to, amount, unspent, attributes)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("%x", tx)

}
