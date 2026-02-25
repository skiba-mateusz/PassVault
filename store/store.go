package store

import (
	"context"
	"database/sql"
)

type Store struct {
	Config interface {
		Get(ctx context.Context) (string, []byte, []byte, []byte, error)
		Save(ctx context.Context, username string, salt, dek, nonce []byte) error
	}
	Password interface {
		List(ctx context.Context) ([]Password, error)
		Add(ctx context.Context, service string, password, nonce []byte) error
		Delete(ctx context.Context, id int64) error
	}
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Config: &ConfigStore{
			db: db,
		},
		Password: &PasswordStore{
			db: db,
		},
	}
}