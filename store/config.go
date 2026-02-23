package store

import (
	"context"
	"database/sql"
	"time"
)

type ConfigStore struct {
	db *sql.DB
}

func (s *ConfigStore) Get(ctx context.Context) (string, []byte, []byte, []byte, error) {
	query := `
		SELECT 
			username, argon_salt, encrypted_dek, dek_nonce 
		FROM vault_config 
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 15 * time.Second)
	defer cancel()

	var salt, dek, nonce []byte
	var username string

	row := s.db.QueryRowContext(ctx, query)
	if err := row.Scan(&username, &salt, &dek, &nonce); err != nil {
		return "", nil, nil, nil, err
	}

	return username, salt, dek, nonce, nil
}

func (s *ConfigStore) Save(ctx context.Context, username string, salt, dek, nonce []byte) (error) {
	query := `
		INSERT INTO vault_config 
			(username, argon_salt, encrypted_dek, dek_nonce) 
		VALUES (?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, 15 * time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, username, salt, dek, nonce); err != nil {
		return err
	}

	return nil
}