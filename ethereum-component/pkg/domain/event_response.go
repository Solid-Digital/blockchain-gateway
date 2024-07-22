package domain

import (
	"encoding/json"
	"strings"

	"bitbucket.org/unchain/ethereum2/pkg/event"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/unchainio/pkg/errors"
)

type EventResponse struct {
	Event *event.Event
	Log   *types.Log
}

func (e *EventResponse) GetValues() ([]map[string]interface{}, error) {
	// inputs of an event define what should be returned
	ev := e.Event.Event
	inputs := ev.Inputs

	// list non-indexed parameters
	values, err := inputs.UnpackValues(e.Log.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unpack event values")
	}

	// list indexed parameters (first topic is hash of the name of the event)
	topics := e.Log.Topics[1:]

	// add values and topics in the right order
	combined := make([]interface{}, len(inputs))
	for i, input := range inputs {
		if input.Indexed {
			// strip leading zero's, this also changes the data structure from
			// common.Hash to string, but that doesn't impact the final result
			stripped := "0x" + strings.TrimLeft(topics[0].Hex(), "0x")
			combined[i], topics = stripped, topics[1:]
		} else {
			combined[i], values = values[0], values[1:]
		}
	}

	// make list with elements based on input types from abi specification
	typedValues := event.EventOutput(ev)

	// generate bytes so we can unmarshal into typedValues
	bytes, err := json.Marshal(combined)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal topics and values")
	}

	// and unmarshal into typedValues
	err = json.Unmarshal(bytes, &typedValues)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal into typedValues")
	}

	// construct the expected result
	var result []map[string]interface{}
	for index, input := range inputs {
		result = append(result, map[string]interface{}{input.Name: typedValues[index]})
	}

	return result, nil
}
