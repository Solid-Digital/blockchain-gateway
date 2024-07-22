package templater_action

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type Message struct {
	Template  string
	Variables map[string]interface{}
}

func NewMessage(input map[string]interface{}) (*Message, error) {
	var msg *Message

	err := mapstructure.Decode(input, &msg)
	if err != nil {
		return nil, err
	}

	if msg.Template == "" {
		return nil, errors.New("no template found")
	}

	return msg, nil
}
