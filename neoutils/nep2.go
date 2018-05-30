package neoutils

import "github.com/o3labs/neo-utils/neoutils/nep2"

type NEP2 struct {
	EncryptedKey string
	Address      string
}

func NEP2Encrypt(wif string, passphrase string) (*NEP2, error) {
	encryptedKey, address, err := nep2.NEP2Encrypt(wif, passphrase)
	return &NEP2{
		EncryptedKey: encryptedKey,
		Address:      address,
	}, err
}

func NEP2Decrypt(key, passphrase string) (s string, err error) {
	return nep2.NEP2Decrypt(key, passphrase)
}
