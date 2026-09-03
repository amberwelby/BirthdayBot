package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"birthday-bot/src"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	// The underscore import registers the driver with database/sql
	_ "github.com/mattn/go-sqlite3"
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
	db, err := src.NewDatabase("./birthdays.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := src.InitializeSchema(db); err != nil {
		log.Fatalf("Failed to initilize schema: %v", err)
	}

	// Handle recieved messages
	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "Birthdays?" {
			month, day := src.GetToday()
			birthdays := src.GetBirthDates(month, day)
			if len(birthdays) == 0 {
				s.ChannelMessageSend(m.ChannelID, "No birthdays today")
			} else {
				for _, person := range birthdays {
					message := fmt.Sprintf("%v %v", person.First_name, person.Last_name)
					s.ChannelMessageSend(m.ChannelID, message)
				}
			}
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
	go src.Scheduler(sess, channelID)

	// Handle closing
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
