package ethereum_client

import (
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const (
	TypeSingleOutput   = "single"
	TypeMultipleOutput = "multiple"
)

func callOutputType(function *abi.Method) string {
	if len(function.Outputs) <= 1 {
		return TypeSingleOutput
	}

	return TypeMultipleOutput
}

func singleOutput(function *abi.Method) interface{} {
	outputs := function.Outputs

	if len(outputs) == 0 {
		return nil
	}

	return reflect.New(outputs[0].Type.Type).Interface()
}

func multipleOutput(function *abi.Method) []interface{} {
	outputs := function.Outputs

	var ret []interface{}
	for i := range outputs {
		ret = append(ret, reflect.New(outputs[i].Type.Type).Interface())
	}

	return ret
}
