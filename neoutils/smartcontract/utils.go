package smartcontract

import (
	"bytes"
	"encoding/binary"
)

func reverseBytes(b []byte) []byte {
	if len(b) < 2 {
		return b
	}

	dest := make([]byte, len(b))
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		dest[i], dest[j] = b[j], b[i]
	}

	return dest
}

func uintToBytes(value uint) []byte {
	countBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(countBytes, uint64(value))
	return bytes.TrimRight(countBytes, "\x00")
}

func uint16ToFixBytes(value uint16) []byte {
	countBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(countBytes, value)
	return countBytes //bytes.TrimRight(countBytes, "\x00")
}
