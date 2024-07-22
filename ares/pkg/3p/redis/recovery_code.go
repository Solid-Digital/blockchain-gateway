package redis

import (
	"fmt"
	"time"

	"github.com/unchainio/pkg/errors"
)

const recoveryFmt = "recovery:%s"
const recoveryExpiration = 1000 * time.Second

func (c *Client) GetEmailByRecoveryCode(code string) (email string, err error) {
	email, err = c.Get(fmt.Sprintf(recoveryFmt, code))
	if err != nil {
		return "", errors.WithMessage(err, "could not get recovery code")
	}

	return email, nil
}

// RemoveRecoveryCode removes a recovery code from redis. No errors are returned if the code doesn't exist.
func (c *Client) RemoveRecoveryCode(code string) error {
	_, err := c.Delete(fmt.Sprintf(recoveryFmt, code))
	if err != nil {
		return errors.WithMessage(err, "could not remove recovery code")
	}

	return nil
}

func (c *Client) StoreRecoveryCode(code, email string) error {
	err := c.Set(fmt.Sprintf(recoveryFmt, code), email, recoveryExpiration)
	if err != nil {
		return errors.Wrap(err, "could not store auth recovery object in kv-store")
	}

	return nil
}
