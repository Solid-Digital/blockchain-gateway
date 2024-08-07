package ethereum_client

import (
	"bitbucket.org/unchain/ethereum2/pkg/contract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Client) boundContract(contract *contract.Contract) *bind.BoundContract {
	return bind.NewBoundContract(*contract.Address, *contract.ABI, c.Client, c.Client, c.Client)
}
