package ethereum_listener

import (
	"bitbucket.org/unchain/ethereum2/pkg/domain"
	"bitbucket.org/unchain/ethereum2/pkg/event"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/unchainio/interfaces/logger"
	"github.com/unchainio/pkg/errors"
)

type Listener struct {
	logger          logger.Logger
	cfg             *Config
	events          event.Events
	client          *ethclient.Client
	ResponseChannel chan *domain.EventResponse
	ErrorChannel    chan error
}

func NewListener(logger logger.Logger, cfg *Config, events event.Events) (*Listener, error) {
	client, err := ethclient.Dial(cfg.Host)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start ethereum client")
	}

	listener := &Listener{
		logger:          logger,
		cfg:             cfg,
		events:          events,
		client:          client,
		ResponseChannel: make(chan *domain.EventResponse),
		ErrorChannel:    make(chan error),
	}

	listener.subscribeAllEvents()

	return listener, nil
}
