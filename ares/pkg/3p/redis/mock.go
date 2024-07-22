package redis

type MockClient struct {
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (MockClient) SetValue(key string, value string, expiration ...interface{}) error {
	return nil
}

func (MockClient) GetValue(key string) (interface{}, error) {
	return nil, nil
}

func (MockClient) Close() error {
	return nil
}
