package smartcontract

import (
	"sort"
)

type UTXO struct {
	Index int
	TXID  string
	Value float64
}

type Balance struct {
	Amount float64
	UTXOs  []UTXO
}

func (b *Balance) TotalAmount() float64 {
	total := float64(0)
	for _, v := range b.UTXOs {
		total += v.Value
	}
	return total
}

func (b *Balance) SortMinFirst() {
	sort.SliceStable(b.UTXOs, func(i, j int) bool {
		return b.UTXOs[i].Value < b.UTXOs[j].Value
	})
}

type Unspent struct {
	Assets map[NativeAsset]*Balance
}
