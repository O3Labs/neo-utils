package neoutils_test

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/o3"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func TestMintTokens(t *testing.T) {
	scripthash := ""
	fee := smartcontract.NetworkFeeAmount(0)
	nep5 := neoutils.UseNEP5WithNetworkFee(scripthash, fee)

	wif := ""
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	unspent := smartcontract.Unspent{}

	remark := "APISIT FROM O3 IS HERE."

	asset := smartcontract.NEO
	amount := float64(10)

	tx, txID, err := nep5.MintTokensRawTransaction(*privateNetwallet, asset, amount, unspent, remark)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("txID = %v", txID)
	log.Printf("%x", tx)
}

func utxo(network string, address string) (smartcontract.Unspent, error) {

	unspent := smartcontract.Unspent{
		Assets: map[smartcontract.NativeAsset]*smartcontract.Balance{},
	}

	client := o3.DefaultO3APIClient()
	if network == "test" {
		client = o3.APIClientWithNEOTestnet()
	}

	if network == "private" {
		client = o3.APIClientWithNEOPrivateNet()
	}

	response := client.GetNEOUTXO(address)
	if response.Code != 200 {
		return unspent, fmt.Errorf("Error cannot get utxo")
	}

	gasBalance := smartcontract.Balance{
		Amount: float64(0) / float64(100000000),
		UTXOs:  []smartcontract.UTXO{},
	}

	neoBalance := smartcontract.Balance{
		Amount: float64(0) / float64(100000000),
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

//TEST with fee. succeeded
func TestTransferNEP5PrivateNet(t *testing.T) {

	//this is NNC token
	scripthash := "0xe8fe7fbaf639722e577a1961b9cc1d43572ed6c3"
	fee := smartcontract.NetworkFeeAmount(0)
	nep5 := neoutils.UseNEP5WithNetworkFee(scripthash, fee)

	wif := ""
	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	log.Printf("wallet address %v", privateNetwallet.Address)
	hash := neoutils.NEOAddressToScriptHashWithEndian(privateNetwallet.Address, binary.LittleEndian)
	log.Printf("wallet hash %v", hash)

	amount := float64(1)
	to := smartcontract.ParseNEOAddress("AeNkbJdiMx49kBStQdDih7BzfDwyTNVRfb")

	unspent, err := utxo("main", privateNetwallet.Address) //smartcontract.Unspent{}
	log.Printf("unspent %+v", unspent)
	if err != nil {
		log.Printf("error %v", err)
		t.Fail()
		return
	}
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	//address hash is a hex string
	// addressScriptHash := neoutils.NEOAddressToScriptHashWithEndian(privateNetwallet.Address, binary.LittleEndian)
	// b, _ := hex.DecodeString(addressScriptHash)
	// attributes[smartcontract.Script] = []byte(b)
	// attributes[smartcontract.Remark1] = []byte(fmt.Sprintf("O3TX%v", time.Now().Unix()))

	tx, txID, err := nep5.TransferNEP5RawTransaction(*privateNetwallet, to, amount, unspent, attributes)
	if err != nil {
		t.Fail()
		return
	}
	log.Printf("txID %v ", txID)
	log.Printf("%x", tx)
}
