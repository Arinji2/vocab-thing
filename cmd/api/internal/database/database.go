package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func SetupDatabase(dbPath string) (*sql.DB, error) {
	fmt.Println("Setting up database")
	connURL := fmt.Sprintf("file:%s?_foreign_keys=1&_journal_mode=WAL", dbPath)
	db, err := sql.Open("sqlite3", connURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("SQLite database successfully opened")
	return db, nil
}
