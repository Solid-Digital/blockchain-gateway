package event_test

import (
	"bitbucket.org/unchain/ethereum2/pkg/event"
)

func (s *TestSuite) TestNewEvents() {
	cfg := s.helper.BytesFromFile("./testdata/config/events.toml")
	eventCfgs := s.factory.EventCfgs(cfg)
	contracts := s.factory.Contracts(cfg)

	events, err := event.NewEvents(s.logger, eventCfgs, contracts)

	s.Require().NoError(err)
	s.Require().Equal(2, len(events))
}
