package event

import (
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func EventOutput(event *abi.Event) []interface{} {
	inputs := event.Inputs

	var ret []interface{}
	for i, _ := range inputs {
		ret = append(ret, reflect.New(inputs[i].Type.Type).Interface())
	}

	return ret
}
