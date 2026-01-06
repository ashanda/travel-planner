package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// create tables
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE,
		name TEXT,
		picture TEXT,
		created_at INTEGER
	);

	CREATE TABLE IF NOT EXISTS usage (
		user_id TEXT PRIMARY KEY,
		generations INTEGER DEFAULT 0,
		updated_at INTEGER
	);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
