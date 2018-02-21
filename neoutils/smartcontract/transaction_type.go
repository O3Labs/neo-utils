package smartcontract

type TransactionType byte

const (
	MinerTransaction      TransactionType = 0x00
	IssueTransaction      TransactionType = 0x01
	ClaimTransaction      TransactionType = 0x02
	EnrollmentTransaction TransactionType = 0x20
	RegisterTransaction   TransactionType = 0x40
	ContractTransaction   TransactionType = 0x80
	StateTransaction      TransactionType = 0x90
	PublishTransaction    TransactionType = 0xd0
	InvocationTransaction TransactionType = 0xd1
)
