package neoutils

import (
	"log"

	"github.com/o3labs/ont-mobile/ontmobile"
)

func OntologyTransfer(endpoint string, gasPrice uint, gasLimit uint, wif string, asset string, to string, amount float64) (string, error) {
	raw, err := ontmobile.Transfer(gasPrice, gasLimit, wif, asset, to, amount)
	if err != nil {
		return "", err
	}
	log.Printf("%x", raw.Data)
	return "", nil
	// txid, err := ontmobile.SendRawTransaction(endpoint, fmt.Sprintf("%x", raw.Data))
	// if err != nil {
	// 	return "", err
	// }

	// return txid, nil
}
