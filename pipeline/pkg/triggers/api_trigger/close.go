package api_trigger

func (t *Trigger) Close() error {
	err := t.client.Shutdown()
	if err != nil {
		return err
	}

	return nil
}
