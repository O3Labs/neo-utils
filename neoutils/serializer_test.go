package neoutils_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils"
)

func TestNEONJSTX(t *testing.T) {

	str := `{
  "type": 209,
  "version": 1,
  "attributes": [
    {
      "usage": 32,
      "data": "06360b85b04bded387aa633bcc4bdda4354b5493"
    },
    {
      "usage": 240,
      "data": "4f33583135363038323937313537343934356533663630343364366562333734"
    }
  ],
  "inputs": [
    {
      "prevHash": "1ebc1e25fe53541adab91e4f97b5f28d4feb70495335790f20df44658937c3fb",
      "prevIndex": 1
    }
  ],
  "outputs": [
    {
      "assetId": "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7",
      "value": "0.1",
      "scriptHash": "06360b85b04bded387aa633bcc4bdda4354b5493"
    },
    {
      "assetId": "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7",
      "value": "9.1989",
      "scriptHash": "887a40cb76f1ca1060e6e77194bbc274d3e2c3d3"
    }
  ],
  "scripts": [
    {
      "invocationScript": "0000",
      "verificationScript": ""
    },
    {
      "invocationScript": "4026588c100e7eb3bc1a573f284981ac0553b752c7b8c960ec111bca5deb87af3cf1b921cb70a33e5d326a0c297294bbd9f2cd4971bd6e5ffb840e41714df98804",
      "verificationScript": "21023966fbe8c68f82a05c3a158c06b564a28661dd094f73bc3a9afbb73132562e5eac"
    }
  ],
  "script": "00c10a6d696e74546f6b656e736793544b35a4dd4bcc3b63aa87d3de4bb0850b3606",
  "gas": 0
}`

	tx := neoutils.NeonJSTransaction{}
	json.Unmarshal([]byte(str), &tx)

	b := neoutils.NeonJSTXSerializer(tx)
	log.Printf("%x", b)

}
