package neoutils

import (
	"fmt"

	"github.com/o3labs/neo-utils/neoutils/smartcontract"
)

type MultiSigInterface interface {
	CreateMultiSignedAddress(publicKeys [][]byte) error
	CreateMultiSigRedeemScript(numerOfRequiredSignature int, publicKeys [][]byte) ([]byte, error)
}

type MultiSig struct{}

var _ MultiSigInterface = (*MultiSig)(nil)

//TODO: finish this
func (m *MultiSig) CreateMultiSignedAddress(publicKeys [][]byte) error {
	return nil
}

//TODO: finish this
func (m *MultiSig) CreateMultiSigRedeemScript(numerOfRequiredSignature int, publicKeys [][]byte) ([]byte, error) {
	numberOfPublicKeys := len(publicKeys)
	if numberOfPublicKeys <= 1 {
		return nil, fmt.Errorf("Number of required Signature must be more than one")
	}
	if numerOfRequiredSignature > numberOfPublicKeys {
		return nil, fmt.Errorf("Number of required Signature is more than public keys provided.")
	}

	sb := smartcontract.NewScriptBuilder()
	sb.EmitPush(numerOfRequiredSignature)
	for _, publicKey := range publicKeys {
		// pub := btckey.PublicKey{}
		// err := pub.FromBytes(publicKey)
		// if err != nil {
		// 	return err
		// }
		//this publicKey is already a compressed bytes so we can add that directly
		sb.EmitPush(publicKey)
	}

	sb.EmitPush(numberOfPublicKeys)
	sb.EmitPush(smartcontract.CHECKMULTISIG)
	return sb.ToBytes(), nil
}
