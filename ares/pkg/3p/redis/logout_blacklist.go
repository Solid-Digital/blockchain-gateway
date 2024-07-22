package redis

import (
	"fmt"
	"time"

	"bitbucket.org/unchain/ares/gen/dto"
)

const tokenFmt = "auth-token:%s"

func (c *Client) IsTokenInBlacklist(token string) bool {
	value, err := c.Get(fmt.Sprintf(tokenFmt, token))
	// If redis errors out, all tokens will be virtually blacklisted.
	return err == nil || value != ""
}

func (c *Client) BlacklistToken(token *dto.Token) error {
	return c.Set(
		fmt.Sprintf(tokenFmt, token.Raw),
		token.Raw,
		getTokenRemainingValidity(token.Expiration),
	)
}

const (
	expireOffset = 3600
)

func getTokenRemainingValidity(timestamp int64) time.Duration {
	tm := time.Unix(timestamp, 0) // FIXME this doesn't look right
	remainder := time.Until(tm)
	if remainder > 0 {
		t := time.Second * time.Duration(remainder.Seconds()+expireOffset)
		return t
	}
	return time.Second * expireOffset
}
