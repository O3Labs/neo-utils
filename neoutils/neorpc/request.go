package neorpc

type JSONRPCRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

func NewRequest(method string, params []interface{}) JSONRPCRequest {
	return JSONRPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
}

type InvokeFunctionStackArg struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewInvokeFunctionStackByteArray(value string) InvokeFunctionStackArg {
	return InvokeFunctionStackArg{Type: "ByteArray", Value: value}
}
