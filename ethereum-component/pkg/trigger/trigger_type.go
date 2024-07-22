package trigger

import (
	"bitbucket.org/unchain/ethereum2/pkg/domain"
	"github.com/unchainio/interfaces/adapter"
)

type Trigger struct {
	cfg             *Config
	stub            adapter.Stub
	responseChannel chan *domain.EventResponse
	errorChannel    chan error
}

func (t *Trigger) Respond(tag string, response map[string]interface{}, responseError error) error {
	return nil
}

func (t *Trigger) Close() error {
	return nil
}
