package src

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func Scheduler(session *discordgo.Session, channelID string) {
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
	BirthdayMessage(session, channelID)

	// Set ticker for next 24 hours
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		BirthdayMessage(session, channelID)
	}
}
