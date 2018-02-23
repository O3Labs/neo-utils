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
