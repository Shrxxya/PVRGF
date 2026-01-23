package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite", "./vault.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	return db
} //zero knowledge proof
//integratoin with websites to fetch their specific pwd validation regex
//openpassword - open source vault
//user pwd encryption - no compromise
