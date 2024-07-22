package ethereum_client

import "github.com/ethereum/go-ethereum/accounts/abi"

// Call functions are non-payable contract calls, payable contract calls are referred to as Transaction.
const (
	Call        = "call"
	Transaction = "transaction"
)

func (c *Client) functionType(function *abi.Method) string {
	if function.Const {
		return Call
	}

	return Transaction
}
