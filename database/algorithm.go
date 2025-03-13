package database

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"clawmark/state"
	"clawmark/types"
)

func GetUserFeed(userID uuid.UUID, limit int) ([]uuid.UUID, []types.Post, error) {
	personalized, err := getPersonalizedFeed(userID, limit/2)
	if err != nil {
		log.Println("Error fetching personalized feed:", err)
	}

	randomPosts, err := getRandomPosts(limit / 2)
	if err != nil {
		log.Println("Error fetching random posts:", err)
	}

	postSet := make(map[uuid.UUID]bool)
	finalFeed := []uuid.UUID{}

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

func getPersonalizedFeed(userID uuid.UUID, limit int) ([]uuid.UUID, error) {
	var tags []string
	err := state.Redis.ZRevRange(state.Context, fmt.Sprintf("user:%s:tag_scores", userID), 0, 4).ScanSlice(&tags)
	if err != nil {
		return nil, err
	}

	var posts []types.Post
	err = state.Pool.Where("tags && ?", tags).Order("created_at DESC").Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	var postIDs []uuid.UUID
	for _, post := range posts {
		score := computePersonalizedScore(userID, post.ID, post.Tags)
		state.Redis.ZAdd(state.Context, fmt.Sprintf("user:%s:feed", userID), redis.Z{
			Score:  score,
			Member: post.ID,
		})
		postIDs = append(postIDs, post.ID)
	}

	return postIDs, nil
}

func getRandomPosts(limit int) ([]uuid.UUID, error) {
	var posts []types.Post
	err := state.Pool.Order("RANDOM()").Limit(limit).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	var postIDs []uuid.UUID
	for _, post := range posts {
		postIDs = append(postIDs, post.ID)
	}

	return postIDs, nil
}

func getPostsByID(postIDs []uuid.UUID) ([]types.Post, error) {
	if len(postIDs) == 0 {
		return []types.Post{}, nil
	}

	var posts []types.Post
	err := state.Pool.Where("id IN ?", postIDs).Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func computePersonalizedScore(userID uuid.UUID, postID uuid.UUID, tags []string) float64 {
	var interactions []types.Like
	err := state.Pool.Where("user_id = ?", userID).Find(&interactions).Error
	if err != nil {
		log.Println("Error fetching interactions:", err)
		return 0
	}

	interactionBoost := 0.0
	for _, interaction := range interactions {
		if interaction.PostID == postID {
			interactionBoost = 5.0
			break
		}
	}

	tagMatchScore := 0.0
	for _, tag := range tags {
		count, _ := state.Redis.ZScore(state.Context, fmt.Sprintf("user:%s:tag_scores", userID), tag).Result()
		tagMatchScore += count
	}

	return tagMatchScore + interactionBoost
}

func notifyUser(userID uuid.UUID) {
	feed, fullPosts, err := GetUserFeed(userID, 10)
	if err != nil {
		log.Println("Error notifying user:", err)
		return
	}

	message := fmt.Sprintf("update:%v - %v", feed, fullPosts)
	notifyClientsForUser(userID, message)
}

func notifyClientsForUser(userID uuid.UUID, message string) {
	fmt.Printf("Sending real-time update to user %s: %s\n", userID, message)
}
