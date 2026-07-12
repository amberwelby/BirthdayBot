package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"encoding/json"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Get bot token from token.txt
	tokenFile, err := os.ReadFile("token.txt")
	if err != nil {
		log.Fatalf("Token not found: %s", err)
	}

	token := string(tokenFile)

	// Handle json file
	birthdayFile, err := os.Open("birthdays.json")

	if err != nil {
		fmt.Println("Read error")
		log.Fatal(err)
	}

	defer birthdayFile.Close()

	/// Unmarshalling JSON
	byteValue, err := io.ReadAll(birthdayFile)
	if err != nil {
		fmt.Println("Byte string error")
		log.Fatal(err)
	}

	var months map[string]map[string][]string

	err = json.Unmarshal(byteValue, &months)
	if err != nil {
		fmt.Println("Unmarshal Error")
		log.Fatal(err)
	}

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
			names := months["January"]["01"]
			for _, name := range names {
				s.ChannelMessageSend(m.ChannelID, name)
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

	// Handle closing
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
