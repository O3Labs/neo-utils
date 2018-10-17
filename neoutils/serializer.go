package neoutils

import (
	"log"
	"strconv"

	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

type NeonJSTransaction struct {
	Sha256 string `json:"sha256"`
	Hash   string `json:"hash"`
	Inputs []struct {
		PrevIndex int    `json:"prevIndex"`
		PrevHash  string `json:"prevHash"`
	} `json:"inputs"`
	Outputs []struct {
		AssetID    string      `json:"assetId"`
		ScriptHash string      `json:"scriptHash"`
		Value      interface{} `json:"value"`
	} `json:"outputs"`
	Script     string `json:"script"`
	Version    int    `json:"version"`
	Type       int    `json:"type"`
	Attributes []struct {
		Usage int    `json:"usage"`
		Data  string `json:"data"`
	} `json:"attributes"`
	Scripts []interface{} `json:"scripts"`
	Gas     int           `json:"gas"`
}

func NeonJSTXSerializer(tx NeonJSTransaction) []byte {

	hexType := strconv.FormatInt(int64(tx.Type), 16)

	transaction := smartcontract.NewTransactionWithType(hex2bytes(hexType)[0], tx.Version)
	//inputs
	inputs := []smartcontract.UTXO{}
	for _, v := range tx.Inputs {
		input := smartcontract.UTXO{
			Index: v.PrevIndex,
			TXID:  v.PrevHash,
		}
		inputs = append(inputs, input)
	}

	inputBuilder := smartcontract.NewScriptBuilder()
	inputBuilder.PushLength(len(inputs))
	for _, input := range inputs {
		inputBuilder.Push(input)
	}
	transaction.Inputs = inputBuilder.ToBytes()
	//end inputs

	//output
	outputs := []smartcontract.TransactionOutput{}
	for _, v := range tx.Outputs {
		address := smartcontract.NEOAddressFromScriptHash(HexTobytes(v.ScriptHash))
		valueFloat := float64(0)
		switch a := v.Value.(type) {
		case float64:
			valueFloat = a
		case string:
			valueFloat, _ = strconv.ParseFloat(a, 64)
		}

		output := smartcontract.TransactionOutput{
			Asset:   smartcontract.NativeAssets[v.AssetID],
			Value:   int64(smartcontract.RoundFixed8(valueFloat) * float64(100000000)),
			Address: address,
		}
		outputs = append(outputs, output)
	}

	outputBuilder := smartcontract.NewScriptBuilder()

	outputBuilder.PushLength(len(outputs))
	for _, output := range outputs {
		outputBuilder.Push(output)
	}
	transaction.Outputs = outputBuilder.ToBytes()
	//end outputs

	//attributes
	attributes := map[smartcontract.TransactionAttribute][]byte{}
	for _, v := range tx.Attributes {
		usageHex := strconv.FormatInt(int64(v.Usage), 16)
		attr := smartcontract.TransactionAttribute(HexTobytes(usageHex)[0])
		attributes[attr] = hex2bytes(v.Data)
	}

	txAttributes, err := smartcontract.NewScriptBuilder().GenerateTransactionAttributes(attributes)

	if err == nil {
		transaction.Attributes = txAttributes
	}
	//end attributes
	scriptBytes := HexTobytes(tx.Script)

	transaction.Data = append([]byte{byte(len(scriptBytes))}, scriptBytes...)
	transaction.GAS = uint64(tx.Gas)

	log.Printf("tx id %v", transaction.ToTXID())

	return transaction.ToBytes()
}
