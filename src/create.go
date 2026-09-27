package src

import (
	"database/sql"
	"fmt"
)

func Create(db *sql.DB, person *Person) string {
	var result sql.Result
	var err error
	var message string

	if (person.Year == ""){
		person.Year = "NA"
	}

	query := `INSERT INTO birthdays (first_name, last_name, month, day, year, deceased) VALUES (?, ?, ?, ?, ?, ?)`
	result, err = db.Exec(query,
		person.First_name,
		person.Last_name,
		person.Month,
		person.Day,
		person.Year,
		person.Deceased,
	)

	if err != nil {
		message = fmt.Sprintf("Error on create: %s", err)
		return message
	}

	rows, err := result.RowsAffected()
	if err != nil {
		message = fmt.Sprintf("Error on validation: %s", err)
	} else {
		message = fmt.Sprintf("%s %s, %d rows added", person.First_name, person.Last_name, rows)
	}

	return message
}
