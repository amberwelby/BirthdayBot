package src

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS birthdays (
	first_name TEXT NOT NULL, 
	last_name TEXT NOT NULL, 
	month TEXT NOT NULL, 
	day TEXT NOT NULL, 
	year TEXT DEFAULT 'NA',
	deceased BOOL DEFAULT 0,
	PRIMARY KEY (first_name, last_name)
);
`

type Person struct {
	First_name string
	Last_name  string
	Month      string
	Day        string
	Year       string
	Deceased   bool
}

// Open configured SQLite connection
func NewDatabase(dbPath string) (*sql.DB, error) {
	// Connection string with configuration options
	// _journal_mode=WAL enables Write-Ahead Logging for better concurrency
	// _busy_timeout=5000 waits up to 5 seconds when database is locked
	// _synchronous=NORMAL balances safety and performance
	// _cache_size=-64000 sets 64MB cache (negative = KB)
	// _foreign_keys=ON enables foreign key constraint enforcement
	dsn := dbPath + "?_journal_mode=WAL&_busy_timeout=5000&_synchronous=NORMAL&_cache_size=-64000&_foreign_keys=ON"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("%v", err)
	}

	db.SetMaxOpenConns(1) // Support only 1 writer at a time
	db.SetConnMaxIdleTime(1)
	// Didn't set max connection time

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// Initialize if necessary
func InitializeSchema(db *sql.DB) error {
	_, err := db.Exec(schema)
	return err
}
