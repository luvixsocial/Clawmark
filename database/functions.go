package database

import (
	"clawmark/state"
	"fmt"
	"log"
	"time"
)

func main() {
	state.Pool.Config()

	userID := "user-123"
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
