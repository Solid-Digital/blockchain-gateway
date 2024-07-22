# Ethereum Action

The Ethereum action allows you to create transactions onto an Ethereum network.

## Required initialization configuration

This configuration is used to configure the component, and will be used when the
pipeline is started. The configuration format used is
[TOML](https://learnxinyminutes.com/docs/toml/)

**To connect to an Ethereum node we need the following information:**

```toml
# A node to connect to
[ethereum]
host =         # http(s) url and port of Ethereum node

# When you enable synchronous mode, execution of the pipeline will be paused until the transaction
# has been committed. When the transaction is still pending after 120 seconds, the pipeline
# will terminate execution. Note that synchronous mode is not suitable for networks with a high
# blocktime like the Ethereum mainnet.
synchronousMode =

# If the gas limit is not specified, it will be estimated by locally executing the transaction.  
gasLimit = 

# If the gas price is not specified, a gas price oracle will be used.
gasPrice = 

# One or multiple accounts
[[accounts]]
address =      # account address / private key
privateKey =   # account private key

# Optional redis configuration in case you want to use nonce management
[redis]
host =         # host address and port  
password =     # optional password
db = 0         # optional db
```

You may use environment variables from your pipeline environment in this
configuration, this is useful for things like URLs, secrets and passwords. For
example \$PRIVATE_KEY.

**If you are calling an existing contract (the common case) you will need to
specify the following details:**

```toml
# A list of contracts to call.
[[contracts]]
address =      # address of contract
abi =          # The abi, or solidity contract definition (JSON string)
```

Note that you may specify the abi as a multiline string like so:

```toml
abi = '''
'''
```

## Input schema

Then, each time this component should take an action, the following parameters
may be specified:

| Key               | Description                                                 | Required              |
| ----------------- | ----------------------------------------------------------- | ----------------------|
| type              | Action to take ('callContractFunction' or 'deployContract') | true                  |
| from              | Address of sender (specified above)                         | true                  |
| to                | Address of contract                                         | true                  |
| function          | Smart contract function to call                             | for call              |
| params            | The parameters for the contract (in JSON)                   | for call (optional)   |
| solidity          | The solidity code                                           | for deploy            |
| constructorParams | Parameters to create the contract with                      | for deploy (optional) |

In the table below you may find some examples of what this could look like:

| Key               | Example (string)                              | Example (from state object)            |
| ----------------- | --------------------------------------------- | -------------------------------------- |
| type              | callContractFunction                          | \$.http-trigger.body.type              |
| from              | 0xb3622a6ad922f8f82690f31a8da358fc55100bc4    | \$.http-trigger.body.from              |
| to                | 0x52133A2eb5Fffdf3e87741B245997b359462273e    | \$.http-trigger.body.to                |
| function          | set                                           | \$.http-trigger.body.function          |
| params            | {"x": 100}                                    | \$.http-trigger.body.params            |
| solidity          | pragma solidity ^0.5.9; contract ...          | \$.http-trigger.body.solidity          |
| constructorParams | { ERC20Token": { "\_name": "Unchain token" }} | \$.http-trigger.body.constructorParams |

## Output schema

| Key                       | Description                                                                    |
| ------------------------- | ------------------------------------------------------------------------------ |
| contractCallOutput        | Single value or array of responses                                             |
| contractTransactionOutput | Object which contains transactionReceipt and some other info                   |
| deployContractOutput      | Object which contains ABI and bytecode of the deployed contract and other info |

<!-- outputSchema = ["", "contractTransactionOutput", "deployContractOutput"] -->

## Usage

This component will use the details you specify to connect to an Ethereum node.
Transactions are signed internally, so your private key will never leave the 
pipeline (server).

If you are making any transaction (which includes deploying a contract) you may
need to have some Ether available on your account.

While some private networks not require this because their gas price is set to
0, public networks will require you to pay some transaction fee in Ether.

### Deploying a contract

To deploy a new contract you will need to specify a valid account, When the
requirements are met, deploying is then as simple as calling this component with
`"type": "deployContract"` and `"solidity": <the contract source code>` It is
possible that the contract requires constructor parameters. You may specify
these like so `"constructorParams": '{"your": "params"}'`

The component will output its response to `deployContractOutput`

### Calling a contract

There are two types of calls you can make: get information (which is free), and
put information into the blockchain (which costs Ether). Both are called in the
same way, like so:

`"type": "callContract"`, `"function": "functionName"`,
`"params": {"x": "foo", "y": "bar"}`

But they will have a different response.

**Free (get) calls** will respond with the retrieved value on the
`contractCallOutput`

**Paid (transaction) calls** cannot happen immediately, because the transaction needs
to be mined onto a block.

Instead, the blockchain, and this component responds with a transactionReceipt
on the `contractTransactionOutput`. It contains a lot of information about the
transaction.
