package neoutils_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
)

func TestSerialize(t *testing.T) {

	data := `
{"sha256":"accb0534cb4c4d9b8594189d31e759ab96ae7488dc90e52a443b44bb2e2b2493","hash":"ef249c579898e3adaee6f4c5df8117cc08b8a2832bdd5978beeb859cef6620c9","inputs":[{"prevIndex":0,"prevHash":"3a963116b572a466819c05bee74782902a51fd9b83be99f25d9edc5b7891049a"},{"prevIndex":18,"prevHash":"b7e77a70481edc1d5156a182d358aac53da51b9f1653683ae7bcb811b779c759"}],"outputs":[{"assetId":"602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7","scriptHash":"e707714512577b42f9a011f8b870625429f93573","value":1e-08}],"script":"0800ca9a3b000000001432e125258b7db0a0dffde5bd03b2b859253538ab144d2c053c1a5911be6253b3cc1a68397feb3f647053c1076465706f73697467823b63e7c70a795a7615a38d1ba67d9e54c195a1","version":1,"type":209,"attributes":[{"usage":32,"data":"4d2c053c1a5911be6253b3cc1a68397feb3f6470"}],"scripts":[],"gas":0}
`
	tx := neoutils.NeonJSTransaction{}
	json.Unmarshal([]byte(data), &tx)

	final := neoutils.NeonJSTXSerializer(tx)
	log.Printf("%x", final)

	w := WalletForSwitcheoTestNet()
	log.Printf("%v", w.Address)
}
