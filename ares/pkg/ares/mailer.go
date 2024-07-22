package ares

type Mailer interface {
	SendMessage(to string, subject string, body string) (err error)
	SendRecoveryMessage(to, requestID, recoveryURL, recoveryCode, username string) (err error)
	SendSignUpMessage(to, requestID, signUpURL string) (err error)
	SendInviteMessage(to, requestID, signUpURL, orgName string) (err error)
}
