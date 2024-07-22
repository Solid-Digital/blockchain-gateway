package ethereum_client

import (
	"encoding/json"
	"math/big"
	"strings"

	"bitbucket.org/unchain/ethereum2/pkg/domain"

	"bitbucket.org/unchain/ethereum2/pkg/account"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/pkg/errors"
)

func (c *Client) deploy(account *account.Account, solidity string, params map[string]map[string]interface{}, nonce *big.Int) (domain.DeployContractResponse, error) {
	opts, err := c.deployOptions(account, nonce)
	if err != nil {
		return nil, err
	}

	contracts, err := compileSolidity(solidity)
	if err != nil {
		return nil, err
	}

	ret := domain.DeployContractResponse{}
	for name, contract := range contracts {
		ABI, ABIDefinition, err := generateABI(contract)
		if err != nil {
			return nil, err
		}

		byteCode := common.FromHex(contract.Code)

		inputs, err := constructorInput(name, &ABI.Constructor, params)
		if err != nil {
			return nil, err
		}

		address, tx, _, err := bind.DeployContract(opts, ABI, byteCode, c.Client, inputs...)
		if err != nil {
			return nil, errors.Wrap(err, "failed to deploy contract")
		}

		//Wait for tx to be committed
		err = c.txCommitted(tx.Hash())
		if err != nil {
			return nil, err
		}

		//Without redis we don't do nonce management
		if c.redis != nil {
			//Keep track of latest nonce
			err = c.setNonce(account, tx.Nonce())
			if err != nil {
				//We don't error since the transaction has been sent to the node successfully
				c.logger.Debugf(err.Error())
			}
		}

		ret.AddContract(&address, tx, ABIDefinition)
	}

	return ret, nil
}

func (c *Client) deployOptions(account *account.Account, nonce *big.Int) (*bind.TransactOpts, error) {
	signerFn, err := account.GetSignerFn()
	if err != nil {
		return nil, err
	}

	// TODO: implement gas, gas limit and maybe context?
	return &bind.TransactOpts{
		From:   *account.Address,
		Signer: signerFn,
		Nonce:  nonce,
	}, nil
}

func generateABI(contract *compiler.Contract) (abi.ABI, string, error) {
	ABIDefinition, err := json.Marshal(contract.Info.AbiDefinition)
	if err != nil {
		return abi.ABI{}, "", errors.Wrap(err, "failed to read ABI")
	}

	ABI, err := abi.JSON(strings.NewReader(string(ABIDefinition)))
	if err != nil {
		return abi.ABI{}, "", errors.Wrap(err, "failed to generate ABI")
	}

	return ABI, string(ABIDefinition), nil
}

// compileSolidity returns a map of contract names to compiled contracts
func compileSolidity(solidity string) (map[string]*compiler.Contract, error) {
	// TODO: use of different compilers? Not sure how this works right now.
	// If solc equals "" it will refer to solc binary on path
	solc := ""
	contracts, err := compiler.CompileSolidityString(solc, solidity)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile solidity source")
	}

	// For some reason the contract names are prefixed with "<stdin>:", we don't
	// want this so we create a new map and return that instead
	cleanContracts := make(map[string]*compiler.Contract, len(contracts))
	for key, value := range contracts {
		cleanContracts[cleanKey(key)] = value
	}

	return cleanContracts, nil
}

func cleanKey(key string) string {
	if len(key) >= 8 && key[:8] == "<stdin>:" {
		return key[8:]
	}

	return key
}
