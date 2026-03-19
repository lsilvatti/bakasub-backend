package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func createTables(db *sql.DB) error {
	createUserConfigTable := `
	CREATE TABLE IF NOT EXISTS user_config (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		default_model TEXT NOT NULL,
		default_preset TEXT NOT NULL,
		remove_sdh_default BOOLEAN NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	createFavoriteFoldersTable := `
	CREATE TABLE IF NOT EXISTS favorite_folders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		alias TEXT NOT NULL,
		path TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createUserConfigTable)

	if err != nil {
		return fmt.Errorf("failed to create user_config table: %w", err)
	}

	_, err = db.Exec(createFavoriteFoldersTable)

	if err != nil {
		return fmt.Errorf("failed to create favorite_folders table: %w", err)
	}

	return nil

}

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = createTables(db)

	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("SQLite database initialized successfully")

	return db, nil
}
