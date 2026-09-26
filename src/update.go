package src

import (
	"database/sql"
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
)

func UpdatePerson(db *sql.DB, person *Person) error {
	exists, err := RetrieveByName(db, person.First_name, person.Last_name)
	if exists == nil {
		return err
	}

	query := `UPDATE birthdays 
	SET month = ?, day = ?, year = ?, deceased = ?
	WHERE first_name = ? AND last_name = ?
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

func UpdateSchedule(db *sql.DB, session *discordgo.Session, channelID string, hour int) {
	// How long until next runtime
	location, _ := time.LoadLocation("Local")
	now := time.Now().Local()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, location)

	// If that's already passed today, set for tomorrow
	if now.After(nextRun) {
		nextRun = nextRun.Add(24 * time.Hour)
	}

	// Sleep until runtime
	initialDelay := time.Until(nextRun)
	time.Sleep(initialDelay)

	// Schedule message
	BirthdayMessage(db, session, channelID)

	// Set ticker for next 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		BirthdayMessage(db, session, channelID)
	}
}
