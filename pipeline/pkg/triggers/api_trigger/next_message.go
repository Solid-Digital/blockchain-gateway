package api_trigger

func (t *Trigger) NextMessage() (tag string, request map[string]interface{}, err error) {
	req := <-t.client.RequestChannel

	return req.Tag, req.Output, req.Error
}
