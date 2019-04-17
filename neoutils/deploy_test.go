package neoutils_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

func TestReverseContractHash(t *testing.T) {
	log.Printf("%x", neoutils.ReverseBytes(neoutils.HexTobytes("5746e5e5e2a5ccd56930e82d771cb6dcc54e5005")))
}

func TestDeploy(t *testing.T) {

	contract := neoutils.SmartContractInfo{
		AVMHEX:      "54c56b6c766b00527ac46c766b51527ac46c766b00c36c766b51c3936c766b52527ac46203006c766b52c3616c7566",
		Name:        "Add",
		Version:     "v8",
		Author:      "James",
		Email:       "support@o3.network",
		Description: "Add",
		InputTypes:  []smartcontract.ParameterType{smartcontract.ByteArray},
		ReturnType:  smartcontract.ByteArray,
		Properties:  smartcontract.NoProperty,
	}

	log.Printf("sc hash %v", contract.GetScriptHash())

	asset := smartcontract.GAS
	amount := float64(100)

	// encryptedKey := ""
	// passphrase := ""
	// wif, _ := neoutils.NEP2Decrypt(encryptedKey, passphrase)
	wif := ""
	wallet, err := neoutils.GenerateFromWIF(wif)
	if err != nil {
		log.Printf("%v", err)
		t.Fail()
		return
	}
	log.Printf("wallet address %v", wallet.Address)

	// unspent, err := utxo("test", wallet.Address)
	// log.Printf("unspent %+v", unspent)
	// if err != nil {
	// 	log.Printf("error %v", err)
	// 	t.Fail()
	// 	return
	// }

	unspent := smartcontract.Unspent{
		Assets: map[smartcontract.NativeAsset]*smartcontract.Balance{},
	}

	gasBalance := smartcontract.Balance{
		Amount: float64(0) / float64(100000000),
		UTXOs:  []smartcontract.UTXO{},
	}

	gasTX := smartcontract.UTXO{
		Index: 0,
		TXID:  "0xbf3498830b8a3722ec0e876ac7abc3b4b802bd53cff1e11aa69a8c69b391ef49",
		Value: 18922,
	}
	gasBalance.UTXOs = append(gasBalance.UTXOs, gasTX)

	unspent.Assets[smartcontract.GAS] = &gasBalance

	attributes := map[smartcontract.TransactionAttribute][]byte{}
	tx, err := neoutils.DeploySmartContractScript(contract, *wallet, asset, amount, unspent, attributes)
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

func TestPropertyList(t *testing.T) {
	contract := neoutils.SmartContractInfo{
		AVMHEX:      "54c56b6c766b00527ac46c766b51527ac46c766b00c36c766b51c3936c766b52527ac46203006c766b52c3616c7566",
		Name:        "Add",
		Version:     "v8",
		Author:      "James",
		Email:       "support@o3.network",
		Description: "Add",
		InputTypes:  []smartcontract.ParameterType{smartcontract.ByteArray},
		ReturnType:  smartcontract.ByteArray,
		Properties:  smartcontract.NoProperty,
	}
	log.Printf("%x", contract.Serialize())
}
