package domain

type Action interface {
	Invoke(stub Stub, input map[string]interface{}) (output map[string]interface{}, err error)
}
