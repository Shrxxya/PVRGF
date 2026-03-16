package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)
type VaultStorage interface {
	Close() error
	GetDB() *sql.DB
}
type SQLStorage struct {
	db *sql.DB
}

func (s *SQLStorage) Close() error {
	return s.db.Close()
}

func (s *SQLStorage) GetDB() *sql.DB {
	return s.db
}

func InitDB() VaultStorage {
	db, err := sql.Open("sqlite", "./vault.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database not reachable:", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Failed to enable foreign keys:", err)
	}

	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`

	passwordTable := `
	CREATE TABLE IF NOT EXISTS passwords (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		label TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err = db.Exec(userTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	_, err = db.Exec(passwordTable)
	if err != nil {
		log.Fatal("Failed to create passwords table:", err)
	}

	return &SQLStorage{db: db}
}
