package neoutils

import (
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/o3labs/neo-utils/neoutils/btckey"

	"github.com/o3labs/neo-utils/neoutils/sss"
)

type Wallet struct {
	PublicKey       []byte
	PrivateKey      []byte
	Address         string
	WIF             string
	HashedSignature []byte
}

func hex2bytes(hexstring string) (b []byte) {
	b, _ = hex.DecodeString(hexstring)
	return b
}

func bytesToHex(b []byte) string {
	return hex.EncodeToString(b)
}

// Generate a wallet from a private key
func GeneratePublicKeyFromPrivateKey(privateKey string) (*Wallet, error) {
	pb := hex2bytes(privateKey)
	var priv btckey.PrivateKey
	err := priv.FromBytes(pb)
	if err != nil {
		return &Wallet{}, err
	}
	wallet := &Wallet{
		PublicKey:       priv.PublicKey.ToBytes(),
		PrivateKey:      priv.ToBytes(),
		Address:         priv.ToNeoAddress(),
		WIF:             priv.ToWIFC(),
		HashedSignature: priv.ToNeoSignature(),
	}
	return wallet, nil
}

// Generate a wallet from a WIF
func GenerateFromWIF(wif string) (*Wallet, error) {
	var priv btckey.PrivateKey
	err := priv.FromWIF(wif)
	if err != nil {
		return &Wallet{}, err
	}

	wallet := &Wallet{
		PublicKey:       priv.PublicKey.ToBytes(),
		PrivateKey:      priv.ToBytes(),
		Address:         priv.ToNeoAddress(),
		WIF:             priv.ToWIFC(),
		HashedSignature: priv.ToNeoSignature(),
	}
	return wallet, nil
}

// Create a new wallet.
func NewWallet() (*Wallet, error) {

	priv, err := btckey.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	wallet := &Wallet{
		PublicKey:       priv.PublicKey.ToBytes(),
		PrivateKey:      priv.ToBytes(),
		Address:         priv.ToNeoAddress(),
		WIF:             priv.ToWIFC(),
		HashedSignature: priv.ToNeoSignature(),
	}
	return wallet, nil
}

//Shared Secret with 2 parts.
type SharedSecret struct {
	First  []byte
	Second []byte
}

// Generate Shamir shared secret to SharedSecret struct.
func GenerateShamirSharedSecret(secret string) (*SharedSecret, error) {
	n := byte(2) // create 2 shares
	k := byte(2) // require 2 of them to combine

	shares, err := sss.Split(n, k, []byte(secret)) // split into 30 shares
	if err != nil {
		return nil, err
	}
	return &SharedSecret{
		First:  shares[byte(1)],
		Second: shares[byte(2)],
	}, nil
}

// Recover the secret from shared secrets.
func RecoverFromSharedSecret(first []byte, second []byte) (string, error) {
	s := SharedSecret{
		First:  []byte(first),
		Second: []byte(second),
	}
	if len(s.First) != len(s.Second) {
		return "", errors.New("Invalid shared secret")
	}
	k := byte(2)
	// select a random subset of the total shares
	subset := make(map[byte][]byte, k)
	subset[byte(1)] = s.First
	subset[byte(2)] = s.Second

	// combine two shares and recover the secret
	recovered := string(sss.Combine(subset))
	return recovered, nil
}

//Compute shared secret using ECDH
func (w *Wallet) ComputeSharedSecret(publicKey []byte) []byte {
	pKey := btckey.PublicKey{}
	pKey.FromBytes(publicKey)
	curve := elliptic.P256()
	x, _ := curve.ScalarMult(pKey.X, pKey.Y, w.PrivateKey)
	return x.Bytes()
}

// Sign data using ECDSA with a private key
func Sign(data []byte, key string) ([]byte, error) {
	return btckey.Sign(data, key)
}
