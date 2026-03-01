package store

import "context"

type MockConfigStore struct {
	salt, dek, nonce []byte
	username string
}

type MockPasswordStore struct {
	passwords []Password
}

func NewMockStore() *Store {
	return &Store{
		Config: &MockConfigStore{},
		Password: &MockPasswordStore{
			passwords: []Password{},
		},
	}
}

func (s *MockConfigStore) Get(ctx context.Context) (string, []byte, []byte, []byte, error) {
	return s.username, s.salt, s.dek, s.nonce, nil
}

func (s *MockConfigStore) Save(ctx context.Context, username string, salt, dek, nonce []byte) error {
	s.salt, s.dek, s.nonce, s.username = salt, dek, nonce, username
	return nil
}

func (s *MockPasswordStore) List(ctx context.Context) ([]Password, error) {
	return s.passwords, nil
}

func (s *MockPasswordStore) Add(ctx context.Context, service string, password, nonce []byte) error {
	s.passwords = append(s.passwords, Password{
		Service: service,
		EncryptedPassword: password,
		Nonce: nonce,
	})
	return nil
}

func (s *MockPasswordStore) Delete(ctx context.Context, id int64) error {
	return nil
}

func (s *MockPasswordStore) Edit(ctx context.Context, id int64, newService string) error {
	return nil
}