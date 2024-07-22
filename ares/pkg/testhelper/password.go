package testhelper

import (
	"bitbucket.org/unchain/ares/pkg/auth"
)

func (h *Helper) CorrectPassword(password, passwordHash string) bool {
	if err := auth.CompareHashAndPassword(passwordHash, password); err != nil {
		return false
	}

	return true
}

func (h *Helper) CorrectUserPassword(userID int64, password string) bool {
	user := h.DBGetUser(userID)

	h.suite.Require().NotNil(user)

	return h.CorrectPassword(password, user.PasswordHash)
}
