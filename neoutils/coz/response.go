package coz

type CoZUnspent struct {
	Index int     `json:"index"`
	Txid  string  `json:"txid"`
	Value float64 `json:"value"`
}
type UnspentBalance struct {
	GAS struct {
		Balance float64      `json:"balance"`
		Unspent []CoZUnspent `json:"unspent"`
	} `json:"GAS"`
	NEO struct {
		Balance int          `json:"balance"`
		Unspent []CoZUnspent `json:"unspent"`
	} `json:"NEO"`
	Address string `json:"address"`
	Net     string `json:"net"`
}

type ClaimResponse struct {
	Address string `json:"address"`
	Claims  []struct {
		Claim  int    `json:"claim"`
		End    int    `json:"end"`
		Index  int    `json:"index"`
		Start  int    `json:"start"`
		Sysfee int    `json:"sysfee"`
		Txid   string `json:"txid"`
		Value  int    `json:"value"`
	} `json:"claims"`
	Net               string `json:"net"`
	TotalClaim        int    `json:"total_claim"`
	TotalUnspentClaim int    `json:"total_unspent_claim"`
}
