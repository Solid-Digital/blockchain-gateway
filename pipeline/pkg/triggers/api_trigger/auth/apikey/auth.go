package apikey

type AuthService struct {
	APIKeys []string
}

func NewService(keys []string) *AuthService {
	return &AuthService{
		APIKeys: keys,
	}
}
