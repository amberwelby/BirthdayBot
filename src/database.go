package src

import (
	"database/sql"
	"errors"
	"log"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS bithdays (
	first_name TEXT NOT NULL, 
	last_name TEXT NOT NULL, 
	month TEXT NOT NULL, 
	day TEXT NOT NULL, 
	year TEXT,
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

// CRUD Operations
// Create
func Create(db *sql.DB, person *Person) error {
	query := `INSERT INTO birthdays (first_name, last_name, month, day, year) VALUES(?????)`

	result, err := db.Exec(query,
		person.First_name,
		person.Last_name,
		person.Month,
		person.Day,
		person.Year,
		person.Deceased,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	} else {
		log.Printf("Rows created: %v", rows)
	}

	return nil
}

// Retrieve
func RetrieveByName(db *sql.DB, first string, last string) (*Person, error) {
	query := `SELECT * FROM birthdays WHERE first_name = ?, last_name = ?`

	person := &Person{}
	err := db.QueryRow(query, first, last).Scan(
		&person.First_name,
		&person.Last_name,
		&person.Month,
		&person.Day,
		&person.Year,
		&person.Deceased,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	return person, nil
}

func RetrieveByMonth(db *sql.DB, month string) ([]Person, error) {
	query := `SELECT * FROM birthdays WHERE month = ?`

	people := []Person{}
	rows, err := db.Query(query, month)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("No birthdays found")
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if rows.Err() == nil {
			person := Person{}
			err = rows.Scan(
				person.First_name,
				person.Last_name,
				person.Month,
				person.Day,
				person.Year,
				person.Deceased,
			)
			if err != nil {
				log.Println(err)
			}
			people = append(people, person)
		}
	}

	return people, nil
}

func RetrieveByYear(db *sql.DB, year string) ([]Person, error) {
	query := `SELECT * FROM birthdays WHERE year = ?`

	people := []Person{}
	rows, err := db.Query(query, year)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("No birthdays found")
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if rows.Err() == nil {
			person := Person{}
			err = rows.Scan(
				person.First_name,
				person.Last_name,
				person.Month,
				person.Day,
				person.Year,
				person.Deceased,
			)
			if err != nil {
				log.Println(err)
			}
			people = append(people, person)
		}
	}

	return people, nil
}

func RetrieveByDate(db *sql.DB, month string, day string) ([]Person, error) {
	query := `SELECT * FROM birthdays WHERE month = ?, day = ?`

	people := []Person{}
	rows, err := db.Query(query, month, day)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("No birthdays found")
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if rows.Err() == nil {
			person := Person{}
			err = rows.Scan(
				person.First_name,
				person.Last_name,
				person.Month,
				person.Day,
				person.Year,
				person.Deceased,
			)
			if err != nil {
				log.Println(err)
			}
			people = append(people, person)
		}
	}

	return people, nil
}

// Update
func Update(db *sql.DB, person *Person) error {
	query := `UPDATE birthdays 
	SET month = ?, day = ?, year = ?, deceased = ?
	WHERE first_name = ?, last_name = ?
	`

	result, err := db.Exec(query,
		person.Month,
		person.Day,
		person.Year,
		person.Deceased,
		person.First_name,
		person.Last_name,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Update not completed")
	}

	return nil
}

// Delete
func Delete(db *sql.DB, first string, last string) error {
	query := `DELETE FROM birthdays WHERE first_name = ?, last_name = ?`

	result, err := db.Exec(query, first, last)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Delete not completed")
	}

	return nil

}
