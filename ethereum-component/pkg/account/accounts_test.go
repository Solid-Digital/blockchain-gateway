package account_test

import (
	"testing"

	"bitbucket.org/unchain/ethereum2/pkg/account"
	"github.com/stretchr/testify/require"
)

func (s *TestSuite) TestNewAccounts() {
	cfg := s.helper.BytesFromFile("./testdata/config/accounts.toml")
	accountCfgs := s.factory.AccountCfgs(cfg)
	accounts, err := account.NewAccounts(s.logger, accountCfgs)

	s.Require().NoError(err)
	s.Require().Equal(3, len(accounts))
}

func (s *TestSuite) TestAccounts_GetAccount() {
	cases := map[string]struct {
		Accounts account.Accounts
		Address  string
		Success  bool
	}{
		"get first account": {
			s.factory.Accounts(s.helper.BytesFromFile("./testdata/config/accounts.toml")),
			"0xfdfa8d41f986c80904bf4825402e788f3121e7af",
			true,
		},
		"get second account": {
			s.factory.Accounts(s.helper.BytesFromFile("./testdata/config/accounts.toml")),
			"0x69091d42c8307d9a24a47b2d92d4506604ae44b9",
			true,
		},
		"get third account": {
			s.factory.Accounts(s.helper.BytesFromFile("./testdata/config/accounts.toml")),
			"0xfe6cca44dec3726aff6d44c39974e821f6d7510e",
			true,
		},
		"account does not exist": {
			s.factory.Accounts(s.helper.BytesFromFile("./testdata/config/accounts.toml")),
			"0xc81630116d2335f39732dc18bb20500b33ee3984",
			false,
		},
		"single account does not require address to get it": {
			s.factory.Accounts(s.helper.BytesFromFile("./testdata/config/single_account.toml")),
			"",
			true,
		},
		"multiple accounts require address to get it": {
			s.factory.Accounts(s.helper.BytesFromFile("./testdata/config/accounts.toml")),
			"",
			false,
		},
		"private key with leading 0x does not error": {
			s.factory.Accounts(s.helper.BytesFromFile("./testdata/config/leading_0x_account.toml")),
			"",
			true,
		},
	}

	for name, tc := range cases {
		s.T().Run(name, func(t *testing.T) {
			account, err := tc.Accounts.GetAccount(tc.Address)

			if tc.Success {
				require.NoError(t, err)
				require.NotNil(t, account)
			} else {
				require.Error(t, err)
				require.Nil(t, account)
			}
		})
	}
}
