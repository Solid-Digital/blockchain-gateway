package event

import (
	"reflect"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/unchainio/pkg/errors"
)

// generate filters of the types conform the ABI specification
func buildFilters(inputs abi.Arguments, filters [][]toml.Primitive) ([][]interface{}, error) {

	var result [][]interface{}

	// filterList contains a list of filter values
	for _, filterList := range filters {

		var filterResult []interface{}

		// value contains a filter value, of which the type corresponds to the
		// event input type from the ABI (in order, so with index i)
		for i, value := range filterList {

			// value should be of this go type
			input := inputs[i]
			goType := input.Type.Type
			valPtr := reflect.New(goType).Interface()

			// FIXME: deprecated, metadata
			err := toml.PrimitiveDecode(value, valPtr)
			if err != nil {
				return nil, errors.Wrap(err, "failed parsing toml filters")
			}

			// dereference pointer to goType
			val := reflect.ValueOf(valPtr).Elem().Interface()
			filterResult = append(filterResult, val)
		}

		result = append(result, filterResult)
	}

	return result, nil
}
