package cron_trigger

func (t *Trigger) Close() error {
	t.cron.Stop()

	return nil
}
