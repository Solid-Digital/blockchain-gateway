package redis

import (
	"fmt"
	"time"

	"github.com/unchainio/pkg/errors"
)

const signUpFmt = "sign-up:%s"
const signUpExpiration = 24 * time.Hour

func (c *Client) StoreSignUpCode(code, email string) error {
	err := c.Set(fmt.Sprintf(signUpFmt, code), email, signUpExpiration)
	if err != nil {
		return errors.Wrap(err, "could not store auth sign up object in kv-store")
	}

	return nil
}

func (c *Client) GetEmailBySignUpCode(code string) (email string, err error) {
	email, err = c.Get(fmt.Sprintf(signUpFmt, code))
	if err != nil {
		return "", err
	}

	return email, nil
}

// RemoveSignUpCode removes a sign-up code from redis. No errors are returned if the code doesn't exist.
func (c *Client) RemoveSignUpCode(code string) error {
	_, err := c.Delete(fmt.Sprintf(signUpFmt, code))
	if err != nil {
		return errors.WithMessage(err, "could not remove recovery code")
	}

	return nil
}
