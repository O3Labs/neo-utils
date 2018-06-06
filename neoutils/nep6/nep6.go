package nep6

// https://github.com/neo-project/proposals/blob/master/nep-6.mediawiki

type NEP6Contract struct {
	//script is the script code of the contract. This field can be null if the contract has been deployed to the blockchain.
	Script string `json:"script,omitempty"`
	//parameters is an array of Parameter objects which describe the details of each parameter in the contract function. For more information about Parameter object, see the descriptions in NEP-3: NeoContract ABI.
	Parameters []interface{} `json:"parameters,omitempty"`
	// deployed indicates whether the contract has been deployed to the blockchain.
	Deployed bool `json:"deployed,omitempty"`
}

type NEP6Account struct {
	// address is the base58 encoded address of the account.
	Address string `json:"address"`
	// label is a label that the user has made to the account.
	Label string `json:"label"`
	// isDefault indicates whether the account is the default change account.
	IsDefault bool `json:"isDefault"`
	// lock indicates whether the account is locked by user. The client shouldn't spend the funds in a locked account.
	Lock bool `json:"lock"`
	// key is the private key of the account in the NEP-2 format. This field can be null (for watch-only address or non-standard address).
	Key string `json:"key"`
	// contract is a Contract object which describes the details of the contract. This field can be null (for watch-only address).
	Contract NEP6Contract `json:"contract"`
	// extra is an object that is defined by the implementor of the client for storing extra data. This field can be null.
	Extra interface{} `json:"extra"`
}

type NEP6Scrypt struct {
	N int `json:"n"` // 16384
	R int `json:"r"` // 8
	P int `json:"p"` // 8
}

type NEP6Wallet struct {
	Name     string        `json:"name"`
	Version  string        `json:"version"` //fixed to 1.0
	Scrypt   NEP6Scrypt    `json:"scrypt"`
	Accounts []NEP6Account `json:"accounts"`
	Extra    interface{}   `json:"extra"`
}

func NewNEP6WithNEP2EncryptedKey(name string, addressLabel string, address string, encryptedKey string) *NEP6Wallet {
	account := NEP6Account{
		Address:   address,
		Label:     addressLabel,
		IsDefault: true,
		Lock:      false,
		Key:       encryptedKey,
		Contract:  NEP6Contract{},
	}
	nep6 := NEP6Wallet{
		Name:    name,
		Version: "1.0",
		Scrypt: NEP6Scrypt{
			N: 16384,
			R: 8,
			P: 8,
		},
	}
	nep6.Accounts = append(nep6.Accounts, account)
	return &nep6
}
