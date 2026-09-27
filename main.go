package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"birthday-bot/src"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Get environment variables
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")
	channelID := os.Getenv("CHANNEL_ID")

	// Create new Discord session
	sess, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		log.Fatal(err)
	}

	// Open database connection
	db, err := src.NewDatabase("./data/birthdays.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := src.InitializeSchema(db); err != nil {
		log.Fatalf("Failed to initilize schema: %v", err)
	}

	// Handle recieved messages
	sess.AddHandler(func(session *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == session.State.User.ID {
			return
		}

		// Slash commands
		prefix := "/"
		args := strings.Split(m.Content, " ")

		// First char in first input word
		if string(args[0][0]) != prefix {
			return
		}

		command := args[0][1:]

		switch command {
		// Create
		case "addbirthday":
			person := &src.Person{}
			person.First_name = args[1]
			person.Last_name = args[2]
			person.Month, person.Day, person.Year, err = src.ValidateDate(args[3])
			if err != nil {
				session.ChannelMessageSend(channelID, err.Error())
			}
			if len(args) == 4 {
				person.Deceased = true
			} else {
				person.Deceased = false
			}
			message := src.Create(db, person)
			session.ChannelMessageSend(channelID, message)
		// Read
		case "getperson":
			person, err := src.RetrieveByName(db, args[1], args[2])
			if err != nil {
				session.ChannelMessageSend(channelID, err.Error())
			} else {
				src.PrintPersonDate(person, session, channelID)
			}
		case "getdate":
			month, day, _, err := src.ValidateDate(args[1])
			people, err := src.RetrieveByDate(db, month, day)
			if err != nil {
				session.ChannelMessageSend(channelID, err.Error())
			} else {
				src.PrintPeople(people, session, channelID)
			}
		case "getmonth":
			month, _, err := src.ValidateMonth(args[1])
			if err != nil {
				session.ChannelMessageSend(channelID, err.Error())
			}
			people, err := src.RetrieveByMonth(db, month)
			if err != nil {
				session.ChannelMessageSend(channelID, err.Error())
			} else {
				src.PrintPeople(people, session, channelID)
			}
		case "getyear":
			year, err := src.ValidateYear(args[1])
			if err != nil {
				session.ChannelMessageSend(channelID, err.Error())
			}
			people, err := src.RetrieveByYear(db, year)
			if err != nil {
				session.ChannelMessageSend(channelID, err.Error())
			} else {
				src.PrintPeople(people, session, channelID)
			}
		// Update
		case "update":
			person := &src.Person{}
			person.First_name = args[1]
			person.Last_name = args[2]
			person.Month, person.Day, person.Year, err = src.ValidateDate(args[3])
			if err != nil {
				session.ChannelMessageSend(channelID, err.Error())
			}
			if len(args) == 4 {
				person.Deceased = true
			}
			err := src.UpdatePerson(db, person)	
			if err != nil {
				session.ChannelMessageSend(channelID, err.Error())
			} else {
				person, err = src.RetrieveByName(db, person.First_name, person.Last_name)
				if err != nil {
					session.ChannelMessageSend(channelID, err.Error())
				} else {
					src.PrintPerson(person, session, channelID)
				}
			}		
		case "schedule":
			hour, err := strconv.Atoi(args[1])
			if err != nil {
				session.ChannelMessageSend(channelID, err.Error())
				hour = 8
			}
			if hour > 23 || hour < 0 {
				session.ChannelMessageSend(channelID, "Given hour outside of 0-23. Defaulted to 8am")
				hour = 8
			}
			src.UpdateSchedule(db, session, channelID, hour)
		// Delete
		case "remove":
			err := src.Delete(db, args[1], args[2])
			if err != nil {
				session.ChannelMessageSend(channelID, err.Error())
			} else {
				session.ChannelMessageSend(channelID, "User deleted")
			}
		// Help
		case "help":
			message := src.Help()
			session.ChannelMessageSend(channelID, message)
		default:
			session.ChannelMessageSend(channelID, "Command not recognized, try /help for options")
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Open websocket
	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	fmt.Println("Birthday bot is online!")

	// Run daily reminder as goroutine/concurrent thread
	go src.UpdateSchedule(db, sess, channelID, 8)

	// Handle closing
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
