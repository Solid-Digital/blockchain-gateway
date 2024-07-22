package ethereum_client_test

import (
	"math/big"
	"testing"

	"bitbucket.org/unchain/ethereum2/pkg/domain"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestClient_ContractCallFunction() {
	solidity1 := s.helper.StringFromFile("./testdata/contract/no_constructor.sol")
	address1, _ := s.helper.DeploySingleContract(s.client, solidity1, nil)

	msg1 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       address1,
		Function: "get",
	}

	solidity2 := s.helper.StringFromFile("./testdata/contract/with_function_params.sol")
	address2, _ := s.helper.DeploySingleContract(s.client, solidity2, nil)
	msg2 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       address2,
		Function: "get",
		Params: map[string]interface{}{
			"_myString": "my string",
		},
	}

	solidity3 := s.helper.StringFromFile("./testdata/contract/multiple_return_values.sol")
	address3, _ := s.helper.DeploySingleContract(s.client, solidity3, nil)
	msg3 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       address3,
		Function: "get",
		Params: map[string]interface{}{
			"_myString": "my string",
		},
	}

	solidity4 := s.helper.StringFromFile("./testdata/contract/no_constructor.sol")
	address4, _ := s.helper.DeploySingleContract(s.client, solidity4, nil)
	msg4 := &domain.MessageCallContractFunction{
		From:     NonRegisteredAccount,
		To:       address4,
		Function: "get",
	}

	msg5 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       "0xd0a6e6c54dbc68db5db3a091b171a77407ff7ccf",
		Function: "get",
	}

	solidity6 := s.helper.StringFromFile("./testdata/contract/no_constructor.sol")
	address6, _ := s.helper.DeploySingleContract(s.client, solidity6, nil)
	msg6 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       address6,
		Function: "invalidFunction",
	}

	solidity7 := s.helper.StringFromFile("./testdata/contract/with_function_params.sol")
	address7, _ := s.helper.DeploySingleContract(s.client, solidity7, nil)
	msg7 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       address7,
		Function: "get",
	}

	solidity8 := s.helper.StringFromFile("./testdata/contract/with_function_params.sol")
	address8, _ := s.helper.DeploySingleContract(s.client, solidity8, nil)
	msg8 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       address8,
		Function: "get",
		Params: map[string]interface{}{
			"_notValid": "my string",
		},
	}

	solidity9 := s.helper.StringFromFile("./testdata/contract/with_function_params.sol")
	address9, _ := s.helper.DeploySingleContract(s.client, solidity9, nil)
	msg9 := &domain.MessageCallContractFunction{
		From:     DefaultAccount,
		To:       address9,
		Function: "get",
		Params: map[string]interface{}{
			"_myString": 12,
		},
	}

	cases := map[string]struct {
		Msg      *domain.MessageCallContractFunction
		Expected interface{}
		Success  bool
	}{
		"call function without params": {
			msg1,
			big.NewInt(int64(123)),
			true,
		},
		"call function with params": {
			msg2,
			"my string",
			true,
		},
		"multiple return values": {
			msg3,
			[]interface{}{
				s.helper.StringPtr("my string"),
				s.helper.Uint8Ptr(10)},
			true,
		},
		"account not registered in client": {
			msg4,
			nil,
			false,
		},
		"contract does not exist": {
			msg5,
			nil,
			false,
		},
		"function does not exist": {
			msg6,
			nil,
			false,
		},
		"missing function params": {
			msg7,
			nil,
			false,
		},
		"invalid function params": {
			msg8,
			nil,
			false,
		},
		"invalid function params type": {
			msg9,
			nil,
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := s.client.CallContractFunction(tc.Msg)

			if tc.Success {
				require.NoError(t, err)
				require.Equal(t, tc.Expected, response)
			} else {
				require.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}
