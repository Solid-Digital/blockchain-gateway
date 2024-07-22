package ethereum_client_test

import (
	"strings"
	"testing"

	"bitbucket.org/unchain/ethereum2/pkg/domain"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestClient_DeployContract() {
	msg1 := &domain.MessageDeployContract{
		From:     DefaultAccount,
		Solidity: s.helper.StringFromFile("./testdata/contract/no_constructor.sol"),
	}

	msg2 := &domain.MessageDeployContract{
		From:              DefaultAccount,
		Solidity:          s.helper.StringFromFile("./testdata/contract/with_constructor_params.sol"),
		ConstructorParams: map[string]map[string]interface{}{"SimpleStorage": {"_foo": 10, "_bar": "baz"}},
	}

	msg3 := &domain.MessageDeployContract{
		From:              DefaultAccount,
		Solidity:          s.helper.StringFromFile("./testdata/contract/with_constructor_params.sol"),
		ConstructorParams: map[string]map[string]interface{}{"InvalidContractName": {"_foo": 10, "_bar": "baz"}},
	}

	msg4 := &domain.MessageDeployContract{
		From:     DefaultAccount,
		Solidity: s.helper.StringFromFile("./testdata/contract/with_constructor_params.sol"),
	}

	msg5 := &domain.MessageDeployContract{
		From:              DefaultAccount,
		Solidity:          s.helper.StringFromFile("./testdata/contract/with_constructor_params.sol"),
		ConstructorParams: map[string]map[string]interface{}{"SimpleStorage": {"_foo": 10}},
	}

	msg6 := &domain.MessageDeployContract{
		From:     DefaultAccount,
		Solidity: s.helper.StringFromFile("./testdata/contract/does_not_compile.sol"),
	}

	msg7 := &domain.MessageDeployContract{
		From:     DefaultAccount,
		Solidity: s.helper.StringFromFile("./testdata/contract/multiple_contracts.sol"),
	}

	msg8 := &domain.MessageDeployContract{
		From:              DefaultAccount,
		Solidity:          s.helper.StringFromFile("./testdata/contract/multiple_contracts.sol"),
		ConstructorParams: map[string]map[string]interface{}{"Greeter": {"_greeting": "foo"}},
	}

	msg9 := &domain.MessageDeployContract{
		From:     NonRegisteredAccount,
		Solidity: s.helper.StringFromFile("./testdata/contract/no_constructor.sol"),
	}

	cases := map[string]struct {
		Msg       *domain.MessageDeployContract
		Contracts int
		Success   bool
	}{
		"deploy contract without params": {
			msg1,
			1,
			true,
		},
		"deploy contract with params": {
			msg2,
			1,
			true,
		},
		"invalid contract name": {
			msg3,
			0,
			false,
		},
		"missing constructor params": {
			msg4,
			0,
			false,
		},
		"invalid constructor params": {
			msg5,
			0,
			false,
		},
		"does not compile": {
			msg6,
			0,
			false,
		},
		"multiple contracts": {
			msg7,
			2,
			true,
		},
		"multiple contracts with constructor parameters": {
			msg8,
			2,
			true,
		},
		"account not registered in client": {
			msg9,
			0,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.client.DeployContracts(tc.Msg)

			if tc.Success {
				require.NoError(t, err)
				require.NotNil(t, response)
				require.Equal(t, tc.Contracts, len(response))

				for address, deployedContract := range response {
					require.Equal(t, address, strings.ToLower(deployedContract.Address.String()))
					require.True(t, s.helper.ContractAtAddress(s.client, deployedContract.Address))
				}
			} else {
				require.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}
