package neoutils_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
)

func TestGAS2CGAS(t *testing.T) {
	wif := ""
	wallet, _ := neoutils.GenerateFromWIF(wif)
	amountToConvert := float64(10)
	networkfee := float64(0)
	tx, txID, err := neoutils.GASToCGAS(*wallet, amountToConvert, networkfee)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("txID: %v\n raw:%x", txID, tx)
}
