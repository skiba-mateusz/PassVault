package store

import "context"

type MockConfigStore struct {
	salt, dek, nonce []byte
	username string
}

func NewMockStore() *Store {
	return &Store{
		Config: &MockConfigStore{},
	}
}

func (s *MockConfigStore) Get(ctx context.Context) ([]byte, []byte, []byte, error) {
	return s.salt, s.dek, s.nonce, nil
}

func (s *MockConfigStore) Save(ctx context.Context, username string, salt, dek, nonce []byte) error {
	s.salt, s.dek, s.nonce, s.username = salt, dek, nonce, username
	return nil
}
