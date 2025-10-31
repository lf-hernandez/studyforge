package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Database wraps the SQL database connection
type Database struct {
	DB *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	log.Println("Database connection established")

	return &Database{DB: db}, nil
}

// RunMigrations executes database migrations
func (d *Database) RunMigrations(migrationsPath string) error {
	log.Println("Running database migrations...")

	// Read migration file
	migrationFile := filepath.Join(migrationsPath, "001_init_schema.sql")
	content, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	if _, err := d.DB.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	log.Println("Closing database connection")
	return d.DB.Close()
}
