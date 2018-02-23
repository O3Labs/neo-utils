package rpc

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
