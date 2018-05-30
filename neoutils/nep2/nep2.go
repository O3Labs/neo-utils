package nep2

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/o3labs/neo-utils/neoutils/btckey"
	"github.com/o3labs/neo-utils/neoutils/crypto"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/text/unicode/norm"
)

//Source
//https://github.com/CityOfZion/neo-go/

// NEP-2 standard implementation for encrypting and decrypting wallets.

// NEP-2 specified parameters used for cryptography.
const (
	n       = 16384
	r       = 8
	p       = 8
	keyLen  = 64
	nepFlag = 0xe0
)

var nepHeader = []byte{0x01, 0x42}

type scryptParams struct {
	N int `json:"n"`
	R int `json:"r"`
	P int `json:"p"`
}

func newScryptParams() scryptParams {
	return scryptParams{
		N: n,
		R: r,
		P: p,
	}
}

// NEP2Encrypt encrypts a the PrivateKey using a given passphrase
// under the NEP-2 standard.
func NEP2Encrypt(wif string, passphrase string) (s string, address string, err error) {
	var privateKey btckey.PrivateKey
	errFromWIF := privateKey.FromWIF(wif)
	if err != nil {
		return "", "", errFromWIF
	}

	address = privateKey.ToNeoAddress()
	addressHash := hashAddress(address)[0:4]

	// Normalize the passphrase according to the NFC standard.
	phraseNorm := norm.NFC.Bytes([]byte(passphrase))
	derivedKey, err := scrypt.Key(phraseNorm, addressHash, n, r, p, keyLen)
	if err != nil {
		return s, "", err
	}
	derivedKey1 := derivedKey[:32]
	derivedKey2 := derivedKey[32:]

	xr := xor(privateKey.ToBytes(), derivedKey1)

	encrypted, err := crypto.AESEncrypt(xr, derivedKey2)
	if err != nil {
		return s, "", err
	}

	buf := new(bytes.Buffer)
	buf.Write(nepHeader)
	buf.WriteByte(nepFlag)
	buf.Write(addressHash)
	buf.Write(encrypted)

	if buf.Len() != 39 {
		return s, "", fmt.Errorf("invalid buffer length: expecting 39 bytes got %d", buf.Len())
	}

	return crypto.Base58CheckEncode(buf.Bytes()), address, nil
}

// NEP2Decrypt decrypts an encrypted key using a given passphrase
// under the NEP-2 standard.
func NEP2Decrypt(key, passphrase string) (s string, err error) {
	encrypted, err := crypto.Base58CheckDecode(key)
	if err != nil {
		return s, nil
	}
	if err := validateNEP2Format(encrypted); err != nil {
		return s, err
	}

	addrHash := encrypted[3:7]

	// Normalize the passphrase according to the NFC standard.
	phraseNorm := norm.NFC.Bytes([]byte(passphrase))
	derivedKey, err := scrypt.Key(phraseNorm, addrHash, n, r, p, keyLen)
	if err != nil {
		return s, err
	}

	derivedKey1 := derivedKey[:32]
	derivedKey2 := derivedKey[32:]
	encryptedBytes := encrypted[7:]

	decrypted, err := crypto.AESDecrypt(encryptedBytes, derivedKey2)
	if err != nil {
		return s, err
	}
	privBytes := xor(decrypted, derivedKey1)
	// Rebuild the private key.
	var privKey btckey.PrivateKey
	err = privKey.FromBytes(privBytes)
	if err != nil {
		return s, err
	}

	if !compareAddressHash(&privKey, addrHash) {
		return s, errors.New("password mismatch")
	}

	return privKey.ToWIFC(), nil
}

func compareAddressHash(priv *btckey.PrivateKey, hash []byte) bool {
	address := priv.ToNeoAddress()
	addrHash := hashAddress(address)[0:4]
	return bytes.Compare(addrHash, hash) == 0
}

func validateNEP2Format(b []byte) error {
	if len(b) != 39 {
		return fmt.Errorf("invalid length: expecting 39 got %d", len(b))
	}
	if b[0] != 0x01 {
		return fmt.Errorf("invalid byte sequence: expecting 0x01 got 0x%02x", b[0])
	}
	if b[1] != 0x42 {
		return fmt.Errorf("invalid byte sequence: expecting 0x42 got 0x%02x", b[1])
	}
	if b[2] != 0xe0 {
		return fmt.Errorf("invalid byte sequence: expecting 0xe0 got 0x%02x", b[2])
	}
	return nil
}

func xor(a, b []byte) []byte {
	if len(a) != len(b) {
		panic("cannot XOR non equal length arrays")
	}
	dst := make([]byte, len(a))
	for i := 0; i < len(dst); i++ {
		dst[i] = a[i] ^ b[i]
	}
	return dst
}

func hashAddress(addr string) []byte {
	sha := sha256.New()
	sha.Write([]byte(addr))
	hash := sha.Sum(nil)
	sha.Reset()
	sha.Write(hash)
	return sha.Sum(nil)
}
