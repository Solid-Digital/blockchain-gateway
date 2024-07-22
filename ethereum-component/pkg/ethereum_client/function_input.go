package ethereum_client

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
)

//build function inputs conform the abi specification
func functionInput(function *abi.Method, params map[string]interface{}) ([]interface{}, error) {
	var inputs []interface{}

	//loop over input params from abi specification
	for _, input := range function.Inputs {
		value, ok := params[input.Name]
		if !ok {
			return nil, errors.New(fmt.Sprintf("missing parameter: %s", input.Name))
		}

		//value should be of this go type
		goType := input.Type.Type
		valPtr := reflect.New(goType).Interface()

		//TODO: try if this can be done in a shorter way
		byteValue, err := json.Marshal(value)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert value to byte: %s", value)
		}

		//assign value
		err = json.Unmarshal(byteValue, valPtr)
		if err != nil {
			return nil, err
		}

		//dereference pointer to goType
		val := reflect.ValueOf(valPtr).Elem().Interface()
		inputs = append(inputs, val)
	}

	return inputs, nil
}

//build constructor inputs conform the abi specification
func constructorInput(contractName string, constructor *abi.Method, params map[string]map[string]interface{}) ([]interface{}, error) {
	if len(constructor.Inputs) == 0 {
		return []interface{}{}, nil
	}

	constructorParams, ok := params[contractName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no constructor params provided for contract '%s'", contractName))
	}

	inputs, err := functionInput(constructor, constructorParams)
	if err != nil {
		return nil, err
	}

	return inputs, nil
}
