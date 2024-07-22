package redis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/unchainio/pkg/errors"
)

const loginAttemptsFmt = "login-attempts:%s:%s"
const loginAttemptsExpiryDuration = 15 * time.Minute

func (c *Client) GetLoginAttempts(ip, email string) (attempts int, err error) {
	attemptsString, err := c.Get(fmt.Sprintf(loginAttemptsFmt, ip, email))
	if err != nil {
		return 0, err
	}

	if attemptsString != "" {
		attempts, err = strconv.Atoi(attemptsString)

		if err != nil {
			return 0, errors.Wrap(err, "")
		}
	}

	return attempts, nil
}

func (c *Client) IncrementLoginAttempts(ip string, email string, attempts int) error {
	attempts++

	attemptsStr := strconv.Itoa(attempts)

	return c.Set(fmt.Sprintf(loginAttemptsFmt, ip, email), attemptsStr, loginAttemptsExpiryDuration)
}

func (c *Client) ClearLoginAttempts(email string, ip string, attempts int) {
	// clears any existing login attempts from kv store
	_, err := c.Delete(fmt.Sprintf(loginAttemptsFmt, ip, email))
	if err != nil {
		c.log.Errorf("failed to remove login attempt for - %v", err)
	}
}
