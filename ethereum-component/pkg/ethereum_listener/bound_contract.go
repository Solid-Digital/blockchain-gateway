package ethereum_listener

import (
	"bitbucket.org/unchain/ethereum2/pkg/contract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Listener) boundContract(contract *contract.Contract) *bind.BoundContract {
	return bind.NewBoundContract(*contract.Address, *contract.ABI, c.client, c.client, c.client)
}
