package database

import (
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"clawmark/state"
	"clawmark/types"
)

func GetUserFeed(userID string, limit int) ([]int, []types.Post, error) {
	personalized, err := getPersonalizedFeed(userID, limit/2)
	if err != nil {
		log.Println("Error fetching personalized feed:", err)
	}

	randomPosts, err := getRandomPosts(limit / 2)
	if err != nil {
		log.Println("Error fetching random posts:", err)
	}

	postSet := make(map[int]bool)
	finalFeed := []int{}

	for _, post := range personalized {
		postSet[post] = true
		finalFeed = append(finalFeed, post)
	}

	for _, post := range randomPosts {
		if !postSet[post] {
			finalFeed = append(finalFeed, post)
		}
	}

	fullPosts, err := getPostsByID(finalFeed)
	if err != nil {
		return nil, nil, err
	}

	return finalFeed, fullPosts, nil
}

func getPersonalizedFeed(userID string, limit int) ([]int, error) {
	tags, err := state.Redis.ZRevRange(state.Context, fmt.Sprintf("user:%s:tag_scores", userID), 0, 4).Result()
	if err != nil {
		return nil, err
	}

	query := `SELECT id, tags FROM posts WHERE tags && $1 ORDER BY created_at DESC LIMIT $2`
	rows, err := state.Pool.Query(state.Context, query, tags, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []int
	for rows.Next() {
		var postID int
		var postTags []string
		if err := rows.Scan(&postID, &postTags); err != nil {
			log.Println("Error scanning post:", err)
			continue
		}

		score := computePersonalizedScore(userID, postID, postTags)

		state.Redis.ZAdd(state.Context, fmt.Sprintf("user:%s:feed", userID), redis.Z{
			Score:  score,
			Member: postID,
		})

		posts = append(posts, postID)
	}

	postIDsStr, err := state.Redis.ZRevRange(state.Context, fmt.Sprintf("user:%s:feed", userID), 0, int64(limit)-1).Result()
	if err != nil {
		return nil, err
	}

	var postIDs []int
	for _, postIDStr := range postIDsStr {
		var postID int
		if _, err := fmt.Sscanf(postIDStr, "%d", &postID); err != nil {
			log.Println("Error parsing post ID:", err)
			continue
		}
		postIDs = append(postIDs, postID)
	}

	return postIDs, nil
}

func getRandomPosts(limit int) ([]int, error) {
	rows, err := state.Pool.Query(state.Context, "SELECT id FROM posts ORDER BY RANDOM() LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postIDs []int
	for rows.Next() {
		var postID int
		if err := rows.Scan(&postID); err != nil {
			log.Println("Error scanning post ID:", err)
			continue
		}
		postIDs = append(postIDs, postID)
	}

	return postIDs, nil
}

func getPostsByID(postIDs []int) ([]types.Post, error) {
	if len(postIDs) == 0 {
		return []types.Post{}, nil
	}

	query := `SELECT id, title, content, tags, created_at FROM posts WHERE id = ANY($1)`
	rows, err := state.Pool.Query(state.Context, query, postIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []types.Post
	for rows.Next() {
		var p types.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Tags, &p.CreatedAt); err != nil {
			log.Println("Error scanning full post:", err)
			continue
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func computePersonalizedScore(userID string, postID int, tags []string) float64 {
	rows, err := state.Pool.Query(state.Context, "SELECT post_id FROM user_interactions WHERE user_id = $1", userID)
	if err != nil {
		log.Println("Error fetching interactions:", err)
		return 0
	}
	defer rows.Close()

	interactionHistory := make(map[int]bool)
	for rows.Next() {
		var interactedPostID int
		if err := rows.Scan(&interactedPostID); err != nil {
			log.Println("Error scanning interaction:", err)
			continue
		}
		interactionHistory[interactedPostID] = true
	}

	interactionBoost := 0.0
	if interactionHistory[postID] {
		interactionBoost = 5.0
	}

	tagMatchScore := 0.0
	for _, tag := range tags {
		count, _ := state.Redis.ZScore(state.Context, fmt.Sprintf("user:%s:tag_scores", userID), tag).Result()
		tagMatchScore += count
	}

	return tagMatchScore + interactionBoost
}

func notifyUser(userID string) {
	feed, fullPosts, err := GetUserFeed(userID, 10)
	if err != nil {
		log.Println("Error notifying user:", err)
		return
	}

	message := fmt.Sprintf("update:%v - %v", feed, fullPosts)
	notifyClientsForUser(userID, message)
}

func notifyClientsForUser(userID, message string) {
	fmt.Printf("Sending real-time update to user %s: %s\n", userID, message)
}
