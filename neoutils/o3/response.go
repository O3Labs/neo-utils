package o3

type Response struct {
	Code int `json:"code"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

type UTXOResponse struct {
	Response
	*ErrorResponse
	Result struct {
		Data []UTXOResultData `json:"data"`
	} `json:"result"`
}

type UTXOResultData struct {
	Asset          string `json:"asset"`
	Index          int    `json:"index"`
	Txid           string `json:"txid"`
	Value          string `json:"value"`
	CreatedAtBlock int    `json:"createdAtBlock"`
}

type ClaimableGASResponse struct {
	Response
	*ErrorResponse
	Result struct {
		Data struct {
			Gas    string `json:"gas"`
			Claims []struct {
				Asset          string `json:"asset"`
				Index          int    `json:"index"`
				Txid           string `json:"txid"`
				Value          string `json:"value"`
				CreatedAtBlock int    `json:"createdAtBlock"`
			} `json:"claims"`
		} `json:"data"`
	} `json:"result"`
}
