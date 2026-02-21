package store

import (
	"context"
	"database/sql"
	"time"
)

type ConfigStore struct {
	db *sql.DB
}

func (s *ConfigStore) Get(ctx context.Context) ([]byte, []byte, []byte, error) {
	query := `
		SELECT 
			argon_salt, encrypted_dek, dek_nonce 
		FROM vault_config 
		LIMIT 1
	`

	ctx, cancel := context.WithTimeout(ctx, 15 * time.Second)
	defer cancel()

	var salt, dek, nonce []byte

	row := s.db.QueryRowContext(ctx, query)
	if err := row.Scan(&salt, &dek, &nonce); err != nil {
		return nil, nil, nil, err
	}

	return salt, dek, nonce, nil
}

func (s *ConfigStore) Save(ctx context.Context, salt, dek, nonce []byte) (error) {
	query := `
		INSERT INTO vault_config 
			(argon_salt, encrypted_dek, dek_nonce) 
		VALUES (?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, 15 * time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, salt, dek, nonce); err != nil {
		return err
	}

	return nil
}