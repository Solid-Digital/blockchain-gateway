# Readme

Ethereum2 implements the Janus V2 interfaces.

## Supported features

The Ethereum connector has a trigger component and an action component. The trigger component is used to listen for 
events, the action component is used for all other features.

### Features Action component

* call smart contract function
	* call function (i.e. non-payable function)
	* transaction (i.e. payable function)
* deploy smart contract
* nonce management
* synchronous mode
* gas price
* gas limit

Configuration of Redis is required to enable nonce management. By using nonce management, 
multiple transactions of the same account can be mined into to same block. It is also possible
to set the nonce manually for each message, in case you want to do nonce management
in some external application.

When you enable synchronous mode, execution of the pipeline will be paused until the transaction
has been committed. When the transaction is still pending after 120 seconds, the pipeline
will terminate execution. Note that synchronous mode is not suitable for networks with a high
blocktime like the Ethereum mainnet.

If the gas limit is not specified, it will be estimated by locally executing the transaction.

If the gas price is not specified, a gas price oracle will be used.

### Features Trigger component

* listen to events

## Testing

Ganache-cli is used for testing: 

https://github.com/trufflesuite/ganache-cli

Run `make test` to run the tests.

## Configuration of Action component

| item      | description                                  |
|-----------|----------------------------------------------|
| host      | address and port of Ethereum node            |
| accounts  | Ethereum accounts                            |
| contracts | Ethereum smart contracts to interact with    |    
| redis     | optional redis database for nonce management |

### Example configuration 

```toml
[ethereum]
host = "http://example.com:8545"
synchronousMode = false
gasPrice = 0
gasLimit = 100000000

[[accounts]]
address = "0xfdfa8d41f986c80904bf4825402e788f3121e7af"
privateKey = "0x5b3a720780cdf4288877eff96b5e4be01f5bddf4fb77f304e12f5a00f17220cc"

[[accounts]]
address = "0x69091d42c8307d9a24a47b2d92d4506604ae44b9"
privateKey = "0x95206d4239e755291d12b7bbc778b388fb765317cd115f8a72e27827629edc72"

[[accounts]]
address = "0xfe6cca44dec3726aff6d44c39974e821f6d7510e"
privateKey = "0x11cce58ced3fce5e69f475418e346c516abc59c32104a398c6e2c409949d30e7"

[[contracts]]
address = "0x5d7de81c4a2009e60d4fea85fe88fdc659ae5cad"
abi = '[{"constant":true,"inputs":[{"name":"s","type":"string"}],"name":"getString","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"pure","type":"function"}]'

[[contracts]]
address = "0xc67206229c1f20ce80f3e13fbc9199e329dea1e1"
abi = '[{"constant":false,"inputs":[{"name":"s","type":"string"}],"name":"setString","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]'

[[redis]]
host = "localhost:6379"
password = "mySecretPassword"
db = 5
```

## Inputs

| item              | description                                                    |
|-------------------|----------------------------------------------------------------|
| from              | address (string) of account you want to use                    |
| to                | address (string) of smart contract you want to interact with   |
| function          | name of function (string) you want to call                     |
| params            | arguments (map of string to bytes)                             |
| nonce             | optional transaction nonce                                     |
| solidity          | solidity source code (string) of contract you want to deploy   |
| constructorParams | arguments (map of map to bytes)                                |
| type              | request type (i.e. callContractFunction or deployContract      |

## Outputs

| key                       | description                                        |
|---------------------------|----------------------------------------------------|
| contractCallOutput        | return values of call function                     |
| contractTransactionOutput | transaction receipt from transaction function      |
| deployContractOutput      | map of deployed contract address (string) to bytes |  

### Call function

When calling a function on a smart contract, the following inputs are used:

* from
* to
* function
* params (optional)

The return values of the contract are outputted to *contractCallOutput*.

### Transaction

When creating a transaction, the following inputs are used:

* from
* to
* function
* params (optional)
* nonce (optional)

The transaction receipt is outputted to *contractTransactionOutput*.

### Deploy contract

When deploying a contract, the following inputs are used:

* from (optional when a single account has been configured)
* solidity
* constructorParams
* nonce (optional)

## Configuration of Trigger component

| item      | description                               |
|-----------|-------------------------------------------|
| host      | address and port of Ethereum node         |
| contracts | Ethereum smart contracts to interact with |
| events    | Ethereum events to listen for             |

### Example configuration 

```toml
[ethereum]
	host = "ws://localhost:8545"

[[contracts]]
	address = "0x5d7de81c4a2009e60d4fea85fe88fdc659ae5cad"
	abi = '[{"anonymous":false,"inputs":[{"indexed":true,"name":"_from","type":"address"},{"indexed":true,"name":"_to","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Transfer","type":"event"}]'

[[contracts]]
	address = "0xc67206229c1f20ce80f3e13fbc9199e329dea1e1"
	abi = '[{"anonymous":false,"inputs":[{"indexed":true,"name":"_from","type":"address"},{"indexed":true,"name":"_to","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Transfer","type":"event"}]'

[[events]]
    contractAddress = "0x5d7de81c4a2009e60d4fea85fe88fdc659ae5cad"
    name = "Transfer"
    filters = []

[[events]]
    contractAddress = "0xc67206229c1f20ce80f3e13fbc9199e329dea1e1"
    name = "Transfer"
    filters = [
        ["0x57b675c9d3751e94bec2d943e6f47b2dfba73a2d", "0x507d5f3241a9d9cc618473cd706670b45ca3a9ac"]
    ]
```

## Outputs

| key    | description                                        |
|--------|----------------------------------------------------|
| event  | name of the event                                  |
| values | values emitted by the event                        |
| log    | log of the event                                   |  
