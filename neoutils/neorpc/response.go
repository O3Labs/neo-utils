package neorpc

type GetContractStateResult struct {
	Version     int      `json:"version"`
	Hash        string   `json:"hash"`
	Script      string   `json:"script"`
	Parameters  []string `json:"parameters"`
	Returntype  string   `json:"returntype"`
	Name        string   `json:"name"`
	CodeVersion string   `json:"code_version"`
	Author      string   `json:"author"`
	Email       string   `json:"email"`
	Description string   `json:"description"`
	Properties  struct {
		Storage       bool `json:"storage"`
		DynamicInvoke bool `json:"dynamic_invoke"`
	} `json:"properties"`
}

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type JSONRPCResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

type GetContractStateResponse struct {
	JSONRPCResponse
	*ErrorResponse                        //optional
	Result         GetContractStateResult `json:"result"`
}

type SendRawTransactionResponse struct {
	JSONRPCResponse
	*ErrorResponse //optional
	Result         bool
}

type GetRawTransactionResponse struct {
	JSONRPCResponse
	*ErrorResponse                         //optional
	Result         GetRawTransactionResult `json:"result"`
}

type GetRawTransactionResult struct {
	Txid       string `json:"txid"`
	Size       int    `json:"size"`
	Type       string `json:"type"`
	Version    int    `json:"version"`
	Attributes []struct {
		Usage string `json:"usage"`
		Data  string `json:"data"`
	} `json:"attributes"`
	Vin []struct {
		Txid string `json:"txid"`
		Vout int    `json:"vout"`
	} `json:"vin"`
	Vout []struct {
		N       int    `json:"n"`
		Asset   string `json:"asset"`
		Value   string `json:"value"`
		Address string `json:"address"`
	} `json:"vout"`
	SysFee  string `json:"sys_fee"`
	NetFee  string `json:"net_fee"`
	Scripts []struct {
		Invocation   string `json:"invocation"`
		Verification string `json:"verification"`
	} `json:"scripts"`
	Script        string `json:"script"`
	Gas           string `json:"gas"`
	Blockhash     string `json:"blockhash"`
	Confirmations int    `json:"confirmations"`
	Blocktime     int    `json:"blocktime"`
}

type GetBlockCountResponse struct {
	JSONRPCResponse
	*ErrorResponse     //optional
	Result         int `json:"result"`
}

type GetBlockResponse struct {
	JSONRPCResponse
	*ErrorResponse                //optional
	Result         GetBlockResult `json:"result"`
}

type GetBlockResult struct {
	Hash              string `json:"hash"`
	Size              int    `json:"size"`
	Version           int    `json:"version"`
	Previousblockhash string `json:"previousblockhash"`
	Merkleroot        string `json:"merkleroot"`
	Time              int    `json:"time"`
	Index             int    `json:"index"`
	Nonce             string `json:"nonce"`
	Nextconsensus     string `json:"nextconsensus"`
	Script            struct {
		Invocation   string `json:"invocation"`
		Verification string `json:"verification"`
	} `json:"script"`
	Tx []struct {
		Txid       string        `json:"txid"`
		Size       int           `json:"size"`
		Type       string        `json:"type"`
		Version    int           `json:"version"`
		Attributes []interface{} `json:"attributes"`
		Vin        []interface{} `json:"vin"`
		Vout       []interface{} `json:"vout"`
		SysFee     string        `json:"sys_fee"`
		NetFee     string        `json:"net_fee"`
		Scripts    []interface{} `json:"scripts"`
		Nonce      int64         `json:"nonce,omitempty"`
		Script     string        `json:"script,omitempty"`
		Gas        string        `json:"gas,omitempty"`
	} `json:"tx"`
	Confirmations int    `json:"confirmations"`
	Nextblockhash string `json:"nextblockhash"`
}
