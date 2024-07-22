# Ethereum Trigger

The Ethereum trigger allows you to take an action based on an event that comes
from an Ethereum network.

**In order to configure events you need the following information:**

A connection to the Ethereum node, because you need to receive events from the
node this cannot be http:// or https://.

```toml
[ethereum]
host =          # websocket (ws://) or secure websocket (wss://) url of the Ethereum node
```

One or multiple contracts on which you want to listen for events.

```toml
[[contracts]]
address =       # The address at which the contract is deployed
abi =           # the ABI for this contract.
```

One or multiple functions to listen to

```toml
[[events]]
contractAddress =   # The address of the contact (same as under contracts)
name =              # The name of the event
filters =           # Array of items to filter on.
```

Filters are useful when you want to listen to events from a given contract, but
for example only want a notification when the event was created by one of the
accounts you are interested in.

**This component will output the following fields**
