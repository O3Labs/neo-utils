package smartcontract

type Properties int

const (
	NoProperty       Properties = 0
	HasStorage       Properties = 1 << 0
	HasDynamicInvoke Properties = 1 << 1
	Payable          Properties = 1 << 2
)
