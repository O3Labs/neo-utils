package neorpc

type GetUnspentsResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Balance []struct {
			Unspent []struct {
				Txid  string `json:"txid"`
				N     int    `json:"n"`
				Value int    `json:"value"`
			} `json:"unspent"`
			AssetHash   string `json:"asset_hash"`
			Asset       string `json:"asset"`
			AssetSymbol string `json:"asset_symbol"`
			Amount      int    `json:"amount"`
		} `json:"balance"`
		Address string `json:"address"`
	} `json:"result"`
}
