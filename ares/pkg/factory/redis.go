package factory

import (
	"github.com/Pallinder/go-randomdata"
)

func (f *Factory) RecoveryCode(email string) string {
	// Create recovery code in kv store for email, and return code
	code := randomdata.RandStringRunes(10)
	err := f.ares.Redis.StoreRecoveryCode(code, email)

	f.suite.Require().NoError(err)

	return code
}
