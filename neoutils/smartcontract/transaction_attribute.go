package smartcontract

type TransactionAttribute byte

const (
	ContractHash TransactionAttribute = 0x00

	ECDH02 TransactionAttribute = 0x02

	ECDH03 TransactionAttribute = 0x03

	Script TransactionAttribute = 0x20

	Vote TransactionAttribute = 0x30

	DescriptionUrl TransactionAttribute = 0x81
	Description    TransactionAttribute = 0x90

	Hash1  TransactionAttribute = 0xa1
	Hash2  TransactionAttribute = 0xa2
	Hash3  TransactionAttribute = 0xa3
	Hash4  TransactionAttribute = 0xa4
	Hash5  TransactionAttribute = 0xa5
	Hash6  TransactionAttribute = 0xa6
	Hash7  TransactionAttribute = 0xa7
	Hash8  TransactionAttribute = 0xa8
	Hash9  TransactionAttribute = 0xa9
	Hash10 TransactionAttribute = 0xaa
	Hash11 TransactionAttribute = 0xab
	Hash12 TransactionAttribute = 0xac
	Hash13 TransactionAttribute = 0xad
	Hash14 TransactionAttribute = 0xae
	Hash15 TransactionAttribute = 0xaf

	Remark   TransactionAttribute = 0xf0
	Remark1  TransactionAttribute = 0xf1
	Remark2  TransactionAttribute = 0xf2
	Remark3  TransactionAttribute = 0xf3
	Remark4  TransactionAttribute = 0xf4
	Remark5  TransactionAttribute = 0xf5
	Remark6  TransactionAttribute = 0xf6
	Remark7  TransactionAttribute = 0xf7
	Remark8  TransactionAttribute = 0xf8
	Remark9  TransactionAttribute = 0xf9
	Remark10 TransactionAttribute = 0xfa
	Remark11 TransactionAttribute = 0xfb
	Remark12 TransactionAttribute = 0xfc
	Remark13 TransactionAttribute = 0xfd
	Remark14 TransactionAttribute = 0xfe
	Remark15 TransactionAttribute = 0xf
)

func (t TransactionAttribute) ToByte() byte {
	return byte(t)
}
