package cron_trigger

func (t *Trigger) NextMessage() (tag string, message map[string]interface{}, err error) {
	req := <-t.RequestChannel

	return req.Tag, nil, nil
}
