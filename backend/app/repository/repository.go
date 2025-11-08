package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

var Db *sql.DB

const directory = "pkg/migrations/database"
const dbPath = "pkg/migrations/database/social.db"


func OpenDb() (*sql.DB, error) {
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=1")
	if err != nil {
		return db, err
	}
	Db = db

	return db, nil
}

func ApplyMigrations(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "pkg/migrations/sqlite",
	}

	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		fmt.Errorf("failed to enable foreign keys PRAGMA: %w", err)
		return nil
	}

	row := db.QueryRow("PRAGMA foreign_keys;")
	var enabled int
	err = row.Scan(&enabled)
	if err != nil {
	} else {
	}

	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("error while executing the migration: %v", err)
	}
	log.Printf("Applied %d migrations successfully!\n", n)
	return nil
}
