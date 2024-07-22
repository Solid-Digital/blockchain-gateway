package api_trigger

import (
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchainio/pkg/errors"
)

func (t *Trigger) Respond(tag string, response map[string]interface{}, err error) error {
	responseChannel, ok := t.client.ResponseChannelMap.Load(tag)
	if !ok {
		return errors.Errorf("failed to find response channel for message with tag %s", tag)
	}

	channel, ok := responseChannel.(chan *domain.Response)
	if !ok {
		return errors.Errorf("response channel type is invalid")
	}

	channel <- &domain.Response{
		Message: response,
		Error:   err,
	}

	return nil
}
