package trigger

import (
	"bitbucket.org/unchain/ethereum2/pkg/domain"
	"github.com/unchainio/interfaces/adapter"
)

func (t *Trigger) Trigger() (tag string, response map[string]interface{}, err error) {
	tag = adapter.NewTag()

	select {
	case eventResponse := <-t.responseChannel:
		output, err := domain.NewEventOutput(eventResponse)
		if err != nil {
			t.stub.Errorf(err.Error())

			return tag, nil, err
		}

		return tag, output, nil
	case err := <-t.errorChannel:
		// What will happen in this case, how does janus handle errors?
		t.stub.Errorf(err.Error())

		return tag, nil, err
	}
}
