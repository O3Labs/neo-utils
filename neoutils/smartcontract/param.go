package smartcontract

type ParameterType byte

const (
	Signature ParameterType = 0x00
	Boolean   ParameterType = 0x01

	Integer ParameterType = 0x02

	Hash160 ParameterType = 0x03

	Hash256 ParameterType = 0x04

	ByteArray ParameterType = 0x05
	PublicKey ParameterType = 0x06
	String    ParameterType = 0x07

	Array            ParameterType = 0x10
	InteropInterface ParameterType = 0xf0
	Void             ParameterType = 0xff
)

func (p ParameterType) Byte() byte {
	return byte(p)
}
