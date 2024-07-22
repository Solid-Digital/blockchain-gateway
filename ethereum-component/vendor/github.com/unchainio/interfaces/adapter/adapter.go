package adapter

import "github.com/unchainio/interfaces/logger"

type ComponentType string

const (
	TriggerComponent = "trigger"
	ActionComponent  = "action"
)

type Trigger interface {
	// Init: must NOT block, start long running processes in a go routine
	Init(stub Stub, config []byte) (err error)

	// Trigger: must block until a new message is received
	Trigger() (tag string, message map[string]interface{}, err error)

	// Respond: is called by the adapter base after the message (with tag `tag`), which was initially received
	// by the Trigger, has passed through the actions in the pipeline. Respond is used to send a response back to
	// the caller of the Trigger. In case of an error during invocation of an action component, the adapter base will
	// discontinue invocation of other actions in the pipeline and call Respond.
	Respond(tag string, response map[string]interface{}, err error) error

	Close() error
}

type Action interface {
	Init(stub Stub, config []byte) (err error)
	Invoke(input map[string]interface{}) (output map[string]interface{}, err error)
}

type Stub interface {
	logger.Logger

	// TODO in the future this interface will also contain a kv store and a secret store
}
