package db

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

func Init() (*sql.DB, error) {
	path, err := getDBPath()
	if err != nil {
		return nil, err
	}

	db, err := connectDB(path)
	if err != nil {
		return nil, err
	}

	if err = migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func connectDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, err
}

func getDBPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dbDir := filepath.Join(dir, "PassVault")
	dbPath := filepath.Join(dbDir, "vault.db")

	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return "", err
	}

	return dbPath, nil
}

func migrate(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS vault_config (
			id INTEGER PRIMARY KEY,
			username VARCHAR(255),
			argon_salt BLOB NOT NULL,
			encrypted_dek BLOB NOT NULL,
			dek_nonce BLOB NOT NULL
		);

		CREATE TABLE IF NOT EXISTS passwords (
			id INTEGER PRIMARY KEY,
			service VARCHAR(255) NOT NULL,
			encrypted_password BLOB NOT NULL,
			password_nonce BLOB NOT NULL
		);
	`

	ctx, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
	defer cancel()

	if _, err := db.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}