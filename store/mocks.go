package store

import "context"

type MockConfigStore struct {
	salt, dek, nonce []byte
}

func NewMockStore() *Store {
	return &Store{
		Config: &MockConfigStore{},
	}
}

func (s *MockConfigStore) Get(ctx context.Context) ([]byte, []byte, []byte, error) {
	return s.salt, s.dek, s.nonce, nil
}

func (s *MockConfigStore) Save(ctx context.Context, salt, dek, nonce []byte) error {
	s.salt, s.dek, s.nonce = salt, dek, nonce
	return nil
}
