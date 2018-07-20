package neoutils

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"

	nep9 "github.com/o3labs/NEP9-go/nep9"
	"github.com/o3labs/neo-utils/neoutils/btckey"
	"golang.org/x/crypto/ripemd160"
)

func ReverseBytes(b []byte) []byte {
	for i := 0; i < len(b)/2; i++ {
		j := len(b) - i - 1
		b[i], b[j] = b[j], b[i]
	}
	return b
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
// This method takes Big Endian Script hash
func ScriptHashToNEOAddress(scriptHash string) string {
	b := HexTobytes(scriptHash)
	//script hash from rpc or anything is always in big endian
	//to convert to a proper neo address
	//we need to reverse it first
	address := btckey.B58checkencodeNEO(0x17, ReverseBytes(b))
	return address
}

// // Convert NEO address to script hash
// func NEOAddressToScriptHash(neoAddress string) string {
// 	v, b, _ := btckey.B58checkdecode(neoAddress)
// 	if v != 0x17 {
// 		return ""
// 	}
// 	//reverse from little endian to big endian
// 	return fmt.Sprintf("%x", ReverseBytes(b))
// }

// Convert NEO address to script hash
func NEOAddressToScriptHashWithEndian(neoAddress string, endian binary.ByteOrder) string {
	v, b, _ := btckey.B58checkdecode(neoAddress)
	if v != 0x17 {
		return ""
	}
	if endian == binary.LittleEndian {
		return fmt.Sprintf("%x", b)
	} else {
		//reverse from little endian to big endian
		return fmt.Sprintf("%x", ReverseBytes(b))
	}
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
	v := new(big.Int).SetBytes(reversed)
	return v
}

type SimplifiedNEP9 struct {
	To     string  `json:"to"`
	Asset  string  `json:"assetID"`
	Amount float64 `json:"amount"`
}

func ParseNEP9URI(uri string) (*SimplifiedNEP9, error) {

	parsed, err := nep9.NewURI(uri)
	if err != nil {
		return nil, err
	}
	return &SimplifiedNEP9{
		To:     parsed.Address,
		Asset:  parsed.Asset,
		Amount: parsed.Amount,
	}, nil
}

func Hash160(data []byte) []byte {
	_, b, err := btckey.B58checkdecode(string(data))
	if err != nil {
		return nil
	}
	shortened := b[1 : len(b)-1]
	hex := bytesToHex(shortened)
	return ReverseBytes([]byte(hex))
}

func Hash256(b []byte) []byte {
	hash := sha256.Sum256(b)
	hash = sha256.Sum256(hash[:])
	return hash[:]
}

func PublicKeyToNEOAddress(publicKeyBytes []byte) string {
	publicKeyBytes = append([]byte{0x21}, publicKeyBytes...)
	publicKeyBytes = append(publicKeyBytes, 0xAC)

	/* SHA256 Hash */
	sha256_h := sha256.New()
	sha256_h.Reset()
	sha256_h.Write(publicKeyBytes)
	pub_hash_1 := sha256_h.Sum(nil)

	/* RIPEMD-160 Hash */
	ripemd160_h := ripemd160.New()
	ripemd160_h.Reset()
	ripemd160_h.Write(pub_hash_1)
	pub_hash_2 := ripemd160_h.Sum(nil)

	program_hash := pub_hash_2

	address := btckey.B58checkencodeNEO(0x17, program_hash)
	return address
}

func VMCodeToNEOAddress(vmCode []byte) string {
	/* SHA256 Hash */
	sha256_h := sha256.New()
	sha256_h.Reset()
	sha256_h.Write(vmCode)
	pub_hash_1 := sha256_h.Sum(nil)

	/* RIPEMD-160 Hash */
	ripemd160_h := ripemd160.New()
	ripemd160_h.Reset()
	ripemd160_h.Write(pub_hash_1)
	pub_hash_2 := ripemd160_h.Sum(nil)

	program_hash := pub_hash_2

	address := btckey.B58checkencodeNEO(0x17, program_hash)
	return address
}
