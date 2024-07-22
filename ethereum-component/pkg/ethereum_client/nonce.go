package ethereum_client

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"bitbucket.org/unchain/ethereum2/pkg/account"
)

const (
	NonceKeyFMT     = "inflight_nonce:%s"
	NonceLockKeyFMT = "inflight_nonce:lock:%s"
	ExpirationTime  = time.Minute
)

func (c *Client) getNonce(account *account.Account) (*big.Int, error) {
	//See if we can get a new nonce from inflight messages
	inflightNonce, err := c.getInflightNonce(account)
	if err != nil {
		//We don't have any inflight messages
		return nil, nil
	}

	// TODO(e-nikolov) should this be PendingNonceAt?
	//Compare inflight nonce with nonce from node
	nodeNonce, err := c.Client.NonceAt(context.Background(), *account.Address, nil)
	if err != nil {
		//We may be able to recover from this...
		c.logger.Debugf("failed to get account nonce: %s", err.Error())
	}

	if inflightNonce > nodeNonce {
		return big.NewInt(int64(inflightNonce)), nil
	}

	return big.NewInt(int64(nodeNonce)), nil
}

func getLockKey(account *account.Account) string {
	return fmt.Sprintf(NonceLockKeyFMT, account.Address.String())
}

func (c *Client) getInflightNonce(account *account.Account) (uint64, error) {
	key := fmt.Sprintf(NonceKeyFMT, account.Address.String())
	val, err := c.redis.Get(key).Result()
	if err != nil {
		return uint64(0), err
	}

	nonce, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return uint64(0), errors.Wrap(err, "failed to parse nonce")
	}

	return nonce + 1, nil
}

func (c *Client) setNonce(account *account.Account, nonce uint64) error {
	key := fmt.Sprintf(NonceKeyFMT, account.Address.String())
	value := fmt.Sprintf("%d", nonce)

	err := c.redis.Set(key, value, ExpirationTime).Err()
	if err != nil {
		return errors.Wrapf(err, "failed to store nonce: %s", value)
	}

	return nil
}
