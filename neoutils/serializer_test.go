package neoutils_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
)

func TestSerialize(t *testing.T) {

	data := `{
  "attributes": [
    {
      "data": "d3c3e2d374c2bb9471e7e66010caf176cb407a88",
      "usage": 32
    }
  ],
  "gas": 0,
  "inputs": [
    {
      "prevHash": "d9e0a58751bf19e31df29723c568fbc385d06f5475f4ebd4e3294ea030125331",
      "prevIndex": 0
    }
  ],
  "outputs": [
    {
      "assetId": "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7",
      "scriptHash": "887a40cb76f1ca1060e6e77194bbc274d3e2c3d3",
      "value": "1"
    }
  ],
  "script": "14d3c3e2d374c2bb9471e7e66010caf176cb407a8851c106726566756e646776db3192722022eb7841038246dc8fa636dcf274",
  "scripts": [
    {
      "invocationScript": "0000",
      "verificationScript": ""
    },
    {
      "invocationScript": "404ca49e325cb30df1d2bebdd0129094c8a89863125251858143e0739b3996960456b1fda8c6cd06fea4160a2a76d131e4cf2c9b251bddca7d88b4c721beb08102",
      "verificationScript": "21023966fbe8c68f82a05c3a158c06b564a28661dd094f73bc3a9afbb73132562e5eac"
    }
  ],
  "type": 209,
  "version": 1
}`
	tx := neoutils.NeonJSTransaction{}
	json.Unmarshal([]byte(data), &tx)

	final := neoutils.NeonJSTXSerializer(tx)
	log.Printf("%x", final)

}
