//WARNING: not finish
package neoutils

import (
	"fmt"
	"sort"

	"github.com/o3labs/neo-utils/neoutils/btckey"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

type MultiSigInterface interface {
	CreateMultiSigRedeemScript() ([]byte, error)
}

type MultiSig struct {
	NumberOfRequiredSignatures int
	PublicKeys                 [][]byte
}

var _ MultiSigInterface = (*MultiSig)(nil)

func (m *MultiSig) CreateMultiSigRedeemScript() ([]byte, error) {

	numberOfPublicKeys := len(m.PublicKeys)
	if numberOfPublicKeys <= 1 {
		return nil, fmt.Errorf("Number of required Signature must be more than one")
	}
	if m.NumberOfRequiredSignatures > numberOfPublicKeys {
		return nil, fmt.Errorf("Number of required Signature is more than public keys provided.")
	}

	//sort public key
	keys := []btckey.PublicKey{}
	for _, pb := range m.PublicKeys {
		publicKey := btckey.PublicKey{}
		publicKey.FromBytes(pb)
		keys = append(keys, publicKey)
	}

	//https://golang.org/pkg/math/big/#Int.Cmp
	sort.SliceStable(keys, func(i, j int) bool { return keys[i].Point.X.Cmp(keys[j].Point.X) == -1 })

	sb := smartcontract.NewScriptBuilder()
	sb.Push(m.NumberOfRequiredSignatures)
	for _, publicKey := range keys {
		sb.Push(publicKey.ToBytes())
	}
	sb.Push(numberOfPublicKeys)
	sb.PushOpCode(smartcontract.CHECKMULTISIG)
	return sb.ToBytes(), nil
}
