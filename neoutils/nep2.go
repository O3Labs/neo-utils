package neoutils

import "github.com/o3labs/neo-utils/neoutils/nep2"

func NEP2Encrypt(wif string, passphrase string) (s string, err error) {
	return nep2.NEP2Encrypt(wif, passphrase)
}

func NEP2Decrypt(key, passphrase string) (s string, err error) {
	return nep2.NEP2Decrypt(key, passphrase)
}
