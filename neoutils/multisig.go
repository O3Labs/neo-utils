//WARNING: not finish
package neoutils

import (
	"fmt"
	"sort"

	"github.com/o3labs/neo-utils/neoutils/btckey"
	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

type MultiSigInterface interface {
	CreateMultiSigRedeemScript(numerOfRequiredSignature int, publicKeys [][]byte) ([]byte, error)
}

type MultiSig struct{}

var _ MultiSigInterface = (*MultiSig)(nil)

func (m *MultiSig) CreateMultiSigRedeemScript(numerOfRequiredSignature int, publicKeys [][]byte) ([]byte, error) {
	numberOfPublicKeys := len(publicKeys)
	if numberOfPublicKeys <= 1 {
		return nil, fmt.Errorf("Number of required Signature must be more than one")
	}
	if numerOfRequiredSignature > numberOfPublicKeys {
		return nil, fmt.Errorf("Number of required Signature is more than public keys provided.")
	}

	//sort public key
	keys := []btckey.PublicKey{}
	for _, pb := range publicKeys {
		publicKey := btckey.PublicKey{}
		publicKey.FromBytes(pb)
		keys = append(keys, publicKey)
	}

	//https://golang.org/pkg/math/big/#Int.Cmp
	sort.SliceStable(keys, func(i, j int) bool { return keys[i].Point.X.Cmp(keys[j].Point.X) == -1 })

	sb := smartcontract.NewScriptBuilder()
	sb.Push(numerOfRequiredSignature)
	for _, publicKey := range keys {
		sb.Push(publicKey.ToBytes())
	}
	sb.Push(numberOfPublicKeys)
	sb.PushOpCode(smartcontract.CHECKMULTISIG)
	return sb.ToBytes(), nil
}
