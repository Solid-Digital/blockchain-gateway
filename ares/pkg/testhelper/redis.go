package testhelper

func (h *Helper) EmailBySignUpCodeExists(code string) bool {
	_, err := h.ares.Redis.GetEmailBySignUpCode(code)

	if err != nil {
		return false
	}

	return true
}

func (h *Helper) RedisGetEmailFromSignUpCode(code string) string {
	email, err := h.ares.Redis.GetEmailBySignUpCode(code)
	h.suite.Require().NoError(err)

	return email
}
