package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func RunMigrations(db *sql.DB) error {
	content, err := migrationsFS.ReadFile("migrations/0001_init.up.sql")
	if err != nil {
		return fmt.Errorf("read migration: %w", err)
	}

	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("run migration: %w", err)
	}

	log.Println("Database migrations applied successfully")
	return nil
}
