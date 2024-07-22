package ethereum_client_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/ethereum/go-ethereum/core/types"

	"bitbucket.org/unchain/ethereum2/pkg/domain"
	"github.com/stretchr/testify/require"
)

// This test experiences random failures, seems to be related to the following:
// https://github.com/trufflesuite/ganache-core/issues/166
// Seems the issue has been resolved by using v6.6.0 of Ganache-core
func (s *TestSuite) TestClient_ContractTransactionFunction() {
	solidity1 := s.helper.StringFromFile("./testdata/contract/transact_function_no_params.sol")
	address1, _ := s.helper.DeploySingleContract(s.client, solidity1, nil)
	msg1 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       address1,
		Function: "increment",
	}

	solidity2 := s.helper.StringFromFile("./testdata/contract/no_constructor.sol")
	address2, _ := s.helper.DeploySingleContract(s.client, solidity2, nil)
	msg2 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       address2,
		Function: "set",
		Params:   map[string]interface{}{"x": 12},
	}

	solidity3 := s.helper.StringFromFile("./testdata/contract/transact_function_no_params.sol")
	address3, _ := s.helper.DeploySingleContract(s.client, solidity3, nil)
	msg3 := &domain.MessageCallContractFunction{
		From:     NonRegisteredAccount,
		To:       address3,
		Function: "increment",
	}

	msg4 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       NonExistingContractAddress,
		Function: "increment",
	}

	solidity5 := s.helper.StringFromFile("./testdata/contract/transact_function_no_params.sol")
	address5, _ := s.helper.DeploySingleContract(s.client, solidity5, nil)
	msg5 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       address5,
		Function: "invalidFunction",
	}

	solidity6 := s.helper.StringFromFile("./testdata/contract/no_constructor.sol")
	address6, _ := s.helper.DeploySingleContract(s.client, solidity6, nil)
	msg6 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       address6,
		Function: "set",
	}

	solidity7 := s.helper.StringFromFile("./testdata/contract/no_constructor.sol")
	address7, _ := s.helper.DeploySingleContract(s.client, solidity7, nil)
	msg7 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       address7,
		Function: "set",
		Params:   map[string]interface{}{"invalidParam": 12},
	}

	solidity8 := s.helper.StringFromFile("./testdata/contract/no_constructor.sol")
	address8, _ := s.helper.DeploySingleContract(s.client, solidity8, nil)
	msg8 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       address8,
		Function: "set",
		Params:   map[string]interface{}{"x": "invalid type"},
	}

	cases := map[string]struct {
		Msg     *domain.MessageCallContractFunction
		Success bool
	}{
		"call function without params": {
			msg1,
			true,
		},
		"call function with params": {
			msg2,
			true,
		},
		"account not registered in client": {
			msg3,
			false,
		},
		"contract does not exist": {
			msg4,
			false,
		},
		"function does not exist": {
			msg5,
			false,
		},
		"missing function params": {
			msg6,
			false,
		},
		"invalid function params": {
			msg7,
			false,
		},
		"invalid function params type": {
			msg8,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.client.CallContractFunction(tc.Msg)

			if tc.Success {
				require.NoError(t, err)
				require.True(t, s.helper.TransactionCommitted(s.client, response.(*types.Transaction)))
			} else {
				require.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}

//Verify that many transactions can be executed in a short period of time (i.e. in the same block)
func (s *TestSuite) TestNonceManagement() {
	solidity := s.helper.StringFromFile("./testdata/contract/nonce_management.sol")
	address, _ := s.helper.DeploySingleContract(s.blockMiningNetworkClient, solidity, nil)

	//Make sure the contract has been committed to the ledger
	time.Sleep(6 * time.Second)

	expected := map[string]string{}
	for i := 0; i < 10; i++ {
		k := fmt.Sprintf("%d", i)
		v := randomdata.City()

		msg := &domain.MessageCallContractFunction{
			From:     DefaultAccount,
			To:       address,
			Function: "set",
			Params:   map[string]interface{}{"_key": k, "_value": v},
		}

		response, err := s.blockMiningNetworkClient.CallContractFunction(msg)
		s.Require().NoError(err)

		_, ok := response.(*types.Transaction)
		s.Require().True(ok)

		expected[k] = v
	}

	//Make sure the transactions are committed
	time.Sleep(6 * time.Second)

	for k, v := range expected {
		msg := &domain.MessageCallContractFunction{
			From:     DefaultAccount,
			To:       address,
			Function: "get",
			Params:   map[string]interface{}{"_key": k},
		}

		response, err := s.blockMiningNetworkClient.CallContractFunction(msg)
		s.Require().NoError(err)

		s.Require().Equal(v, response)
	}
}

//Verify that transactions are executed in synchronous mode
//This test is similar to the nonce management test, but does not have
//the sleep instructions in between, yet runs without problems because of sync mode.
func (s *TestSuite) TestSynchronousMode() {
	solidity := s.helper.StringFromFile("./testdata/contract/nonce_management.sol")
	address, _ := s.helper.DeploySingleContract(s.syncModeClient, solidity, nil)

	for i := 0; i < 5; i++ {
		k := fmt.Sprintf("%d", i)
		v := randomdata.City()

		msg := &domain.MessageCallContractFunction{
			From:     DefaultAccount,
			To:       address,
			Function: "set",
			Params:   map[string]interface{}{"_key": k, "_value": v},
		}

		response, err := s.syncModeClient.CallContractFunction(msg)
		s.Require().NoError(err)

		tx, ok := response.(*types.Transaction)
		s.Require().True(ok)
		s.helper.TransactionCommitted(s.syncModeClient, tx)
	}
}
