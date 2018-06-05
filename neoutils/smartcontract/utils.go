package smartcontract

import (
	"bytes"
	"encoding/binary"
	"math"
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

func RoundFixed8(val float64) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(8))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= 0.5 {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
