package src

import (
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func BirthdayMessage(db *sql.DB, session *discordgo.Session, channelID string) error {
	month, day := GetToday()
	message := fmt.Sprintf("%s/%s \n", month, day)

	birthdays, err := RetrieveByDate(db, month, day)
	if err != nil {
		return err
	}

	if len(birthdays) == 0 {
		message += "No birthdays today"
	} else {
		for _, person := range birthdays {
			message += fmt.Sprintln(person.First_name, person.Last_name)
		}
	}

	session.ChannelMessageSend(channelID, message)

	return nil
}

func PrintPerson(person *Person, session *discordgo.Session, channelID string){
	message := ""
	year := ""

	if person.Year != "NA" {
		year = fmt.Sprintf("%s/", person.Year)
	}
	message += fmt.Sprintf("%s %s %s%s/%s", person.First_name, person.Last_name, year, person.Month, person.Day)

	if person.Deceased {
		message += fmt.Sprintf("%s, deceased", message)
	} 

	session.ChannelMessageSend(channelID, message) 			
}

func PrintPersonName(person *Person, session *discordgo.Session, channelID string) {
	message := fmt.Sprintf("%s %s\n", person.First_name, person.Last_name)
	session.ChannelMessageSend(channelID, message) 
}

func PrintPersonDate(person *Person, session *discordgo.Session, channelID string) {
	message := fmt.Sprintf("%s/%s\n", person.Month, person.Day)
	if (person.Year != "NA") {
		message = fmt.Sprintf("%s/%s", person.Year, message)
	}
	session.ChannelMessageSend(channelID, message) 
}

func PrintPeople(people []*Person, session *discordgo.Session, channelID string){
	message := ""

	if len(people) == 0 {
		message += "No people found"
	} else {
		for _, person := range people {
			PrintPerson(person, session, channelID)
		}
	}
	
	session.ChannelMessageSend(channelID, message)
}

func Help() string {
	message := "1. \"/addbirthday first last yyyy/mm/dd\" adds a birthday entry\n 2. \"/getperson first last\" retrieves a person record\n 3. \"/getdate mm dd\" retrieves all people born on that date\n 4. \"/getmonth mm\" retrieves all people born in the month\n 5. \"/getyear yyyy\" retrieves all people born in that year\n 6. \"/update first last mm dd {yyyy, deceased}\" updates record of the person with name \"first last\"\n 7. \"/schedule h\" updates the scheduler to the 24h time specified\n 8. \"/remove first last\" deletes the record of a person with name \"first last\"\n 9. \"/help\" will display the list of available commands and their arguments\n "

	return message
}