package auth

import (
	"strconv"
	"strings"

	"bitbucket.org/unchain/ares/gen/dto"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/strfmt"
	"github.com/unchainio/pkg/errors"
)

func (s *Service) Authenticate(token string) (*dto.User, error) {
	user, err := s.authenticate(token)
	if err != nil {
		s.log.Errorf("%+v", err)

		return nil, err
	}

	return user, nil
}

func (s *Service) authenticate(token string) (*dto.User, error) {
	token = trimPrefixIgnoreCase(token, "Bearer ")

	t, err := s.TokenAuth.Decode(token)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid token %q", token)
	}

	if t == nil || !t.Valid || t.Claims.Valid() != nil || s.kv.IsTokenInBlacklist(t.Raw) {
		return nil, errors.Errorf("invalid token %q", token)
	}

	userID, ok := t.Claims.(jwt.MapClaims)["id"].(string)
	if !ok {
		return nil, errors.New("Unable to get user id from jwt")
	}

	userEmail, ok := t.Claims.(jwt.MapClaims)["email"].(string)
	if !ok {
		return nil, errors.New("Unable to get user email from jwt")
	}

	expiration, ok := t.Claims.(jwt.MapClaims)["exp"].(float64)
	if !ok {
		return nil, errors.New("Unable to get token expiration from jwt")
	}

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "\"%s\" is not a valid id", userID)
	}

	return &dto.User{
		Email: strfmt.Email(userEmail),
		ID:    id,
		Token: &dto.Token{
			Expiration: int64(expiration),
			Raw:        t.Raw,
		},
	}, nil
}

// hasPrefixIgnoreCase tests whether the string s begins with prefix while ignoring the casing.
func hasPrefixIgnoreCase(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.ToLower(s[0:len(prefix)]) == strings.ToLower(prefix)
}

// trimPrefixIgnoreCase returns s without the provided leading prefix string.
// If s doesn't start with prefix, s is returned unchanged.
func trimPrefixIgnoreCase(s, prefix string) string {
	if hasPrefixIgnoreCase(s, prefix) {
		return s[len(prefix):]
	}
	return s
}
