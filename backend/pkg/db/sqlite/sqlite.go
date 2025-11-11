package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

var Db *sql.DB

const (
	dbPath = "pkg/db/sqlite/social.db"
)

func OpenDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=1&_busy_timeout=5000")
	if err != nil {
		return db, err
	}
	Db = db

	return db, nil
}

func ApplyMigrations(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "pkg/db/migrations/sqlite",
	}

	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("failed to enable foreign keys PRAGMA: %w", err)
	}

	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("error while executing the migration: %v", err)
	}
	log.Printf("Applied %d migrations successfully!\n", n)
	return nil
}
