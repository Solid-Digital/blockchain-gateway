package domain

type Trigger interface {
	// Init: must NOT block, start long running processes in a go routine
	Init(stub Stub, config []byte) (err error)

	// NextMessage: request the next message that was produced/receveid by the trigger
	NextMessage() (tag string, message map[string]interface{}, err error)

	// Respond: is called by the adapter base after the message (with tag `tag`), which was initially received
	// by the Trigger, has passed through the actions in the pipeline. Respond is used to send a response back to
	// the caller of the Trigger. In case of an error during invocation of an action component, the adapter base will
	// discontinue invocation of other actions in the pipeline and call Respond.
	Respond(tag string, response map[string]interface{}, err error) error

	Close() error
}
