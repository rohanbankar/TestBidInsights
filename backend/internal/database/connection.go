package database

import (
	"database/sql"
	"log"

	_ "github.com/marcboeker/go-duckdb"
)

func Connect(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("duckdb", dbPath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Printf("Connected to DuckDB at %s", dbPath)
	return db, nil
}

func Close(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}