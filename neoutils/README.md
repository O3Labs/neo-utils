# NEO utilities

This package contains useful functions that make your life easier when working with NEO blockchain.

#### What you can use this package for
#### Wallet

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
- Reverse bytes
- Hex string to bytes
- Bytes to hex string
- Convert script hash to NEO address
- Convert NEO Address to script hash
- Validate NEO Address
- Convert Byte array to big int
- Parse NEP9 URI


#### NEO JSON RPC
- Get contract state with smart contract's script hash
- Send raw transaction
- Get raw transaction with TXID

#### City of Zion APIs
- Get unspent data by NEO Address

#### NEO Smart contract
- Generate invocation script data
- Generate invocation inputs data
- Generate invocation output data
- Generate transaction attributes data
- Generate invocation and verification script with signatures
- Parse raw transaction's script to operation name and args
- Get invoked operation from raw transaction's script data
- Get invoked smart contract script hash from raw transaction's script data
- Get invoked params from raw transaction's script data
- Generate ready-for-sendrawtransaction smart contract invocation data
- Invoke NEO Smart contract
