package contract_test

import (
	"testing"

	"bitbucket.org/unchain/ethereum2/pkg/contract"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestContract_GetFunction() {
	contracts := s.factory.Contracts(s.helper.BytesFromFile("./testdata/config/single_contract.toml"))
	singleContract := contracts[s.helper.KeyFromMap(contracts)]

	cases := map[string]struct {
		Contract     *contract.Contract
		FunctionName string
		Success      bool
	}{
		"get function": {
			singleContract,
			"get",
			true,
		},
		"function does not exist": {
			singleContract,
			"set",
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := tc.Contract.GetFunction(tc.FunctionName)

			if tc.Success {
				require.NoError(t, err)
				require.NotNil(t, response)
				require.Equal(t, tc.FunctionName, response.Name)
			} else {
				require.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}

func (s *TestSuite) TestContract_GetEvent() {
	contracts := s.factory.Contracts(s.helper.BytesFromFile("./testdata/config/contract_with_event.toml"))
	contractWithEvent := contracts[s.helper.KeyFromMap(contracts)]

	cases := map[string]struct {
		Contract  *contract.Contract
		EventName string
		Success   bool
	}{
		"get event": {
			contractWithEvent,
			"Transfer",
			true,
		},
		"event does not exist": {
			contractWithEvent,
			"InvalidEvent",
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			response, err := tc.Contract.GetEvent(tc.EventName)

			if tc.Success {
				require.NoError(t, err)
				require.NotNil(t, response)
				require.Equal(t, tc.EventName, response.Name)
			} else {
				require.Error(t, err)
				require.Nil(t, response)
			}
		})
	}
}

func (s *TestSuite) TestContract_AddressString() {
	contracts := s.factory.Contracts(s.helper.BytesFromFile("./testdata/config/single_contract.toml"))
	singleContract := contracts[s.helper.KeyFromMap(contracts)]

	// Contains upper case characters
	s.Require().Equal("0x76BC9E61A1904B82CbF70d1fd9C0f8a120483BbB", singleContract.Address.String())

	// All lower case
	s.Require().Equal("0x76bc9e61a1904b82cbf70d1fd9c0f8a120483bbb", singleContract.AddressString())
}
