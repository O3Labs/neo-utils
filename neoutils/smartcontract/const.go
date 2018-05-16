package smartcontract

import "encoding/hex"

type NativeAsset string
type TradingVersion byte //currently 0

type NetworkFeeAmount float64

type Operation string

const scripthashLength = 20

const (
	NEO NativeAsset = "c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b"
	GAS NativeAsset = "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7"
)

var NativeAssets = map[string]NativeAsset{
	"c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b": NEO,
	"602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7": GAS,
}

func (n NativeAsset) ToLittleEndianBytes() []byte {
	b, err := hex.DecodeString(string(n))
	if err != nil {
		return nil
	}
	return reverseBytes(b)
}

const (
	NEOTradingVersion           TradingVersion = 0x00
	NEOTradingVersionPayableGAS TradingVersion = 0x01
)
