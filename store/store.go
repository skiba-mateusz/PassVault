package store

import (
	"context"
	"database/sql"
)

type Store struct {
	Config interface {
		Get(ctx context.Context) ([]byte, []byte, []byte, error)
		Save(ctx context.Context, username string, salt, dek, nonce []byte) error
	}
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Config: &ConfigStore{
			db: db,
		},
	}
}