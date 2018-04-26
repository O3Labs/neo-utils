package neoutils

import (
	"fmt"

	"github.com/o3labs/neo-utils/neoutils/coz"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

// This class contains simplified method designed specifically for gomobile bind
// gomobile bind doesn't support slice argument or return

func utxoFromNEONWalletDB(neonWalletDBEndpoint string, address string) (smartcontract.Unspent, error) {
	//"http://localhost:5000/"
	cozClient := coz.NewClient(neonWalletDBEndpoint)

	unspentCoz, err := cozClient.GetUnspentByAddress(address)
	if err != nil {
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
		gasBalance.UTXOs = append(gasBalance.UTXOs, gasTX1)
	}

	for _, v := range unspentCoz.NEO.Unspent {
		tx := smartcontract.UTXO{
			Index: v.Index,
			TXID:  v.Txid,
			Value: v.Value,
		}
		neoBalance.UTXOs = append(neoBalance.UTXOs, tx)
	}

	unspent := smartcontract.Unspent{
		Assets: map[smartcontract.NativeAsset]*smartcontract.Balance{},
	}

	unspent.Assets[smartcontract.GAS] = &gasBalance
	unspent.Assets[smartcontract.NEO] = &neoBalance
	return unspent, nil
}

type RawTransaction struct {
	TXID string
	Data []byte
}

func MintTokensRawTransactionMobile(utxoEndpoint string, scriptHash string, wif string, sendingAssetID string, amount float64, remark string, networkFeeAmountInGAS float64) (*RawTransaction, error) {
	rawTransaction := &RawTransaction{}
	fee := smartcontract.NetworkFeeAmount(networkFeeAmountInGAS)
	nep5 := UseNEP5WithNetworkFee(scriptHash, fee)
	wallet, err := GenerateFromWIF(wif)
	if err != nil {
		return nil, err
	}

	unspent, err := utxoFromNEONWalletDB(utxoEndpoint, wallet.Address)
	if err != nil {
		return nil, err
	}

	nativeAsset := smartcontract.NativeAssets[sendingAssetID]
	if nativeAsset == "" {
		return nil, fmt.Errorf("invalid assetID")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("Invalid amount. cannot be zero or less than zero")
	}

	data, txIDString, err := nep5.MintTokensRawTransaction(*wallet, nativeAsset, amount, unspent, remark)
	if err != nil {
		return nil, err
	}
	rawTransaction.Data = data
	rawTransaction.TXID = txIDString
	return rawTransaction, nil
}
