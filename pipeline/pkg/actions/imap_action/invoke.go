package imap_action

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/unchain/pipeline/pkg/domain"
	"github.com/unchainio/pkg/xconfig"
	"strings"
)

const (
	ConfigInput = "config"
	Function = "function"
	Params = "params"
)

func Invoke(stub domain.Stub, input map[string]interface{}) (output map[string]interface{}, err error) {
	// Parse config
	cfg := new(Config)

	err = xconfig.Load(
		cfg,
		xconfig.FromReaders(
			"toml",
			strings.NewReader(
				fmt.Sprintf("%v", input[ConfigInput]),
			),
		),
	)
	if err != nil {
		stub.Errorf("could not load configuration: %v\n", err)
		return nil, errors.Wrap(err, "could not load configuration")
	}

	// Get client
	client, err := NewClient(stub, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "could not get new client")
	}

	// Perform action based on switch case decided upon instruction
	switch fmt.Sprintf("%s", input[Function]) {
	case "GetNewMessageAttachments":
		// -- Get new message attachments
		attachments, err := client.GetNewMessageAttachments()
		if err != nil {
			return nil, errors.Wrap(err, "could not get new message attachments")
		}
		err = client.Client.Terminate()
		if err != nil {
			return nil, errors.Wrapf(err, "could not terminate client")
		}
		return attachments, nil
	case "MarkMessageAsRead":
		params, ok := input[Params].(map[string]interface{})
		if (!ok) {
			return nil, errors.New("could not cast params to map[string]interface{}")
		}
		seqNum, ok := params["seqNum"].(int)
		if (!ok) {
			return nil, errors.New("could not cast seqNum to int")
		}
		err = client.MarkMessageAsRead(uint32(seqNum))
		if err != nil {
			return nil, errors.Wrap(err, "could not mark message as read")
		}
		err = client.Client.Terminate()
		if err != nil {
			return nil, errors.Wrapf(err, "could not terminate client")
		}
		return nil, nil
	case "MoveFailedMessage":
		params, ok := input[Params].(map[string]interface{})
		if (!ok) {
			return nil, errors.New("could not cast params to map[string]interface{}")
		}
		seqNum, ok := params["seqNum"].(int)
		if (!ok) {
			return nil, errors.New("could not cast seqNum to int")
		}
		err = client.MoveFailedMessage(uint32(seqNum))
		if err != nil {
			return nil, errors.Wrap(err, "could not move failed message")
		}
		err = client.Client.Terminate()
		if err != nil {
			return nil, errors.Wrapf(err, "could not terminate client")
		}
		return nil, nil
	case "NoOp":
		err = client.Client.Terminate()
		if err != nil {
			return nil, errors.Wrapf(err, "could not terminate client")
		}
		return nil, nil
	default:
		return nil, errors.New(fmt.Sprintf(`
		Unknown command. 

The following commands are available:
- GetNewMessageAttachments
- MarkMessageAsRead 
- MoveFailedMessage
		`))
	}
}
