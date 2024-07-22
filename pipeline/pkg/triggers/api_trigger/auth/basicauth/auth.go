package basicauth

import "encoding/base64"

type AuthService struct {
	AuthStrings []string
}

type Credentials struct {
	Username string
	Password string
}

func NewService(credentials []Credentials) *AuthService {
	var authStrings []string
	for _, c := range credentials {
		authStrings = append(authStrings, createBasicAuthString(c))
	}

	return &AuthService{
		AuthStrings: authStrings,
	}
}

func createBasicAuthString(c Credentials) string {
	auth := c.Username + ":" + c.Password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}