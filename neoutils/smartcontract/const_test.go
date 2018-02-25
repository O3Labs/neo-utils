package smartcontract

import (
	"log"
	"testing"
)

func TestNativeAsset(t *testing.T) {
	log.Printf("%x", NEO.ToLittleEndianBytes())
}
