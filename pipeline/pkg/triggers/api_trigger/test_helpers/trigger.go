package test_helpers

import (
	"github.com/unchain/pipeline/pkg/triggers/api_trigger"
	"github.com/unchainio/pkg/errors"
	"time"
)

const (
	TestPort = "8888"
)

func (t *TestHelpers) InitializedTrigger(cfg []byte) *api_trigger.Trigger {
	trigger := new(api_trigger.Trigger)
	trigger.Init(t.logger, cfg)

	return trigger
}

func (t *TestHelpers) InitializedTriggerWithError(cfg []byte) (*api_trigger.Trigger, error) {
	trigger := new(api_trigger.Trigger)
	err := trigger.Init(t.logger, cfg)
	if err != nil {
		return nil, err
	}
	return trigger, nil
}

func (t *TestHelpers) TriggerResponse(trigger *api_trigger.Trigger, seconds int) (string, map[string]interface{}, error) {
	type triggerResponse struct {
		tag      string
		response map[string]interface{}
		err      error
	}

	responseChan := make(chan triggerResponse)
	defer close(responseChan)

	go func() {
		tag, response, err := trigger.NextMessage()
		responseChan <- triggerResponse{
			tag:      tag,
			response: response,
			err:      err,
		}
	}()

	for i := 0; i < seconds; i++ {
		select {
		case response := <-responseChan:
			return response.tag, response.response, response.err
		default:
			time.Sleep(time.Second)
		}
	}

	return "", nil, errors.New("trigger did not receive anything")
}

func (t *TestHelpers) RespondResponse(trigger *api_trigger.Trigger, tag string, input map[string]interface{}, err error, seconds int) error {
	responseChan := make(chan error)
	defer close(responseChan)

	go func() {
		responseChan <- trigger.Respond(tag, input, err)
	}()

	for i := 0; i < seconds; i++ {
		select {
		case response := <-responseChan:
			return response
		default:
			time.Sleep(time.Second)
		}
	}

	return errors.New("respond did not receive anything")
}
