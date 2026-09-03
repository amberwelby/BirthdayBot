package src

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func BirthdayMessage(session *discordgo.Session, channelID string) {
	month, day := GetToday()
	message := fmt.Sprintf("%s/%s \n", month, day)

	birthdays := GetBirthDates(month, day)
	if len(birthdays) == 0 {
		message += "No birthdays today"
	} else {
		for _, name := range birthdays {
			message += fmt.Sprintln(name)
		}
	}

	session.ChannelMessageSend(channelID, message)
}

func GetBirthDates(month string, day string) []Person {
	// Open database
	db, err := NewDatabase("./birthdays.db")
	if err != nil {
		log.Printf("Database error: %v", err)
		return nil
	}

	// Look up birthdays
	names, err := RetrieveByDate(db, month, day)
	if err != nil {
		log.Printf("Retrieve error: %v", err)
		return nil
	}

	return names
}
