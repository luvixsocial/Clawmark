package database

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func main() {
	userID, err := uuid.Parse("user-123")
	if err != nil {
		log.Fatal("Invalid UUID:", err)
	}
	limit := 10

	feed, fullPosts, err := GetUserFeed(userID, limit)
	if err != nil {
		log.Fatal("Failed to fetch feed:", err)
	}

	fmt.Println("Generated Feed IDs:", feed)
	fmt.Println("Generated Full Posts:", fullPosts)

	// Simulate real-time update
	go func() {
		for {
			time.Sleep(10 * time.Second)
			notifyUser(userID)
		}
	}()
}
