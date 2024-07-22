package event

import (
	"bitbucket.org/unchain/ethereum2/pkg/contract"
	"github.com/unchainio/interfaces/logger"
)

type Events []*Event

func NewEvents(logger logger.Logger, cfgs []*Config, contracts contract.Contracts) (Events, error) {
	events := Events{}

	for _, cfg := range cfgs {
		event, err := newEvent(logger, cfg, contracts)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}
