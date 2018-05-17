package smartcontract

import (
	"bytes"
	"encoding/binary"
)

func reverseBytes(b []byte) []byte {
	for i := 0; i < len(b)/2; i++ {
		j := len(b) - i - 1
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func uintToBytes(value uint) []byte {
	countBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(countBytes, uint64(value))
	return bytes.TrimRight(countBytes, "\x00")
}

func uint16ToFixBytes(value uint16) []byte {
	countBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(countBytes, value)
	return countBytes
}
