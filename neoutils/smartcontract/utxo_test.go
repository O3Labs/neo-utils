package smartcontract

import (
	"log"
	"testing"
)

func TestUTXOStruct(t *testing.T) {
	gasTX2 := UTXO{
		Index: 0,
		TXID:  "ad8d65c22de1873dea36587a989a4563c7264c48ed20a6edbe957bbe428984c0",
		Value: 40.0,
	}
	gasTX1 := UTXO{
		Index: 0,
		TXID:  "1b640fc70e127a74ab6785afe155f089e08a153b2effc7a4bed8b6690cfc65fe",
		Value: 7608.0,
	}

	gasBalance := Balance{
		Amount: 7648.0,
		UTXOs:  []UTXO{gasTX1, gasTX2},
	}

	log.Printf("total amount = %v", gasBalance.TotalAmount())

	neoTX1 := UTXO{
		Index: 0,
		TXID:  "e8b8bf4f98490368fc1caa86f8646e7383bb52751ffc3a1a7e296d715c4382ed",
		Value: 100000000,
	}

	neoBalance := Balance{
		Amount: 100000000,
		UTXOs:  []UTXO{neoTX1},
	}

	log.Printf("total amount = %v", neoBalance.TotalAmount())

	unspent := Unspent{
		Assets: map[NativeAsset]*Balance{},
	}
	unspent.Assets[NEO] = &neoBalance
	unspent.Assets[GAS] = &gasBalance

	log.Printf("%+v", unspent)
}
func TestSortUTXOMinFirst(t *testing.T) {

	gasTX1 := UTXO{
		Index: 0,
		TXID:  "1b640fc70e127a74ab6785afe155f089e08a153b2effc7a4bed8b6690cfc65fe",
		Value: 7608.0,
	}

	gasTX2 := UTXO{
		Index: 0,
		TXID:  "ad8d65c22de1873dea36587a989a4563c7264c48ed20a6edbe957bbe428984c0",
		Value: 40.0,
	}

	gasBalance := Balance{
		Amount: 7648.0,
		UTXOs:  []UTXO{gasTX1, gasTX2},
	}
	log.Printf("before sort %+v", gasBalance)
	gasBalance.SortMinFirst()
	log.Printf("after sort %+v", gasBalance)
}
