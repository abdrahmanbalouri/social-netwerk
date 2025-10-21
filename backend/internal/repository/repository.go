package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
)

var Db *sql.DB

const dbPath = "internal/repository/forum.db"

func OpenDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=1")
	if err != nil {
		return db, err
	}
	db.SetMaxOpenConns(10)
	Db = db

	return db, nil
}

func ApplyMigrations(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "pkg/migrations/sqlite",
	}

	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("error while executing the migration: %v", err)
	}
	log.Printf("Applied %d migrations successfully!\n", n)
	return nil
}
