package store

import (
	"context"
	"fmt"
)

type MockConfigStore struct {
	salt, dek, nonce []byte
	username string
}

type MockPasswordStore struct {
	passwords []Password
	id int64
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
	if len(s.passwords) == 0 {
		return nil, fmt.Errorf("No passwords")
	}

	return s.passwords, nil
}

func (s *MockPasswordStore) Add(ctx context.Context, service string, password, nonce []byte) error {
	s.passwords = append(s.passwords, Password{
		ID: s.id,
		Service: service,
		EncryptedPassword: password,
		Nonce: nonce,
	})
	s.id++
	return nil
}

func (s *MockPasswordStore) Delete(ctx context.Context, id int64) error {
	for idx, pass := range s.passwords {
		if pass.ID == id {
			s.passwords = append(s.passwords[:idx], s.passwords[idx+1:]...)
			return nil
		}
	}

	return nil
}

func (s *MockPasswordStore) Edit(ctx context.Context, id int64, newService string) error {
	for idx, pass := range s.passwords {
		if pass.ID == id {
			s.passwords[idx].Service = newService
			return nil
		}
	}
	
	return nil
}