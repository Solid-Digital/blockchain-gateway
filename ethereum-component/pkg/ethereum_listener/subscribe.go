package ethereum_listener

import (
	"bitbucket.org/unchain/ethereum2/pkg/domain"
	"bitbucket.org/unchain/ethereum2/pkg/event"
	"github.com/unchainio/pkg/errors"
)

func (l *Listener) subscribeAllEvents() {
	for _, event := range l.events {
		go l.subscribe(event)
	}
}

func (l *Listener) subscribe(event *event.Event) {
	boundContract := l.boundContract(event.Contract)

	logChannel, subscription, err := boundContract.WatchLogs(nil, event.Event.Name, event.Filters...)
	if err != nil {
		l.ErrorChannel <- errors.Wrap(err, "failed setting up subscription")
		return
	}

	for {
		select {
		case err := <-subscription.Err():
			// FIXME: somehow a lot errors are put on this channel; can we just ignore them?
			// FIXME: or maybe just log them, instead of putting it on the error channel
			l.ErrorChannel <- errors.Wrap(err, "subscription error")
			subscription.Unsubscribe()
			l.subscribe(event)
		case log := <-logChannel:
			eventResponse := &domain.EventResponse{
				Event: event,
				Log:   &log,
			}

			l.ResponseChannel <- eventResponse
		}
	}
}
