# NEO utilities

This package contains useful functions that make your life easier when working with NEO blockchain.

#### What you can use this package for
##### Wallet
- Create a new wallet address for NEO blockchain
- Restore a wallet with WIF
- Restore a wallet with raw private key

#### Encryption
- Sign data using ECDSA
- Encrypt data using AES
- Decrypt AES encrypted data 
- Public key encryption using ECDH
- Create N-parts shared secret using [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing)
- Restore data from shared secret using [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing)

#### NEO Nodes utilities
- Select NEO best node by measuring the latency between caller and the nodes concurrently (this is really fast!)

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
