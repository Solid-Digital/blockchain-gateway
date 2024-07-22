package smtp_action

import (
	"github.com/mitchellh/mapstructure"
)

type Input struct {
	Username   string
	Password   string
	Hostname   string
	Port       string
	From       string
	Recipients []string
	Message    []byte
}

func NewInput(input map[string]interface{}) (*Input, error) {
	var res *Input
	err := mapstructure.Decode(input, &res)
	if err != nil {
		return nil, err
	}

	// check for required fields
	// TODO

	return res, nil
}
