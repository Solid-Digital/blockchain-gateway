package domain

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/core/types"
)

type EventOutput struct {
	Event  string
	Values []map[string]interface{}
	Log    *types.Log
}

func NewEventOutput(response *EventResponse) (map[string]interface{}, error) {
	values, err := response.GetValues()
	if err != nil {
		return nil, err
	}

	output := &EventOutput{
		Event:  response.Event.Event.Name,
		Values: values,
		Log:    response.Log,
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	out := map[string]interface{}{}
	err = json.Unmarshal(bytes, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
