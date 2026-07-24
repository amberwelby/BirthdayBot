package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func birthdayMessage(s *discordgo.Session, channelID string) {
	month, day := getToday()
	message := fmt.Sprintf("%s/%s \n", month, day)

	birthdays := getBirthDates(month, day)
	if len(birthdays) == 0 {
		message += "No birthdays today"
	} else {
		for _, name := range birthdays {
			message += fmt.Sprintln(name)
		}
	}

	s.ChannelMessageSend(channelID, message)
}

func getBirthDates(month string, day string) []string {
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

func getToday() (string, string) {
	date := strings.Split(time.Now().Format(time.DateOnly), "-")
	month := date[1]
	day := date[2]

	return month, day
}

func scheduler(s *discordgo.Session, channelID string) {
	// How long until next runtime
	location, _ := time.LoadLocation("Local")
	now := time.Now().Local()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), 23, 45, 0, 0, location)

	// If that's already passed today, set for tomorrow
	if now.After(nextRun) {
		nextRun = nextRun.Add(24 * time.Hour)
	}

	// Sleep until runtime
	initialDelay := time.Until(nextRun)
	time.Sleep(initialDelay)

	// Schedule message
	birthdayMessage(s, channelID)

	// Set ticker for next 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		birthdayMessage(s, channelID)
	}
}

func main() {
	// Get environment variables
	godotenv.Load()
	token := os.Getenv("BOT_TOKEN")
	channelID := os.Getenv("CHANNEL_ID")

	// Create new session
	sess, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		log.Fatal(err)
	}

	// Handle recieved messages
	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "Birthdays?" {
			month, day := getToday()
			birthdays := getBirthDates(month, day)
			if len(birthdays) == 0 {
				s.ChannelMessageSend(m.ChannelID, "No birthdays today")
			} else {
				for _, name := range birthdays {
					s.ChannelMessageSend(m.ChannelID, name)
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
	go scheduler(sess, channelID)

	// Handle closing
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
