package store

import (
	"context"
	"database/sql"
)

type Password struct {
	ID                  int64
	Service           	string
	EncryptedPassword 	[]byte
	Nonce 				[]byte
	Password          	string
}

type PasswordStore struct {
	db *sql.DB
}

func (s *PasswordStore) List(ctx context.Context) ([]Password, error) {
	query := `
		SELECT 
			id, service, encrypted_password, password_nonce
		FROM passwords
	`

	var passwords []Password
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var password Password
		if err := rows.Scan(&password.ID, &password.Service, &password.EncryptedPassword, &password.Nonce); err != nil {
			return nil, err
		}
		passwords = append(passwords, password)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return passwords, nil
}

func (s *PasswordStore) Add(ctx context.Context, service string, password, nonce []byte) error {
	query := `
		INSERT INTO passwords (service, encrypted_password, password_nonce)
		VALUES(?, ?, ?)
	`

	if _, err := s.db.ExecContext(ctx, query, service, password, nonce); err != nil {
		return err
	}

	return nil
}

func (s *PasswordStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM passwords
		WHERE id = ?
	`

	if _, err := s.db.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}