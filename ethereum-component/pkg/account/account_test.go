package account_test

// we only test the happy path here
func (s *TestSuite) TestAccount_GetSignerFn() {
	cfg := s.helper.BytesFromFile("./testdata/config/accounts.toml")
	accounts := s.factory.Accounts(cfg)
	account := accounts[s.helper.KeyFromMap(accounts)]

	_, err := account.GetSignerFn()

	s.Suite.Require().NoError(err)
}
