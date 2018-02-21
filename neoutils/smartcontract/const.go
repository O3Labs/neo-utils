package smartcontract

type NativeAsset string
type TradingVersion byte //currently 0
const (
	neo NativeAsset = "c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b"
	gas NativeAsset = "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7"
)

const (
	NEOTradingVersion TradingVersion = 0x00
)
