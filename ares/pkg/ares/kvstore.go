package ares

import "bitbucket.org/unchain/ares/gen/dto"

type KVStore interface {
	IncrementLoginAttempts(ip string, email string, attempts int) error
	ClearLoginAttempts(email string, ip string, attempts int)
	GetLoginAttempts(ip string, email string) (int, error)

	BlacklistToken(token *dto.Token) error
	IsTokenInBlacklist(token string) bool

	StoreSignUpCode(code, email string) error
	GetEmailBySignUpCode(code string) (email string, err error)
	RemoveSignUpCode(code string) error

	StoreRecoveryCode(code, email string) error
	GetEmailByRecoveryCode(code string) (email string, err error)
	RemoveRecoveryCode(code string) error
}
