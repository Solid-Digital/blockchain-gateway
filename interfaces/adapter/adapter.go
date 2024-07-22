package adapter

import "github.com/unchainio/interfaces/logger"

type Endpoint interface {
	// Init: must NOT block, start long running processes in a go routine
	Init(stub Stub, config []byte) (err error)

	// Send: must block until sending is complete
	Send(stub Stub, message *Message) (response *Message, err error)

	// Receive: must block until a new message is received
	Receive(stub Stub) (message *TaggedMessage, err error)

	// Ack is called by the adapter base after the message (with tag `tag`), which was initially received
	// by this input endpoint, has been successfully passed through the actions in the pipeline, sent
	// over the output endpoint, and a response has been returned from the output endpoint and passed
	// through the actions in the response pipeline.
	Ack(stub Stub, tag uint64, response *Message) error

	// Nack is called by the adapter base if anything goes wrong while processing the message with tag `tag`
	Nack(stub Stub, tag uint64, err error) error

	Close(stub Stub) error
}

type Action interface {
	Init(stub Stub, config []byte) (err error)
	Invoke(stub Stub, message *Message) (err error)
}

type Stub interface {
	logger.Logger

	// TODO in the future this interface will also contain a kv store and a secret store
}
