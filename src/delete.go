package src

import (
	"database/sql"
	"errors"
)

func Delete(db *sql.DB, first string, last string) error {
	exists, err := RetrieveByName(db, first, last)
	if exists == nil {
		return err
	}

	query := `DELETE FROM birthdays WHERE first_name = ? AND last_name = ?`

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
