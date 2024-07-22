package domain

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

const (
	CallContractFunction = "callContractFunction"
	DeployContract       = "deployContract"
)

type MessageCallContractFunction struct {
	From     string
	To       string
	Function string
	Params   map[string]interface{}
	Nonce    uint64
}

type MessageDeployContract struct {
	From              string
	Solidity          string
	ConstructorParams map[string]map[string]interface{}
	Nonce             uint64
}

func NewMessage(input map[string]interface{}) (interface{}, error) {
	msgType, err := messageType(input)
	if err != nil {
		return nil, err
	}

	var msg interface{}
	switch msgType {
	case CallContractFunction:
		msg, err = newMessageCallContractFunction(input)
	case DeployContract:
		msg, err = newMessageDeployContract(input)
	default:
		return nil, errors.New(fmt.Sprintf("invalid message type '%s'", msgType))
	}

	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *MessageCallContractFunction) validate() error {
	if m.From == "" {
		return errors.New("required field from is empty")
	}

	if m.To == "" {
		return errors.New("required field to is empty")
	}

	if m.Function == "" {
		return errors.New("required field function is empty")
	}

	return nil
}

func (m *MessageDeployContract) validate() error {
	if m.From == "" {
		return errors.New("required field from is empty")
	}

	if m.Solidity == "" {
		return errors.New("required field solidity is empty")
	}

	return nil
}

func newMessageCallContractFunction(input map[string]interface{}) (*MessageCallContractFunction, error) {
	msg := new(MessageCallContractFunction)
	err := mapstructure.Decode(input, msg)
	if err != nil {
		return nil, errors.Wrap(err, "input for callContractFunction not valid")
	}

	err = msg.validate()
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func newMessageDeployContract(input map[string]interface{}) (*MessageDeployContract, error) {
	msg := new(MessageDeployContract)
	err := mapstructure.Decode(input, msg)
	if err != nil {
		return nil, errors.Wrap(err, "input for deployContract not valid")
	}

	err = msg.validate()
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func messageType(input map[string]interface{}) (string, error) {
	type message struct {
		Type string
	}

	msg := new(message)
	err := mapstructure.Decode(input, msg)
	if err != nil {
		return "", errors.Wrap(err, "failed to get message type")
	}

	return msg.Type, nil
}
