package neoutils

import (
	"log"

	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

// This class contains simplified method designed specifically for gomobile bind
// gomobile bind doesn't support slice argument or return

func MintTokens(scriptHash string, wallet Wallet, assetToSend smartcontract.NativeAsset, amount float64, unspent smartcontract.Unspent, remark string) (bool, error) {
	sc := UseSmartContract(scriptHash)

	operation := "mintTokens"
	args := []interface{}{}
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	attributes[smartcontract.Remark1] = []byte(remark)

	tx, err := sc.GenerateInvokeFunctionRawTransactionWithAmountToSend(wallet, assetToSend, amount, unspent, attributes, operation, args)
	if err != nil {
		return false, err
	}
	log.Printf("%x", tx)
	return false, nil
}
