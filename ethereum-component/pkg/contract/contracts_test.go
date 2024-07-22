package contract_test

import (
	"testing"

	"bitbucket.org/unchain/ethereum2/pkg/contract"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestNewContracts() {
	contractCfgs := s.factory.ContractCfgs(s.helper.BytesFromFile("./testdata/config/contracts.toml"))
	contracts, err := contract.NewContracts(s.logger, contractCfgs)

	s.Require().NoError(err)
	s.Require().Equal(2, len(contracts))
}

func (s *TestSuite) TestContracts_GetContract() {
	cases := map[string]struct {
		Contracts contract.Contracts
		Address   string
		Success   bool
	}{
		"get first contract": {
			s.factory.Contracts(s.helper.BytesFromFile("./testdata/config/contracts.toml")),
			"0x5d7de81c4a2009e60d4fea85fe88fdc659ae5cad",
			true,
		},
		"get second contract": {
			s.factory.Contracts(s.helper.BytesFromFile("./testdata/config/contracts.toml")),
			"0x76bc9e61a1904b82cbf70d1fd9c0f8a120483bbb",
			true,
		},
		"contract does not exist": {
			s.factory.Contracts(s.helper.BytesFromFile("./testdata/config/contracts.toml")),
			"0xd0a6e6c54dbc68db5db3a091b171a77407ff7ccf",
			false,
		},
		"single contract does not require address to get it": {
			s.factory.Contracts(s.helper.BytesFromFile("./testdata/config/single_contract.toml")),
			"",
			true,
		},
		"multiple contracts require address to get it": {
			s.factory.Contracts(s.helper.BytesFromFile("./testdata/config/contracts.toml")),
			"",
			false,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			contract, err := tc.Contracts.GetContract(tc.Address)

			if tc.Success {
				require.NoError(t, err)
				require.NotNil(t, contract)
			} else {
				require.Error(t, err)
				require.Nil(t, contract)
			}
		})
	}
}
