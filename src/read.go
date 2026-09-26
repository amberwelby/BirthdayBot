package src

import (
	"database/sql"
	"errors"
	"log"
)

func RetrieveByName(db *sql.DB, first string, last string) (*Person, error) {
	query := `SELECT * FROM birthdays WHERE first_name = ? AND last_name = ?`

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

func RetrieveByMonth(db *sql.DB, month string) ([]*Person, error) {
	query := `SELECT * FROM birthdays WHERE month = ?`

	people := []*Person{}
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
				&person.First_name,
				&person.Last_name,
				&person.Month,
				&person.Day,
				&person.Year,
				&person.Deceased,
			)
			if err != nil {
				log.Println(err)
			}
			people = append(people, &person)
		}
	}

	return people, nil
}

func RetrieveByYear(db *sql.DB, year string) ([]*Person, error) {
	query := `SELECT * FROM birthdays WHERE year = ?`

	people := []*Person{}
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
				&person.First_name,
				&person.Last_name,
				&person.Month,
				&person.Day,
				&person.Year,
				&person.Deceased,
			)
			if err != nil {
				log.Println(err)
			}
			people = append(people, &person)
		}
	}

	return people, nil
}

func RetrieveByDate(db *sql.DB, month string, day string) ([]*Person, error) {
	query := `SELECT * FROM birthdays WHERE month = ? AND day = ?`

	people := []*Person{}
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
				&person.First_name,
				&person.Last_name,
				&person.Month,
				&person.Day,
				&person.Year,
				&person.Deceased,
			)
			if err != nil {
				log.Println(err)
			}
			people = append(people, &person)
		}
	}

	return people, nil
}
