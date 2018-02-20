package neoutils

import (
	"encoding/hex"
	"fmt"
	"math/big"

	nep9 "github.com/o3labs/NEP9-go"
	"github.com/o3labs/neo-utils/neoutils/btckey"
)

func ReverseBytes(b []byte) []byte {
	// Protect from big.Ints that have 1 len bytes.
	if len(b) < 2 {
		return b
	}

	dest := make([]byte, len(b))
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		dest[i], dest[j] = b[j], b[i]
	}

	return dest
}

// Simple hex string to bytes
func HexTobytes(hexstring string) (b []byte) {
	b, _ = hex.DecodeString(hexstring)
	return b
}

// Simple bytes to Hex
func BytesToHex(b []byte) string {
	return hex.EncodeToString(b)
}

// Convert script hash to NEO address
func ScriptHashToNEOAddress(scriptHash string) string {
	b := HexTobytes(scriptHash)
	//script hash from rpc or anything is always in big endian
	//to convert to a proper neo address
	//we need to reverse it first
	address := btckey.B58checkencodeNEO(0x17, ReverseBytes(b))
	return address
}

// Convert NEO address to script hash
func NEOAddressToScriptHash(neoAddress string) string {
	v, b, _ := btckey.B58checkdecode(neoAddress)
	if v != 0x17 {
		return ""
	}
	//reverse from little endian to big endian
	return fmt.Sprintf("%x", ReverseBytes(b))
}

// Validate NEO address
func ValidateNEOAddress(address string) bool {
	//NEO address version is 23
	//https://github.com/neo-project/neo/blob/427a3cd08f61a33e98856e4b4312b8147708105a/neo/protocol.json#L4
	ver, _, err := btckey.B58checkdecode(address)
	if err != nil {
		return false
	}
	if ver != 23 {
		return false
	}
	return true
}

// Convert byte array to big int
func ConvertByteArrayToBigInt(hexString string) *big.Int {
	b, err := hex.DecodeString(hexString)
	if err != nil {
		return nil
	}
	reversed := ReverseBytes(b)
	reversedHex := hex.EncodeToString(reversed)
	v, _ := new(big.Int).SetString(reversedHex, 16)
	return v
}

type SimplifiedNEP9 struct {
	To      string  `json:"to"`
	AssetID string  `json:"assetID"`
	Amount  float64 `json:"amount"`
}

func ParseNEP9URI(uri string) (*SimplifiedNEP9, error) {
	parsed, err := nep9.ParseNEP9URI(uri)
	if err != nil {
		return nil, err
	}
	return &SimplifiedNEP9{
		To:      parsed.To,
		AssetID: parsed.AssetID,
		Amount:  parsed.Amount,
	}, nil
}
