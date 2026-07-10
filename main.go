package main

import (
	"log"
	"os"
	"fmt"
)

func main() {
	// Get bot token from token.txt
	tokenFile, err := os.ReadFile("token.txt")
	if err != nil {
		log.Fatalf("Token not found: %s", err)
	}

	token := string(tokenFile)
	fmt.Println(token)
}
