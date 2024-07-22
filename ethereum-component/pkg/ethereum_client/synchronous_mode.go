package ethereum_client

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	Interval = time.Second
	Retries  = 120
)

func (c *Client) txCommitted(hash common.Hash) error {
	//We don't execute in synchronous mode
	if c.cfg.SynchronousMode == false {
		return nil
	}

	for i := 0; i < Retries; i++ {
		_, pending, err := c.Client.TransactionByHash(context.Background(), hash)
		if err != nil {
			return errors.Wrapf(err, "failed to get transaction by hash: %s", hash)
		}

		if pending == false {
			return nil
		}

		time.Sleep(Interval)
	}

	return errors.Errorf("tx is still pending: %s", hash)
}
