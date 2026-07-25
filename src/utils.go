package src

import (
	"strings"
	"time"
)

func GetToday() (string, string) {
	date := strings.Split(time.Now().Format(time.DateOnly), "-")
	month := date[1]
	day := date[2]

	return month, day
}
