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

func NEP2DecryptToWallet(key, passphrase string) (*Wallet, error) {
	priv, err := nep2.NEP2DecryptToPrivateKey(key, passphrase)
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
