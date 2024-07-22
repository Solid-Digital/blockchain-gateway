package auth

func (s *Service) GetConnectURL() string {
	return s.cfg.ConnectURL
}
