package neoutils_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func TestDeploy(t *testing.T) {

	contract := neoutils.SmartContractInfo{
		AVMHEX:      "",
		Name:        "Test Contract",
		Version:     "1.0",
		Author:      "Apisit",
		Email:       "apisit@o3.network",
		Description: "https://o3.network",
		InputTypes:  []smartcontract.ParameterType{smartcontract.String, smartcontract.Array},
		ReturnType:  smartcontract.ByteArray,
		Properties:  smartcontract.HasStorage + smartcontract.Payable,
	}

	log.Printf("sc hash %v", contract.GetScriptHash())

	asset := smartcontract.GAS
	amount := float64(490)

	encryptedKey := ""
	passphrase := ""
	wif, _ := neoutils.NEP2Decrypt(encryptedKey, passphrase)

	privateNetwallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	log.Printf("wallet address %v", privateNetwallet.Address)

	unspent, err := utxo("main", privateNetwallet.Address)
	log.Printf("unspent %+v", unspent)
	if err != nil {
		log.Printf("error %v", err)
		t.Fail()
		return
	}
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	tx, err := neoutils.DeploySmartContractScript(contract, *privateNetwallet, asset, amount, unspent, attributes)
	if err != nil {
		log.Printf("error %v", err)
		return
	}
	log.Printf("tx %x", tx)

}

func TestVarInt(t *testing.T) {
	buff := new(bytes.Buffer)
	neoutils.WriteVarUint(buff, uint64(286))
	log.Printf("%x", buff.Bytes())
}
