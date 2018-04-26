# NEO utilities

This package contains useful functions that make your life easier when working with NEO blockchain.

### List of functions available in this library

#### Wallet
```go
import "github.com/o3labs/neo-utils/neoutils"
```

##### Wallet struct
```go
type Wallet struct {
	PublicKey       []byte
	PrivateKey      []byte
	Address         string
	WIF             string
	HashedSignature []byte
}

```

##### Create a new wallet address for NEO blockchain
```go
neoutils.NewWallet() (*Wallet, error)
```
##### Restore a wallet from WIF
```go
neoutils.GenerateFromWIF(wif string) (*Wallet, error)
```

##### Restore a wallet from raw private key string
```go
neoutils.GenerateFromPrivateKey(privateKey string) (*Wallet, error)
```
---


#### Encryption
```go
import "github.com/o3labs/neo-utils/neoutils"
```
##### Sign data using ECDSA
```go
neoutils.Sign(data []byte, key string) ([]byte, error) 
```
##### Encrypt data using AES
```go
neoutils.Encrypt(key []byte, text string) string 
```

##### Decrypt AES encrypted data 
```go
neoutils.Decrypt(key []byte, encryptedText string) string
```

##### Public key encryption using ECDH
```go
(w *Wallet) ComputeSharedSecret(publicKey []byte) []byte
```

##### Create N-parts shared secret using [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing)
```go
neoutils.GenerateShamirSharedSecret(secret string) (*SharedSecret, error)
```

##### Restore data from shared secret using [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing)
```go
neoutils.RecoverFromSharedSecret(first []byte, second []byte) (string, error)
```

#### NEO Nodes utilities
```go
import "github.com/o3labs/neo-utils/neoutils"
```
##### Select best node
Select NEO best node by measuring the latency between caller and the nodes concurrently.
```go
type SeedNodeResponse struct {
	URL          string
	BlockCount   int
	ResponseTime int64 //milliseconds
}

neoutils.SelectBestSeedNode(commaSeparatedURLs string) *SeedNodeResponse 
```
--- 
#### Utilities methods
```go
import "github.com/o3labs/neo-utils/neoutils"
```
##### Reverse bytes
```go
neoutils.ReverseBytes(b []byte) []byte 
```
##### Hex string to bytes
```go
neoutils.HexTobytes(hexstring string) (b []byte)
```
##### Bytes to hex string
```go
neoutils.BytesToHex(b []byte) string
```

##### Convert script hash to NEO address
```go
neoutils.ScriptHashToNEOAddress(scriptHash string) string
```
##### Convert NEO Address to script hash
```go
neoutils.NEOAddressToScriptHash(neoAddress string) string 
```
##### Validate NEO Address
```go
neoutils.ValidateNEOAddress(address string) bool 
```
##### Convert Byte array to big int
```go
neoutils.ConvertByteArrayToBigInt(hexString string) *big.Int
```
##### Parse NEP9 URI
```go
type SimplifiedNEP9 struct {
	To      string  `json:"to"`
	AssetID string  `json:"assetID"`
	Amount  float64 `json:"amount"`
}
neoutils.ParseNEP9URI(uri string) (*SimplifiedNEP9, error) 
```

---

#### NEO JSON RPC
```go
import "github.com/o3labs/neo-utils/neoutils/neorpc"
```

##### Get contract state with smart contract's script hash
```go
client := neorpc.NewClient("http://localhost:30333")
if client == nil {
	return
}

result := client.GetContractState("ce575ae1bb6153330d20c560acb434dc5755241b")
```

##### Send raw transaction
```go
client := neorpc.NewClient("http://localhost:30333")
if client == nil {
	return
}
raw := ""
result := client.SendRawTransaction(raw)
```

##### Get raw transaction with TXID
```go
client := neorpc.NewClient("http://localhost:30333")
if client == nil {
	return
}
txID := "bde02f8c6482e23d5b465259e3e438f0acacaba2a7a938d5eecd90bba0e9d1ad"
result := client.GetRawTransaction(txID)
```
---

#### City of Zion APIs
```go
import "github.com/o3labs/neo-utils/neoutils/coz"
```
##### Get unspent data by NEO Address
```go
client := coz.NewClient("http://127.0.0.1:5000")
unspent, err := client.GetUnspentByAddress("AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y")
for _, v := range unspent.GAS.Unspent {
	log.Printf("%+v", v)
}	
```
---

#### NEO Smart contract
```go
import "github.com/o3labs/neo-utils/neoutils/smartcontract"
```

##### Generate invocation script data
```go
smartcontract.GenerateContractInvocationData(scriptHash ScriptHash, operation string, args []interface{}) []byte
```
##### Generate invocation inputs data
```go
smartconract.GenerateTransactionInput(unspent Unspent, assetToSend NativeAsset, amountToSend float64) ([]byte, error)
```

##### Generate invocation output data
```go
smartcontract.GenerateTransactionOutput(sender NEOAddress, receiver NEOAddress, unspent Unspent, assetToSend NativeAsset, amountToSend float64) ([]byte, error)
```
##### Generate transaction attributes data
```go
smartcontract.GenerateTransactionAttributes(attributes map[TransactionAttribute][]byte) ([]byte, error)
```
##### Generate invocation and verification script with signatures
```go
smartcontract.GenerateVerificationScripts(signatures []TransactionSignature) []byte
```
##### Parse raw transaction's script to operation name and args
```go
p := smartcontract.NewParserWithScript("51143acefb110cba488ae0d809f5837b0ac9c895405e52c10c6d696e74546f6b656e73546f67b17f078543788c588ce9e75544e325a050f8c1b7")

//you can define the known method signature here
type methodSignature struct {
	Operation smartcontract.Operation  //operation
	To        smartcontract.NEOAddress //args[0]
	Amount    int                      //args[1]
}
m := methodSignature{}
list, err := p.Parse(&m)
if err != nil {
	return
}
for _, v := range list {
	log.Printf("%+v", v.(*methodSignature))
}
```

##### Parse raw transaction's script that contains multiple appcall
```go
script := `0500bca06501145a936d7abbaae28579dd36609f910f9b50de972f147bee835ff211327677c453d5f19b693e70a361ab53c1087472616e7366657267b6155db85e53298f01e0280cc2f21a0f40c4e808f10400e1f505147e548ecd2a87dd58731e6171752b1aa11494c62f147bee835ff211327677c453d5f19b693e70a361ab53c1087472616e7366657267b6155db85e53298f01e0280cc2f21a0f40c4e808f10500dc5c240214c10704464fade3197739536450ec9531a1f24a37147bee835ff211327677c453d5f19b693e70a361ab53c1087472616e7366657267b6155db85e53298f01e0280cc2f21a0f40c4e808f166b2263911344b5b15`
p := smartcontract.NewParserWithScript(script)
type methodSignature struct {
	Operation smartcontract.Operation  //operation
	From      smartcontract.NEOAddress //args[0]
	To        smartcontract.NEOAddress //args[1]
	Amount    int                      //args[2]
}
m := methodSignature{}
list, err := p.Parse(&m)
if err != nil {
	return
}

for _, v := range list {
	log.Printf("%+v", v.(*methodSignature))
}
```

##### Get invoked operations from raw transaction's script data
```go
p := smartcontract.NewParserWithScript("51143acefb110cba488ae0d809f5837b0ac9c895405e52c10c6d696e74546f6b656e73546f67b17f078543788c588ce9e75544e325a050f8c1b7")
result, err := p.GetListOfOperations()
if err != nil {
	return
}

log.Printf("result = %v", result)
```

##### Get invoked smart contract script hashes from raw transaction's script data
```go
p := smartcontract.NewParserWithScript("51143acefb110cba488ae0d809f5837b0ac9c895405e52c10c6d696e74546f6b656e73546f67b17f078543788c588ce9e75544e325a050f8c1b7")
result, err := p.GetListOfScriptHashes()
if err != nil {
	return
}

log.Printf("result = %v", result)
```

##### Generate ready-for-sendrawtransaction smart contract invocation data
```go
var validSmartContract = neoutils.UseSmartContract("b7c1f850a025e34455e7e98c588c784385077fb1")
validSmartContract.GenerateInvokeFunctionRawTransaction(wallet Wallet, unspent smartcontract.Unspent, attributes map[smartcontract.TransactionAttribute][]byte, operation string, args []interface{}) ([]byte, error)
```
