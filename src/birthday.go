package src

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

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

func GetBirthDates(month string, day string) []string {
	// Handle json file
	birthdayFile, err := os.Open("birthdays.json")

	if err != nil {
		fmt.Println("Read error")
		log.Fatal(err)
	}

	defer birthdayFile.Close()

	// Unmarshalling JSON
	byteValue, err := io.ReadAll(birthdayFile)
	if err != nil {
		fmt.Println("Byte string error")
		log.Fatal(err)
	}

	birthdayFile.Close()

	var birthdays map[string]map[string][]string

	err = json.Unmarshal(byteValue, &birthdays)
	if err != nil {
		fmt.Println("Unmarshal Error")
		log.Fatal(err)
	}

	// Look up birthday
	names := birthdays[month][day]

	return names
}
